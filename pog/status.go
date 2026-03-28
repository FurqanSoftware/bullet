package pog

import "github.com/fatih/color"

type Status struct {
	IconVal  byte
	TextVal  string
	ColorVal *color.Color
	ThrobVal bool
}

func (s Status) Icon() byte          { return s.IconVal }
func (s Status) Text() string        { return s.TextVal }
func (s Status) Color() *color.Color { return s.ColorVal }
func (s Status) Throb() bool         { return s.ThrobVal }
