package main

import (
	"io"
	"sync"
	"testing"

	"github.com/FurqanSoftware/bullet/scope"
	. "github.com/onsi/gomega"
)

func TestSelectorNode(t *testing.T) {
	for _, c := range []struct {
		nodes []scope.Node
		in    string
		want  scope.Node
	}{
		{
			nodes: nodesAlpha,
			in:    "1\n",
			want:  nodesAlpha[0],
		},
		{
			nodes: nodesAlpha,
			in:    "2\n",
			want:  nodesAlpha[1],
		},
		{
			nodes: nodesAlpha,
			in:    "3\n",
			want:  nodesAlpha[2],
		},
	} {
		t.Run("", func(t *testing.T) {
			g := NewWithT(t)

			stdinr, stdinw := io.Pipe()
			stdoutr, stdoutw := io.Pipe()

			r := Selector{
				stdin:  stdinr,
				stdout: stdoutw,
			}

			wg := sync.WaitGroup{}
			wg.Add(2)
			go func() {
				defer wg.Done()
				io.ReadAll(stdoutr)
			}()
			go func() {
				defer wg.Done()
				io.WriteString(stdinw, c.in)
			}()

			s := r.Node(scope.Scope{
				Nodes: c.nodes,
			})
			g.Expect(s.Nodes).To(HaveLen(1))
			g.Expect(s.Nodes[0]).To(Equal(c.want))

			stdoutw.Close()

			wg.Wait()
		})
	}
}

func TestSelectorNodes(t *testing.T) {
	for _, c := range []struct {
		nodes []scope.Node
		in    string
		want  []scope.Node
	}{
		{
			nodes: nodesAlpha,
			in:    "1\n",
			want:  []scope.Node{nodesAlpha[0]},
		},
		{
			nodes: nodesAlpha,
			in:    "1-2\n",
			want:  []scope.Node{nodesAlpha[0], nodesAlpha[1]},
		},
		{
			nodes: nodesAlpha,
			in:    "1,3\n",
			want:  []scope.Node{nodesAlpha[0], nodesAlpha[2]},
		},
	} {
		t.Run("", func(t *testing.T) {
			g := NewWithT(t)

			stdinr, stdinw := io.Pipe()
			stdoutr, stdoutw := io.Pipe()

			r := Selector{
				stdin:  stdinr,
				stdout: stdoutw,
			}

			wg := sync.WaitGroup{}
			wg.Add(2)
			go func() {
				defer wg.Done()
				io.ReadAll(stdoutr)
			}()
			go func() {
				defer wg.Done()
				io.WriteString(stdinw, c.in)
			}()

			s := r.Nodes(scope.Scope{
				Nodes: c.nodes,
			})
			g.Expect(s.Nodes).To(Equal(c.want))

			stdoutw.Close()

			wg.Wait()
		})
	}
}

var (
	nodesAlpha = []scope.Node{
		{
			Name: "alpha-192.168.0.3",
			Host: "alpha-192.168.0.3",
			Port: 22,
		},
		{
			Name: "alpha-192.168.0.4",
			Host: "alpha-192.168.0.4",
			Port: 22,
		},
		{
			Name: "alpha-192.168.0.5",
			Host: "alpha-192.168.0.5",
			Port: 22,
		},
	}
)
