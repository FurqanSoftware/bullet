package core

import (
	"fmt"

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

func pogConnecting(addr string) pogStatus {
	return pogStatus{
		icon:  '~',
		text:  fmt.Sprintf("Connecting to %s", addr),
		throb: true,
	}
}
