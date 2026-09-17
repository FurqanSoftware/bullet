package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
)

// EnvironPush uploads an environment file to the selected nodes. If the file
// has changed and restart is true, all running containers are recreated so
// that they pick up the new environment.
func EnvironPush(s scope.Scope, g cfg.Configuration, filename string, restart bool) error {
	for _, n := range s.Nodes {
		pog.SetStatus(pogConnecting(n))
		c, err := sshDial(n, g)
		if err != nil {
			return err
		}
		pog.Infof("Connected to %s", n.Label())
		pog.SetStatus(nil)

		d, err := distro.New(c)
		if err != nil {
			return err
		}

		changed, hash, err := pushEnviron(c, d, s, filename)
		if err != nil {
			return err
		}
		if !changed || !restart {
			continue
		}

		_, err = restartNode(d, s)
		if err != nil {
			return err
		}
		err = markEnvironApplied(d, s, hash)
		if err != nil {
			return err
		}
	}
	return nil
}

func pushEnviron(c *ssh.Client, d distro.Distro, s scope.Scope, filename string) (bool, string, error) {
	hash, err := sha256File(filename)
	if err != nil {
		return false, "", err
	}

	applied, _ := d.ReadFile(environHashPath(s))
	appliedHash := strings.TrimSpace(string(applied))
	if appliedHash == "" {
		// No record of an applied environment yet. If the file on the server
		// already matches, record it instead of treating it as a change.
		remoteHash, err := hashRemoteEnviron(d, s)
		if err == nil && remoteHash == hash {
			err = markEnvironApplied(d, s, hash)
			if err != nil {
				return false, "", err
			}
			appliedHash = hash
		}
	}
	if appliedHash == hash {
		pog.Info("Environment file unchanged")
		return false, hash, nil
	}

	err = uploadEnvironFile(c, s, filename)
	if err != nil {
		return false, "", err
	}
	return true, hash, nil
}

func markEnvironApplied(d distro.Distro, s scope.Scope, hash string) error {
	return d.WriteFile(environHashPath(s), []byte(hash))
}

func hashRemoteEnviron(d distro.Distro, s scope.Scope) (string, error) {
	b, err := d.ReadFile(fmt.Sprintf("/opt/%s/env", s.Spec.Application.Identifier))
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:]), nil
}

func environHashPath(s scope.Scope) string {
	return fmt.Sprintf("/opt/%s/env.hash", s.Spec.Application.Identifier)
}

func uploadEnvironFile(c *ssh.Client, s scope.Scope, filename string) error {
	pog.SetStatus(pogText("Uploading environment file"))
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	err = c.Push(fmt.Sprintf("/opt/%s/env", s.Spec.Application.Identifier), fi.Mode(), fi.Size(), f, nil)
	if err != nil {
		return err
	}
	pog.Info("Uploaded environment file")
	pog.SetStatus(nil)
	return nil
}
