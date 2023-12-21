package merkle

func (t *Tree) ProofForBlock(block string) (proof []ProofHash) {
	blockHash := t.hashFn(block)

	var findProof func(node *Node) bool
	findProof = func(node *Node) bool {
		if node == nil {
			return false
		}

		if node.Data == blockHash {
			return true
		}

		if findProof(node.Left) {
			proof = append(proof, ProofHash{node.Right.Data, "L"})

			return true
		}

		if findProof(node.Right) {
			proof = append(proof, ProofHash{node.Left.Data, "R"})

			return true
		}

		return false
	}

	findProof(t.Root)

	return proof
}

func VerifyProof(rootHash string, block string, proof []ProofHash, hashFn hashFn) bool {
	currentHash := hashFn(block)

	for _, p := range proof {
		if p.Position == "L" {
			currentHash = hashFn(currentHash + p.Hash)
		} else if p.Position == "R" {
			currentHash = hashFn(p.Hash + currentHash)
		}
	}

	return currentHash == rootHash
}
