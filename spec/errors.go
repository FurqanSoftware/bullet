package spec

import "strconv"

type Error struct {
	Name string
	Err  error
}

func (e *Error) Error() string {
	return "spec: " + strconv.Quote(e.Name) + ": " + e.Err.Error()
}
