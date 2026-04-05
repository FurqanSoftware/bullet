package spec

import "strconv"

// Error represents a spec parsing error with the operation name and cause.
type Error struct {
	Name string
	Err  error
}

// Error returns the formatted error string.
func (e *Error) Error() string {
	return "spec: " + strconv.Quote(e.Name) + ": " + e.Err.Error()
}
