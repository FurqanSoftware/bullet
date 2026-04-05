package core

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strconv"
	"time"
)

// Release represents a deployable release, identified by timestamp and content hash.
type Release struct {
	Time    string
	Hash    string
	Tarball Tarball
}

// Tarball holds the local path and size of a release archive.
type Tarball struct {
	Path string
	Size int64
}

// NewRelease creates a Release from a tarball path, computing its SHA256 hash.
func NewRelease(tarPath string) (*Release, error) {
	fi, err := os.Stat(tarPath)
	if err != nil {
		return nil, err
	}

	s, err := sha256Tarball(tarPath)
	if err != nil {
		return nil, err
	}

	return &Release{
		Time: strconv.FormatInt(time.Now().Unix(), 10),
		Hash: s,
		Tarball: Tarball{
			Path: tarPath,
			Size: fi.Size(),
		},
	}, nil
}

func sha256Tarball(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	s := h.Sum(nil)

	return hex.EncodeToString(s[:]), nil
}
