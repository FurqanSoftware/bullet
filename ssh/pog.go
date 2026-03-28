package ssh

import (
	"fmt"

	bpog "github.com/FurqanSoftware/bullet/pog"
)

func pogForward(n int) bpog.Status {
	return bpog.Status{
		IconVal:  '~',
		TextVal:  fmt.Sprintf("%d active", n),
		ThrobVal: true,
	}
}
