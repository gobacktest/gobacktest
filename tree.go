package gobacktest

// NodeHandler defines the basic node functionality.
type NodeHandler interface {
	WeightHandler
	Name() string
	SetName(string) NodeHandler
	Root() bool
	SetRoot(bool)
	Children() ([]NodeHandler, bool)
	SetChildren(...NodeHandler) NodeHandler
}

// WeightHandler defines weight functionality.
type WeightHandler interface {
	Weight() float64
	SetWeight(float64)
	Tolerance() float64
	SetTolerance(float64)
}

// Node implements NodeHandler. It represents the base information of each tree node.
// This is the main building block of the tree.
type Node struct {
	root      bool
	name      string
	weight    float64
	tolerance float64
	children  []NodeHandler
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

// Weight returns the weight of this node within this strategy level.
func (n Node) Weight() float64 {
	return n.weight
}

// SetWeight of the node
func (n *Node) SetWeight(w float64) {
	n.weight = w
}

// Tolerance spcifies the possible tollerance from the weight.
func (n Node) Tolerance() float64 {
	return n.tolerance
}

// SetTolerance of the Node
func (n *Node) SetTolerance(t float64) {
	n.tolerance = t
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
