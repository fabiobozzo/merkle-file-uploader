# Merkle File Uploader
This project implements a file verification system where a client uploads a set of (small text) files to a server, deletes its local copies, and later downloads an arbitrary file from the server and verifies that the file is correct and has not been tampered with. This is achieved using a Merkle tree.

## Overview
The client computes a single Merkle tree root hash for the set of files and persists it on disk after uploading the files to the server and deleting its local copies. The client can request a specific file, by index, and a Merkle proof for it from the server. The client uses the proof to compute a root hash and compares it with the persisted root hash. If they match, the file is correct and can be displayed or stored locally.

The server stores the files and the Merkle tree, and provides an interface for uploading files, downloading files, and generating Merkle proofs.

## Implementation
The project is implemented in Go and uses the standard library for networking, allowing it to be deployed across multiple machines. The Merkle tree is implemented from scratch, with the standard library's `crypto/sha256` package used for the underlying hash function.

The project is structured into three main components:
- Client, which handles file uploading, downloading, and Merkle proof verification.
- Server, which handles file storage, and Merkle proof generation.
- MerkleTree, which provides methods for constructing the tree and generating proofs.

## Usage
The project is containerized using Docker and can be deployed using Docker Compose. To start the system, just run:
```
make start-server
```
This will start the server, listening on port `8080`.

Then, the client can be used to interact with the server.  
Two example calls are available at: 
```
make test-upload    # mfu client upload resources/ 
make test-download  # mfu client download 1
```

Please, feel free to explore with the `mfu` CLI yourself, by uploading other files (list of varargs) or folders, and download index (incl. invalid or non-existent indexes).

Finally, tear down the server environment via:
```
make stop-server
```

## Design Choices
- The dependency tree is kept to the bare minimum:
  - `spf13/cobra` to build a CLI tool
  - `gorilla/mux` to facilitate the REST paths handling
  - `stretchr/testify` for unit test assertions
  - `aws/aws-sdk-go-v2` as S3 client


- There are abstractions in place to prepare the ground for future developments: 
  - The upload/download protocol has a HTTP implementation. My next step would be to implement a gRPC-protobuf -based protocol (mainly to leverage streams, because sending files one by one over HTTP would not scale well IRL)
  - The `Storage` interface used by the server has two basic implementations:
    - a naive, in-memory one
    - a more realistic, S3 bucket

## Limitations and future improvements
⚠️ **Disclaimer** !  
This project is a working Proof-of-Concept, for the sake of showing a remote Merkle Proof verification systems. There are several areas where it could be further developed and prepared to be production-ready:

1. **Coverage**: I wrote the (happy flow) unit tests for the Merkle tree and its proof generation and verification. That's the juicy part. For the sake of full coverage, though, the boilerplate testing of the http-based protocol (mocking `Storage`) and utility functions should be added, too. On top of that, comprehensive integration and performance tests. 
2. **Workflow**: The client-server is single-shot. Only one batch of files is managed, at a time. It would be more useful to create separate batches upload, and reference them as separate entities. Further down the road, I'd like to explore how to add/remove to an existing set of files. 
3. **Server Storage**: The uploaded files are stored based on their index, and retrieved sequentially, as needed. The merkle tree is not stored and it's rebuilt when a merkle proof must be generated. The 1st improvement would be to serialize and store the merkle tree along with the (leaf) raw files it represents. The 2nd improvement would be to organize the files in a better way (e.g. keeping the indexes in Redis for fast location lookups).
4. **Synchronization and Concurrency**: The server does not currently handle concurrent requests, which could lead to inconsistencies in the Merkle tree. A future improvement could be to add locking or use a concurrent data structure for the Merkle tree, or even relying on transactions. 
5. **Performance**: I consider the time/space complexity of the Merkle proof generation good enough for this use case, although it would be interesting to increase the algorithm and space complexity a bit, and try to speed it up by concurrently searching the left and right subtrees in parallel. 

