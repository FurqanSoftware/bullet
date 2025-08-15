package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/FurqanSoftware/bullet/scope"
)

type Selector struct {
	stdin  io.Reader
	stdout io.Writer
}

func NewSelector() *Selector {
	return &Selector{
		stdin:  os.Stdin,
		stdout: os.Stdout,
	}
}

func (r *Selector) Node(s scope.Scope) scope.Scope {
	if len(s.Nodes) == 1 {
		return s
	}
	selector := 1
	for i, n := range s.Nodes {
		fmt.Fprintf(r.stdout, "%d. %s\n", i+1, n.Label())
	}
	fmt.Fprintf(r.stdout, "? [%d] ", selector)
	fmt.Fscanf(r.stdin, "%d", &selector)
	s.Nodes = []scope.Node{s.Nodes[selector-1]}
	return s
}

func (r *Selector) Nodes(s scope.Scope) scope.Scope {
	if len(s.Nodes) == 1 {
		return s
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
				log.Fatal(err)
			}
			selected = append(selected, s.Nodes[i-1])
		} else {
			parts := strings.SplitN(r, "-", 2)
			l, err := strconv.Atoi(parts[0])
			if err != nil {
				log.Fatal(err)
			}
			r, err := strconv.Atoi(parts[1])
			if err != nil {
				log.Fatal(err)
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
	return s
}
