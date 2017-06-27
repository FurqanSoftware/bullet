package ssh

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	Client *ssh.Client
}

func Dial(addr string) (*Client, error) {
	c, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(publicKeys),
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

func publicKeys() ([]ssh.Signer, error) {
	key, err := ioutil.ReadFile(os.ExpandEnv("$HOME/.ssh/id_rsa"))
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	return []ssh.Signer{signer}, nil
}
