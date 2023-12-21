package merkle

import "crypto/sha256"

// Node contains data, representing the hash of
// the file or the combined hash of its children,
// and pointers to its left and right children.
type Node struct {
	Data  []byte
	Left  *Node
	Right *Node
}

// NewNode takes two nodes as input, hash their data together,
// and returns a new node with the resulting hash.
func NewNode(left, right *Node, data []byte) (n *Node) {
	n = &Node{}

	h := sha256.New()
	if left == nil && right == nil {
		// this is a leaf node
		h.Write(data)
		n.Data = h.Sum(nil)
	} else if right == nil {
		// this is a special case where we only have one node
		h.Write(left.Data)
		n.Data = h.Sum(nil)
	} else {
		// this is a non-leaf node
		h.Write(left.Data)
		h.Write(right.Data)
		n.Data = h.Sum(nil)
	}

	n.Left = left
	n.Right = right

	return
}
