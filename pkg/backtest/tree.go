package backtest

// NodeHandler defines the basic node functionality.
type NodeHandler interface {
	Name() string
	SetName(string) NodeHandler
	Root() bool
	SetRoot(bool)
	Children() ([]NodeHandler, bool)
	SetChildren(...NodeHandler) NodeHandler
	Run() error
}

// Node implements NodeHandler. It represents the base information of each tree node.
// This is the main building block of the tree.
type Node struct {
	root     bool
	name     string
	children []NodeHandler
}

// Name returns the name of Node.
func (n Node) Name() string {
	return n.name
}

// SetName sets the name of Node.
func (n *Node) SetName(s string) NodeHandler {
	n.name = s
	return n
}

// Root checks if this Node is a root node.
func (n Node) Root() bool {
	return n.root
}

// SetRoot sets the root status of this Node.
func (n *Node) SetRoot(b bool) {
	n.root = b
}

// Children returns the children of this Node.
func (n Node) Children() ([]NodeHandler, bool) {
	if n.children == nil {
		return []NodeHandler{}, false
	}
	return n.children, true
}

// SetChildren sets the Children of this Node.
func (n *Node) SetChildren(children ...NodeHandler) NodeHandler {
	for _, child := range children {
		child.SetRoot(false)
	}
	n.children = children
	return n
}

// Run is an empty function to satisfy the interface.
func (n Node) Run() error {
	return nil
}
