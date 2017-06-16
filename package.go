package bullet

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/FurqanSoftware/bullet/spec"
)

type Release struct {
	Name    string
	Tarball Tarball
}

type Tarball struct {
	Path string
}

func Package(spec *spec.Spec) (*Release, error) {
	rel := Release{
		Name: fmt.Sprintf("%s-master-%d", spec.Application.Identifier, time.Now().Unix()),
	}

	tarPath := fmt.Sprintf("/tmp/bullet-%s.tar", rel.Name)
	f, err := os.Create(tarPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rel.Tarball = Tarball{
		Path: tarPath,
	}

	w := tar.NewWriter(f)
	err = makeTarball(w, spec.Package.Contents)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return &rel, nil
}

func makeTarball(w *tar.Writer, paths []string) error {
	for _, p := range paths {
		parts := strings.SplitN(p, ":", 2)
		if len(parts[0]) == 1 {
			parts = append(parts, parts[0])
		}
		src := os.ExpandEnv(parts[1])
		dst := parts[0]

		err := filepath.Walk(src, func(name string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			return addToTarball(w, path.Join(dst, strings.TrimPrefix(name, src)), name, info)
		})
		if err != nil {
			println(err)
			return err
		}
	}
	return nil
}

func addToTarball(w *tar.Writer, dst, src string, info os.FileInfo) error {
	hdr, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	hdr.Name = dst

	err = w.WriteHeader(hdr)
	if err != nil {
		return err
	}

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	return err
}
