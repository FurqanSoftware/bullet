package core

import (
	"strconv"
	"strings"
)

type Node struct {
	Host string
	Port int

	Identity string
}

func ParseNodeSet(hosts string, port int, identity string) ([]Node, error) {
	nodes := []Node{}
	for _, h := range strings.Split(hosts, ",") {
		nodes = append(nodes, Node{
			Host:     h,
			Port:     port,
			Identity: identity,
		})
	}
	return nodes, nil
}

func (n Node) Addr() string {
	return n.Host + ":" + strconv.Itoa(n.Port)
}
