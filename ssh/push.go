package ssh

import (
	"io"
)

type PushStatus struct {
	N, Size int64
}

type writeTracker struct {
	w    io.Writer
	size int64
	n    int64
	ch   chan PushStatus
}

func (t *writeTracker) Write(p []byte) (int, error) {
	n, err := t.w.Write(p)
	t.n += int64(n)
	if t.ch != nil {
		select {
		case t.ch <- PushStatus{t.n, t.size}:
		default:
		}
	}
	return n, err
}

func (t *writeTracker) Stop() {
	if t.ch != nil {
		close(t.ch)
	}
}
