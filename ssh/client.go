package ssh

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/pog"
	"github.com/avast/retry-go"
	"github.com/mattn/go-tty"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	Client *ssh.Client
}

func Dial(addr, identity string) (*Client, error) {
	keyPaths := []string{}
	if identity != "" {
		keyPaths = append(keyPaths, identity)
	}
	keyPaths = append(keyPaths, os.ExpandEnv("$HOME/.ssh/id_rsa"))

	var c *ssh.Client
	err := retry.Do(func() (err error) {
		c, err = ssh.Dial("tcp", addr, &ssh.ClientConfig{
			User: "root",
			Auth: []ssh.AuthMethod{
				ssh.PublicKeysCallback(publicKeys(keyPaths)),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
		return
	},
		retry.Attempts(uint(cfg.Current.SSHRetries+1)),
		retry.OnRetry(func(n uint, err error) {
			pog.Debugf(".. Retrying")
		}))
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: c,
	}, nil
}

func (c Client) Run(cmd string, echo bool) error {
	sess, err := c.Client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	if echo {
		sess.Stdout = os.Stdout
		sess.Stderr = os.Stderr
	}
	return sess.Run(cmd)
}

func (c Client) RunPTY(cmd string) error {
	sess, err := c.Client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	tty, err := tty.Open()
	if err != nil {
		return err
	}
	defer tty.Close()

	restore, err := tty.Raw()
	if err != nil {
		return err
	}
	defer restore()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	w, h, err := tty.Size()
	if err != nil {
		return err
	}

	err = sess.RequestPty("xterm", h, w, modes)
	if err != nil {
		return err
	}

	sizeCh := tty.SIGWINCH()
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case size := <-sizeCh:
				sess.WindowChange(size.H, size.W)
			case <-doneCh:
				return
			}
		}
	}()
	defer close(doneCh)

	sess.Stdin = tty.Input()
	sess.Stdout = tty.Output()
	sess.Stderr = tty.Output()
	return sess.Run(cmd)
}

func (c Client) Output(cmd string) ([]byte, error) {
	sess, err := c.Client.NewSession()
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	return sess.Output(cmd)
}

func (c Client) Push(name string, mode os.FileMode, size int64, r io.Reader, chstatus chan PushStatus) error {
	sess, err := c.Client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	w, err := sess.StdinPipe()
	if err != nil {
		return err
	}
	wb := bufio.NewWriter(w)
	go func() {
		wt := writeTracker{w: wb, size: size, ch: chstatus}
		defer wt.Stop()
		defer w.Close()
		fmt.Fprintf(wb, "C%#o %d %s\n", mode, size, path.Base(name))
		wb.Flush()
		io.Copy(&wt, r)
		wb.Flush()
		fmt.Fprint(w, "\x00")
		wb.Flush()
	}()

	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr
	return sess.Run(fmt.Sprintf("scp -t %s", name))
}

func (c Client) Forward(local, remote int) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", local))
	if err != nil {
		return err
	}

	pog.Infof("Listening on :%d", local)
	pog.SetStatus(pogForward(0))
	deltaCh := make(chan int)
	go func() {
		n := 0
		for d := range deltaCh {
			n += d
			pog.SetStatus(pogForward(n))
		}
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		deltaCh <- 1
		go func() {
			c.forwardConnect(conn, remote)
			deltaCh <- -1
		}()
	}
}

func (c Client) forwardConnect(conn net.Conn, remote int) {
	sess, err := c.Client.Dial("tcp", fmt.Sprintf("localhost:%d", remote))
	if err != nil {
		log.Println("dial:", err)
		return
	}
	defer sess.Close()

	go c.forwardPump(conn, sess)
	c.forwardPump(sess, conn)
}

func (c Client) forwardPump(w io.Writer, r io.Reader) {
	_, err := io.Copy(w, r)
	if err != nil {
		log.Println("pump:", err)
	}
}

func publicKeys(paths []string) func() ([]ssh.Signer, error) {
	return func() ([]ssh.Signer, error) {
		signers := []ssh.Signer{}
		for _, path := range paths {
			key, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}

			signer, err := ssh.ParsePrivateKey(key)
			if err != nil {
				return nil, err
			}
			signers = append(signers, signer)
		}
		return signers, nil
	}
}
