package LyneParser

import (
	"fmt"

	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

type Node struct {
	id   string
	data string

	parent      *Node
	firstChild  *Node
	nextSibling *Node
}

func (n *Node) String() string {
	if n == nil {
		return "Node[nil]"
	}

	return fmt.Sprintf("Node[id=%s, data=%s]", n.id, n.data)
}

func NewNode(id string, data string) *Node {
	return &Node{
		id:   id,
		data: data,
	}
}

func (p *Node) GetID() string {
	return p.id
}

func (p *Node) GetData() string {
	return p.data
}

func (p *Node) LastChild() *Node {
	if p.firstChild == nil {
		return nil
	}

	var lastChild *Node

	for lastChild = p.firstChild; lastChild.nextSibling != nil; lastChild = lastChild.nextSibling {
	}

	return lastChild
}

func (p *Node) AddChildren(children ...*Node) {
	// 1. Remove the nil children
	children = slext.SliceFilter(children, func(c *Node) bool {
		return c != nil
	})

	if len(children) == 0 {
		return
	}

	// 2. Set the parent of the children
	for _, child := range children {
		child.parent = p
	}

	if p.firstChild == nil {
		p.firstChild, children = children[0], children[1:]
	}

	lastChild := p.LastChild()

	// 4. Set all the siblings of the last child
	for _, child := range children {
		lastChild.nextSibling, lastChild = child, child
	}
}
