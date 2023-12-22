package protocol

import "merkle-file-uploader/internal/merkle"

type UploadedFile struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
}

type UploadedFilesResponse struct {
	UploadedFiles []UploadedFile `json:"uploadedFiles"`
}

type MerkleProofResponse struct {
	MerkleProof []merkle.ProofHash `json:"merkleProof"`
}
