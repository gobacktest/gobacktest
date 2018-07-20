package gobacktest

// Asset is a data building block for the tree structure, eg. Stock, Option, Cash etc.
// It implements the NodeHandler interface via the promoted Node field.
type Asset struct {
	Node
}

// NewAsset return a new strategy node ready to use.
func NewAsset(name string) *Asset {
	var asset = &Asset{}
	asset.SetName(name)
	asset.SetRoot(false)
	return asset
}

// Children returns an empty slice and false, an Asset is not allowed to have children.
func (a Asset) Children() ([]NodeHandler, bool) {
	return []NodeHandler{}, false
}

// SetChildren return itself without change, as an Asset ist not allowed to have children.
func (a *Asset) SetChildren(c ...NodeHandler) NodeHandler {
	return a
}
