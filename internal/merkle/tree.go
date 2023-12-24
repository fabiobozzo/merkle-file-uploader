package merkle

import (
	"errors"
)

var (
	ErrEmptyTreeInput = errors.New("cannot build a merkle tree from empty data")
)

type HashFn func(string) string

// Tree contains the root node of a merkle tree
type Tree struct {
	Root *Node
	HashFn
}

// NewTree creates a new Merkle tree from a slice of blocks using a given hash function.
// It returns a pointer to the new tree and any error encountered.
// The resulting Merkle tree is binary and balanced, with each leaf node containing one of the input blocks.
func NewTree(blocks []string, hashFn HashFn) (tree *Tree, err error) {
	tree = &Tree{HashFn: hashFn}

	if len(blocks) == 0 {
		return nil, ErrEmptyTreeInput
	}

	var nodes []Node

	// Create a leaf node for each block
	for _, block := range blocks {
		nodes = append(nodes, *NewNode(nil, nil, block, hashFn))
	}

	// Repeatedly combine pairs of nodes to create a new level in the tree,
	// until there is only one node left, which is the root of the Merkle tree.
	// This process ensures that the tree is binary (each non-leaf node has two children)
	// and balanced (all leaf nodes are at the same depth).
	for len(nodes) > 1 {
		if len(nodes)%2 == 1 {
			nodes = append(nodes, nodes[len(nodes)-1])
		}

		// Combine pairs of nodes to create the next level of the tree
		var level []Node
		for i := 0; i < len(nodes); i += 2 {
			level = append(level, *NewNode(&nodes[i], &nodes[i+1], "", hashFn))
		}

		// Replace the current level with the next level
		nodes = level
	}

	// The last remaining node is the root of the tree
	tree.Root = &nodes[0]

	return
}
