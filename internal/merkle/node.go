package merkle

// Node contains data, representing the hash of
// the file or the combined hash of its children,
// and pointers to its left and right children.
type Node struct {
	Data  string
	Left  *Node
	Right *Node
}

// NewNode takes two nodes as input, hash their data together,
// and returns a new node with the resulting hash.
func NewNode(left, right *Node, data string, hash HashFn) (n *Node) {
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
