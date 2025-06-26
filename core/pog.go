package core

import (
	"fmt"

	"github.com/FurqanSoftware/bullet/spec"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
)

type pogStatus struct {
	icon  byte
	text  string
	color *color.Color
	throb bool
}

func (s pogStatus) Icon() byte          { return s.icon }
func (s pogStatus) Text() string        { return s.text }
func (s pogStatus) Color() *color.Color { return s.color }
func (s pogStatus) Throb() bool         { return s.throb }

func pogConnecting(n Node) pogStatus {
	return pogStatus{
		icon:  '~',
		text:  fmt.Sprintf("Connecting to %s", n.Label()),
		throb: true,
	}
}

func pogUploadTarball(n, size int64) pogStatus {
	return pogStatus{
		icon:  '~',
		text:  fmt.Sprintf("Uploading tarball, %s of %s (%d%%) done", humanize.Bytes(uint64(n)), humanize.Bytes(uint64(size)), n*100/size),
		throb: true,
	}
}

func pogReloadingContainer(p spec.Program, no int) pogStatus {
	return pogStatus{
		icon:  '~',
		text:  fmt.Sprintf("Reloading container %s:%d", p.Key, no),
		throb: true,
	}
}

func pogText(s string) pogStatus {
	return pogStatus{
		icon:  '~',
		text:  s,
		throb: true,
	}
}
