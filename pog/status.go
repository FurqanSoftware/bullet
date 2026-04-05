package pog

import "github.com/fatih/color"

// Status implements the pog.Status interface for displaying progress updates.
type Status struct {
	IconVal  byte
	TextVal  string
	ColorVal *color.Color
	ThrobVal bool
}

// Icon returns the status icon character.
func (s Status) Icon() byte { return s.IconVal }

// Text returns the status message.
func (s Status) Text() string { return s.TextVal }

// Color returns the color for the status display.
func (s Status) Color() *color.Color { return s.ColorVal }

// Throb returns whether the status indicator should animate.
func (s Status) Throb() bool { return s.ThrobVal }
