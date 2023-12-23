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

	// Repeatedly combines pairs of nodes to create a new level in the tree,
	// until there is only one node left, which is the root of the Merkle tree.
	for len(nodes) > 1 {
		if len(nodes)%2 == 1 {
			nodes = append(nodes, nodes[len(nodes)-1])
		}

		var level []Node
		for i := 0; i < len(nodes); i += 2 {
			level = append(level, *NewNode(&nodes[i], &nodes[i+1], "", hashFn))
		}

		nodes = level
	}

	tree.Root = &nodes[0]

	return
}
