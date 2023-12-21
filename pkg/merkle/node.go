package merkle

// Node contains data, representing the h of
// the file or the combined h of its children,
// and pointers to its left and right children.
type Node struct {
	Data  string
	Left  *Node
	Right *Node
}

// NewNode takes two nodes as input, h their data together,
// and returns a new node with the resulting h.
func NewNode(left, right *Node, data string, hash hashFn) (n *Node) {
	n = &Node{}

	if left == nil && right == nil {
		// this is a leaf node
		n.Data = hash(data)
	} else if right == nil {
		// this is a special case where we only have one node
		n.Data = hash(left.Data)
	} else {
		n.Data = hash(left.Data + right.Data)
	}

	n.Left = left
	n.Right = right

	return
}
