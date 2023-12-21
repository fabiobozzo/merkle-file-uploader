package merkle

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"merkle-file-uploader/pkg/utils"
)

var h = utils.Sha256

func TestMerkleTree(t *testing.T) {
	var traverseTree func(*Node, int, map[string]struct{}) int
	traverseTree = func(node *Node, depth int, blockHashes map[string]struct{}) int {
		// Leaf nodes must be among the hashes of the blocks
		if node.Left == nil && node.Right == nil {
			if _, ok := blockHashes[node.Data]; !ok {
				t.Errorf("leaf node hash is not a block hash: %s", node.Data)
			}

			return depth
		}

		leftDepth := traverseTree(node.Left, depth+1, blockHashes)
		rightDepth := traverseTree(node.Right, depth+1, blockHashes)
		if leftDepth != rightDepth {
			t.Error("the tree is not balanced")
		}

		// Non-leaf nodes must be the hashes of their children's hashes
		expectedHash := h(node.Left.Data + node.Right.Data)
		if node.Data != expectedHash {
			t.Errorf("node hash does not match children's hashes: got %s, want %s", node.Data, expectedHash)
		}

		return leftDepth
	}

	cases := map[string]struct {
		blocks  []string
		wantErr error
	}{
		"empty tree": {
			blocks:  []string{},
			wantErr: ErrEmptyTreeInput,
		},
		"even-sized tree": {
			blocks: []string{"A", "B", "C", "D", "E", "F"},
		},
		"odd-sized tree": {
			blocks: []string{"A", "B", "C", "D", "E"},
		},
		"duplicates": {
			blocks: []string{"A", "B", "B", "A"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			emptyTree, err := NewMerkleTree([]string{}, nil)
			assert.ErrorIs(t, err, ErrEmptyTreeInput)
			assert.Nil(t, emptyTree)

			blockHashes := map[string]struct{}{}
			for _, block := range tc.blocks {
				blockHashes[h(block)] = struct{}{}
			}

			tree, err := NewMerkleTree(tc.blocks, h)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)

				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, tree.Root)

			traverseTree(tree.Root, 0, blockHashes)
		})
	}

}
