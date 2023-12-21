package merkle

import (
	"bytes"
	"errors"
)

var (
	ErrEmptyTreeInput = errors.New("cannot build a merkle tree from empty data")
	ErrHashNotFound   = errors.New("hash not found in merkle tree")
)

// Tree contains the root node of a merkle tree
type Tree struct {
	Root *Node
}

func NewMerkleTree(blocks [][]byte) (tree *Tree, err error) {
	tree = &Tree{}

	if len(blocks) == 0 {
		return nil, ErrEmptyTreeInput
	}

	var nodes []Node

	// create a leaf node for each block
	for _, datum := range blocks {
		nodes = append(nodes, *NewNode(nil, nil, datum))
	}

	for len(nodes) > 1 {
		if len(nodes)%2 == 1 {
			nodes = append(nodes, nodes[len(nodes)-1])
		}

		// repeatedly combines pairs of nodes to create a new level in the tree
		var level []Node
		for i := 0; i < len(nodes); i += 2 {
			level = append(level, *NewNode(&nodes[i], &nodes[i+1], nil))
		}

		nodes = level
	}

	// until there is only one node left, which is the root of the Merkle tree
	tree.Root = &nodes[0]

	return
}

func (t *Tree) Find(hash []byte) (node *Node, path []*Node, err error) {
	found := t.findRecursive(t.Root, hash, &path)
	if found {
		return path[len(path)-1], path, nil
	}

	return nil, nil, ErrHashNotFound
}

func (t *Tree) findRecursive(node *Node, hash []byte, path *[]*Node) bool {
	if node == nil {
		return false
	}

	*path = append(*path, node)
	if bytes.Equal(node.Data, hash) {
		return true
	}

	if t.findRecursive(node.Left, hash, path) || t.findRecursive(node.Right, hash, path) {
		return true
	}
	*path = (*path)[:len(*path)-1] // backtrack

	return false
}

func (t *Tree) GenerateProof(path []*Node) [][]byte {
	var proof [][]byte
	// We start from the end of the path (the leaf node) and move up towards the root.
	for i := len(path) - 1; i > 0; i-- {
		node := path[i]     // The current node in the path
		parent := path[i-1] // The parent of the current node

		var sibling *Node
		// Determine the sibling of the current node.
		// If the current node is the left child of its parent, the sibling is the right child, and vice versa.
		if parent.Left == node {
			sibling = parent.Right
		} else {
			sibling = parent.Left
		}

		// If the sibling exists, add its data (hash) to the proof.
		// The sibling's data is part of the proof because it's necessary to compute the parent's data.
		if sibling != nil {
			proof = append(proof, sibling.Data)
		}
	}
	// The proof is a list of hashes that allows us to compute the root hash from the leaf node.
	return proof
}
