package merkle

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"merkle-file-uploader/pkg/utils"
)

var h = utils.Sha256

func TestMerkleTree(t *testing.T) {
	emptyTree, err := NewMerkleTree([]string{}, nil)
	assert.ErrorIs(t, err, ErrEmptyTreeInput)
	assert.Nil(t, emptyTree)

	blocks := []string{"A", "B", "C", "D", "E"}

	tree, err := NewMerkleTree(blocks, h)
	assert.NoError(t, err)
	assert.Equal(t, h("A"), tree.Root.Left.Left.Left.Data)
	assert.Equal(t, h("B"), tree.Root.Left.Left.Right.Data)
	assert.Equal(t, h("C"), tree.Root.Left.Right.Left.Data)
	assert.Equal(t, h("D"), tree.Root.Left.Right.Right.Data)
	assert.Equal(t, h("E"), tree.Root.Right.Left.Left.Data)
	assert.Equal(t, h("E"), tree.Root.Right.Left.Right.Data)
	assert.Equal(t, h("E"), tree.Root.Right.Right.Left.Data)
	assert.Equal(t, h("E"), tree.Root.Right.Right.Right.Data)

	assert.Equal(t, h(h("A")+h("B")), tree.Root.Left.Left.Data)
	assert.Equal(t, h(h("C")+h("D")), tree.Root.Left.Right.Data)
	assert.Equal(t, h(h("E")+h("E")), tree.Root.Right.Left.Data)
	assert.Equal(t, h(h("E")+h("E")), tree.Root.Right.Right.Data)

	assert.Equal(t, h(h(h("A")+h("B"))+h(h("C")+h("D"))), tree.Root.Left.Data)
	assert.Equal(t, h(h(h("E")+h("E"))+h(h("E")+h("E"))), tree.Root.Right.Data)
}

func TestGenerateProof(t *testing.T) {
	blocks := []string{"A", "B", "C", "D", "E"}

	tree, err := NewMerkleTree(blocks, h)
	assert.NoError(t, err)

	for _, b := range blocks {
		proof := tree.ProofForBlock(b)
		assert.True(t, VerifyProof(tree.Root.Data, b, proof, h))
	}
}
