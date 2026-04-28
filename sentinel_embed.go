package main

import (
	"embed"
	"fmt"
)

//go:embed embed
var sentinelFS embed.FS

func sentinelBinary(arch string) ([]byte, error) {
	data, err := sentinelFS.ReadFile("embed/bullet-sentinel-linux-" + arch)
	if err != nil {
		return nil, fmt.Errorf("sentinel binary for %q is not embedded; run `task sentinel:embed` and rebuild bullet", arch)
	}
	return data, nil
}
