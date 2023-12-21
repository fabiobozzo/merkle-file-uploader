package merkle

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	emptyTree, err := NewMerkleTree([][]byte{})
	assert.ErrorIs(t, err, ErrEmptyTreeInput)
	assert.Nil(t, emptyTree)

	tree, err := NewMerkleTree([][]byte{
		[]byte("A"),
		[]byte("B"),
		[]byte("C"),
		[]byte("D"),
		[]byte("E"),
	})
	assert.NoError(t, err)

	for _, block := range []string{"A", "B", "C", "D", "E"} {
		hash := sha256.Sum256([]byte(block))
		found, path, err := tree.Find(hash[:])
		assert.NoError(t, err)
		assert.Equal(t, hash[:], found.Data)
		assert.True(t, len(path) >= 3)
		assert.Equal(t, hash[:], path[len(path)-1].Data)
		assert.Equal(t, tree.Root.Data, path[0].Data)
	}

	nonExistentHash := sha256.Sum256([]byte("F"))
	found, _, err := tree.Find(nonExistentHash[:])
	assert.ErrorIs(t, err, ErrHashNotFound)
	assert.Nil(t, found)
}
