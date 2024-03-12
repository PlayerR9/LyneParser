package LyneLexer

type matchNode struct {
	token  *Token
	pos    int
	isDone bool

	parent      *matchNode
	firstChild  *matchNode
	nextSibling *matchNode
}

func newMatchNode(token *Token, pos int) *matchNode {
	return &matchNode{
		token: token,
		pos:   pos,
	}
}

/*
func (n *matchNode) getChildren() []*matchNode {
	if n.firstChild == nil {
		return nil
	}

	children := make([]*matchNode, 0)

	for node := n.firstChild; node != nil; node = node.nextSibling {
		children = append(children, node)
	}

	return children
}
*/

func (n *matchNode) addChildren(children ...*matchNode) {
	if len(children) == 0 {
		return
	}

	var lastChild *matchNode = nil

	if n.firstChild != nil {
		for lastChild = n.firstChild; lastChild.nextSibling != nil; lastChild = lastChild.nextSibling {
		}
	}

	if lastChild == nil {
		n.firstChild = children[0]
	} else {
		lastChild.nextSibling = children[0]
	}

	for i := 1; i < len(children); i++ {
		children[i].parent = n
		children[i-1].nextSibling = children[i]
	}
}

// findBranchingPoint returns the parent and the first sibling that is the highest parent
// of the node that has a sibling.
func (n *matchNode) findBranchingPoint() (*matchNode, *matchNode) {
	if n.parent == nil {
		return nil, nil
	}

	if n.parent.firstChild.nextSibling != nil {
		return n.parent, n
	}

	for node := n.parent; node != nil; node = node.parent {
		if node.nextSibling != nil {
			return node.parent, node
		}
	}

	return nil, nil
}

func (n *matchNode) removeBranch() {
	if n.firstChild == nil {
		// This is a leaf node
		return
	}

	node := n.firstChild

	node.removeBranch()
	node.parent = nil

	prev := node

	for node = node.nextSibling; node != nil; node = node.nextSibling {
		node.removeBranch()

		node.parent = nil

		prev.nextSibling = nil
		prev = node
	}

	n.firstChild = nil
}

func (n *matchNode) getLeaves() []*matchNode {
	if n.firstChild == nil {
		return []*matchNode{n}
	}

	leaves := make([]*matchNode, 0)

	for node := n.firstChild; node != nil; node = node.nextSibling {
		leaves = append(leaves, node.getLeaves()...)
	}

	return leaves
}

/*
func (n *matchNode) traverse(f func(*matchNode)) {
	f(n)

	for node := n.firstChild; node != nil; node = node.nextSibling {
		node.traverse(f)
	}
}
*/

func (n *matchNode) snakeTraversal() [][]*matchNode {
	if n.firstChild == nil {
		return [][]*matchNode{
			{n},
		}
	}

	result := make([][]*matchNode, 0)

	for node := n.firstChild; node != nil; node = node.nextSibling {
		for _, tmp := range node.snakeTraversal() {
			result = append(result, append([]*matchNode{n}, tmp...))
		}
	}

	return result
}
