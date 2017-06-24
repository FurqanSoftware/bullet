package core

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"
)

type Release struct {
	Hash    string
	Tarball Tarball
}

type Tarball struct {
	Path string
}

func NewRelease(tarPath string) (*Release, error) {
	s, err := sha1Tarball(tarPath)
	if err != nil {
		return nil, err
	}

	return &Release{
		Hash: fmt.Sprintf("%d-%s", time.Now().Unix(), s),
		Tarball: Tarball{
			Path: tarPath,
		},
	}, nil
}

func sha1Tarball(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	s := h.Sum(nil)

	return hex.EncodeToString(s[:]), nil
}
