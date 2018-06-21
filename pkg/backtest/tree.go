package backtest

// NodeHandler defines the basic node functionality.
type NodeHandler interface {
	Name() string
	SetName(string) NodeHandler
	Parent() (NodeHandler, bool)
	SetParent(NodeHandler) NodeHandler
	Children() ([]NodeHandler, bool)
	SetChildren(...NodeHandler) NodeHandler
	IsRoot() bool
	IsChild() bool
	Run()
}

// Node implements NodeHandler. It represents the base information of each tree node.
// This is the main building block of the tree.
type Node struct {
	name     string
	parent   NodeHandler
	children []NodeHandler
}

// Name returns the name of Node
func (n Node) Name() string {
	return n.name
}

// SetName sets the name of Node
func (n *Node) SetName(s string) NodeHandler {
	n.name = s
	return n
}

// Parent return the parent of this Node
func (n Node) Parent() (NodeHandler, bool) {
	if n.parent == nil {
		return &Node{}, false
	}
	return n.parent, true
}

// SetParent sets the parent of this Node
func (n *Node) SetParent(p NodeHandler) NodeHandler {
	n.parent = p
	return n
}

// Children returns the children of this Node
func (n Node) Children() ([]NodeHandler, bool) {
	if n.children == nil {
		return []NodeHandler{}, false
	}
	return n.children, true
}

// SetChildren sets the Children of this Node
func (n *Node) SetChildren(children ...NodeHandler) NodeHandler {
	for _, child := range children {
		child.SetParent(n)
	}
	n.children = children
	return n
}

// IsRoot checks if this Node is a root node
func (n Node) IsRoot() bool {
	if n.parent != nil {
		return false
	}
	return true
}

// IsChild checks if this Node is a child of another node
func (n Node) IsChild() bool {
	if n.parent == nil {
		return false
	}
	return true
}

// Run is an empty function to satisfy the interface
func (n Node) Run() {}
