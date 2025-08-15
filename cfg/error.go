package cfg

import "strconv"

type Error struct {
	Name string
	Err  error
}

func (e *Error) Error() string {
	return "cfg: " + strconv.Quote(e.Name) + ": " + e.Err.Error()
}
