package merkle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerkleProof(t *testing.T) {
	cases := map[string]struct {
		blocks  []string
		wantErr error
	}{
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
			tree, err := NewTree(tc.blocks, h)
			assert.NoError(t, err)

			for _, b := range tc.blocks {
				proof := tree.ProofForBlock(b)
				assert.True(t, VerifyProof(tree.Root.Data, b, proof, h))
			}

			proof := tree.ProofForBlock("X")
			assert.False(t, VerifyProof(tree.Root.Data, "X", proof, h))
		})
	}
}
