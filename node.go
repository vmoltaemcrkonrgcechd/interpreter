package main

type node struct {
	typ      int
	value    string
	children []*node
}

func newNode(typ int, value string, children ...*node) *node {
	return &node{typ: typ, value: value, children: children}
}
