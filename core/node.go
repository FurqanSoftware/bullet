package core

import (
	"strconv"
	"strings"
)

type Node struct {
	Host string
	Port int
}

func ParseNodeSet(hosts string) ([]Node, error) {
	nodes := []Node{}
	for _, h := range strings.Split(hosts, ",") {
		nodes = append(nodes, Node{
			Host: h,
			Port: 22,
		})
	}
	return nodes, nil
}

func (n Node) Addr() string {
	return n.Host + ":" + strconv.Itoa(n.Port)
}
