package core

import (
	"fmt"

	bpog "github.com/FurqanSoftware/bullet/pog"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/dustin/go-humanize"
)

func pogConnecting(n scope.Node) bpog.Status {
	return bpog.Status{
		IconVal:  '~',
		TextVal:  fmt.Sprintf("Connecting to %s", n.Label()),
		ThrobVal: true,
	}
}

func pogUploadTarball(n, size int64) bpog.Status {
	return bpog.Status{
		IconVal:  '~',
		TextVal:  fmt.Sprintf("Uploading tarball, %s of %s (%d%%) done", humanize.Bytes(uint64(n)), humanize.Bytes(uint64(size)), n*100/size),
		ThrobVal: true,
	}
}

func pogReloadingContainer(p spec.Program, no int) bpog.Status {
	return bpog.Status{
		IconVal:  '~',
		TextVal:  fmt.Sprintf("Reloading container %s:%d", p.Key, no),
		ThrobVal: true,
	}
}

func pogRestartingContainer(p spec.Program, no int) bpog.Status {
	return bpog.Status{
		IconVal:  '~',
		TextVal:  fmt.Sprintf("Restarting container %s:%d", p.Key, no),
		ThrobVal: true,
	}
}

func pogScalingProgram(p spec.Program) bpog.Status {
	return bpog.Status{
		IconVal:  '~',
		TextVal:  fmt.Sprintf("Scaling program %s", p.Key),
		ThrobVal: true,
	}
}

func pogText(s string) bpog.Status {
	return bpog.Status{
		IconVal:  '~',
		TextVal:  s,
		ThrobVal: true,
	}
}
