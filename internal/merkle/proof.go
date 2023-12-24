package merkle

type ProofHash struct {
	Hash     string
	Position string
}

// ProofForBlock generates a Merkle proof for a given block.
// The proof consists of a slice of hashes from the Merkle tree.
func (t *Tree) ProofForBlock(block string) (proof []ProofHash) {
	// Compute the hash of the block with the same hash function used to build the tree
	blockHash := t.HashFn(block)

	// Define a recursive function to find the proof
	var findProof func(node *Node) bool
	findProof = func(node *Node) bool {
		// If the node is nil, we've hit a leaf node without finding the block
		if node == nil {
			return false
		}

		// If the node's data matches the block hash, we've found the block
		if node.Data == blockHash {
			return true
		}

		// Recursively search the left subtree
		if findProof(node.Left) {
			// If the block was found in the left subtree, add the right sibling to the proof
			proof = append(proof, ProofHash{node.Right.Data, "L"})

			return true
		}

		// Recursively search the right subtree
		if findProof(node.Right) {
			// If the block was found in the right subtree, add the left sibling to the proof
			proof = append(proof, ProofHash{node.Left.Data, "R"})

			return true
		}

		return false
	}

	// Start the search at the root of the tree
	findProof(t.Root)

	return proof
}

// VerifyProof verifies a Merkle proof for a given block and root hash.
// It returns true if the proof is valid, and false otherwise.
func VerifyProof(rootHash string, block string, proof []ProofHash, hashFn HashFn) bool {
	currentHash := hashFn(block)

	// Iterate over the proof hashes
	for _, p := range proof {
		// Depending on the position of the sibling in the tree,
		// concatenate it with the current hash and compute the new current hash
		if p.Position == "L" {
			currentHash = hashFn(currentHash + p.Hash)
		} else if p.Position == "R" {
			currentHash = hashFn(p.Hash + currentHash)
		}
	}

	// The proof is valid if the final computed hash matches the root hash
	return currentHash == rootHash
}
