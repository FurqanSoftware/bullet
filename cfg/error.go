package cfg

import "strconv"

// Error represents a configuration error with the operation name and cause.
type Error struct {
	Name string
	Err  error
}

// Error returns the formatted error string.
func (e *Error) Error() string {
	return "cfg: " + strconv.Quote(e.Name) + ": " + e.Err.Error()
}
