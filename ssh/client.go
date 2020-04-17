package ssh

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path"

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

	c, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(publicKeys(keyPaths)),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: c,
	}, nil
}

func (c Client) Run(cmd string) error {
	sess, err := c.Client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr
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

	sigch := make(chan os.Signal)
	signal.Notify(sigch, os.Interrupt)
	defer close(sigch)
	go func() {
		for sig := range sigch {
			switch sig {
			case os.Interrupt:
				sess.Signal(ssh.SIGINT)
			}
		}
	}()

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

func (c Client) Push(name string, mode os.FileMode, size int64, r io.Reader) error {
	sess, err := c.Client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	w, err := sess.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer w.Close()
		fmt.Fprintf(w, "C%#o %d %s\n", mode, size, path.Base(name))
		io.Copy(w, r)
		fmt.Fprint(w, "\x00")
	}()

	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr
	return sess.Run(fmt.Sprintf("scp -t %s", name))
}

func publicKeys(paths []string) func() ([]ssh.Signer, error) {
	return func() ([]ssh.Signer, error) {
		signers := []ssh.Signer{}
		for _, path := range paths {
			key, err := ioutil.ReadFile(path)
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
