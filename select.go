package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/FurqanSoftware/bullet/scope"
)

// Selector prompts the user to select one or more nodes from a scope.
type Selector struct {
	stdin  io.Reader
	stdout io.Writer
}

// NewSelector returns a Selector that reads from stdin and writes to stdout.
func NewSelector() *Selector {
	return &Selector{
		stdin:  os.Stdin,
		stdout: os.Stdout,
	}
}

// Node prompts the user to select a single node from the scope.
func (r *Selector) Node(s scope.Scope) (scope.Scope, error) {
	if len(s.Nodes) == 1 {
		return s, nil
	}
	selector := 1
	for i, n := range s.Nodes {
		fmt.Fprintf(r.stdout, "%d. %s\n", i+1, n.Label())
	}
	fmt.Fprintf(r.stdout, "? [%d] ", selector)
	fmt.Fscanf(r.stdin, "%d", &selector)
	if selector < 1 || selector > len(s.Nodes) {
		return s, fmt.Errorf("invalid node number %d, must be between 1 and %d", selector, len(s.Nodes))
	}
	s.Nodes = []scope.Node{s.Nodes[selector-1]}
	return s, nil
}

// Nodes prompts the user to select one or more nodes from the scope.
func (r *Selector) Nodes(s scope.Scope) (scope.Scope, error) {
	if len(s.Nodes) == 1 {
		return s, nil
	}
	selected := []scope.Node{}
	selector := fmt.Sprintf("1-%d", len(s.Nodes))
	for i, n := range s.Nodes {
		fmt.Fprintf(r.stdout, "%d. %s\n", i+1, n.Label())
	}
	fmt.Fprintf(r.stdout, "? [%s] ", selector)
	fmt.Fscanf(r.stdin, "%s", &selector)
	ranges := strings.Split(selector, ",")
	for _, r := range ranges {
		if !strings.Contains(r, "-") {
			i, err := strconv.Atoi(r)
			if err != nil {
				return s, err
			}
			if i < 1 || i > len(s.Nodes) {
				return s, fmt.Errorf("invalid node number %d, must be between 1 and %d", i, len(s.Nodes))
			}
			selected = append(selected, s.Nodes[i-1])
		} else {
			parts := strings.SplitN(r, "-", 2)
			l, err := strconv.Atoi(parts[0])
			if err != nil {
				return s, err
			}
			r, err := strconv.Atoi(parts[1])
			if err != nil {
				return s, err
			}
			if l < 1 || r > len(s.Nodes) {
				return s, fmt.Errorf("invalid node range %d-%d, must be between 1 and %d", l, r, len(s.Nodes))
			}
			if l > r {
				continue
			}
			for i := l; i <= r; i++ {
				selected = append(selected, s.Nodes[i-1])
			}
		}
	}
	s.Nodes = selected
	return s, nil
}
