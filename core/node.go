package core

import (
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type Node struct {
	Host string
	Port int
	Tags []string

	Identity string
}

func ParseNodeSet(hosts string, port int, identity string) ([]Node, error) {
	if strings.HasPrefix(hosts, "@") {
		return parseNodeSetManifest(hosts[1:], port, identity)
	}
	return parseNodeSetString(hosts, port, identity)
}

func (n Node) Addr() string {
	return n.Host + ":" + strconv.Itoa(n.Port)
}

func (n Node) HasTag(tag string) bool {
	for _, t := range n.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (n Node) HasTags(tags []string) bool {
	for _, tag := range tags {
		if !n.HasTag(tag) {
			return false
		}
	}
	return true
}

func parseNodeSetManifest(hosts string, port int, identity string) ([]Node, error) {
	parts := strings.SplitN(hosts, ":", 2)
	b, err := os.ReadFile(parts[0])
	if err != nil {
		return nil, err
	}
	nodes := []Node{}
	err = yaml.Unmarshal(b, &nodes)
	if err != nil {
		return nil, err
	}
	if len(parts) > 1 {
		nodes = filterNodeSet(nodes, strings.Split(parts[1], ","))
	}
	for i := range nodes {
		nodes[i].Identity = identity
	}
	return nodes, nil
}

func parseNodeSetString(hosts string, port int, identity string) ([]Node, error) {
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

func filterNodeSet(nodes []Node, clauses []string) []Node {
	filtered := nodes[:0]
L:
	for _, node := range nodes {
		for _, clause := range clauses {
			if node.HasTags(strings.Split(clause, "+")) {
				filtered = append(filtered, node)
				continue L
			}
		}
	}
	return filtered
}
