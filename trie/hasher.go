package trie

import (
	"hash"
	"sync"

	"github.com/QuoineFinancial/vertex/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

type hasher struct {
	tmp sliceBuffer
	sha keccakState
}

// keccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type keccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

type sliceBuffer []byte

func (b *sliceBuffer) Write(data []byte) (n int, err error) {
	*b = append(*b, data...)
	return len(data), nil
}

func (b *sliceBuffer) Reset() { *b = (*b)[:0] }

// hashers live in a global db.
var hasherPool = sync.Pool{
	New: func() interface{} {
		return &hasher{
			tmp: make(sliceBuffer, 0, 550), // cap is as large as a full fullNode.
			sha: sha3.NewLegacyKeccak256().(keccakState),
		}
	},
}

func newHasher() *hasher { return hasherPool.Get().(*hasher) }

func returnHasherToPool(h *hasher) { hasherPool.Put(h) }

func (h *hasher) hash(n Node, db *db.RocksDB, force bool) (Node, Node, error) {
	if hash, dirty := n.cache(); hash != nil {
		if db == nil {
			return hash, n, nil
		}
		if !dirty {
			switch n.(type) {
			case *branchNode, *shortNode:
				return hash, hash, nil
			default:
				return hash, n, nil
			}
		}
	}

	// Trie not processed yet or needs storage, walk the children
	collapsed, cached, err := h.hashChildren(n, db)
	if err != nil {
		return hashNode{}, n, err
	}
	hashed, err := h.store(collapsed, db, force)
	if err != nil {
		return hashNode{}, n, err
	}

	// Cache the hash of the node for later reuse and remove
	// the dirty flag in commit mode. It's fine to assign these values directly
	// without copying the node first because hashChildren copies it.
	cachedHash, _ := hashed.(hashNode)
	switch cachedNode := cached.(type) {
	case *shortNode:
		cachedNode.flags.hash = cachedHash
		if db != nil {
			cachedNode.flags.dirty = false
		}
	case *branchNode:
		cachedNode.flags.hash = cachedHash
		if db != nil {
			cachedNode.flags.dirty = false
		}
	}
	return hashed, cached, nil
}

// hashChildren replaces the children of a node with their hashes if the encoded
// size of the child is larger than a hash, returning the collapsed node as well
// as a replacement for the original node with the child hashes cached in.
func (h *hasher) hashChildren(original Node, db *db.RocksDB) (Node, Node, error) {
	var err error

	switch node := original.(type) {
	case *shortNode:
		// Hash the short node's child, caching the newly hashed subtree
		collapsed, cached := node.copy(), node.copy()
		collapsed.Key = hexToCompact(node.Key)
		cached.Key = common.CopyBytes(node.Key)

		if _, ok := node.Value.(valueNode); !ok {
			collapsed.Value, cached.Value, err = h.hash(node.Value, db, false)
			if err != nil {
				return original, original, err
			}
		}

		return collapsed, cached, nil

	case *branchNode:
		// Hash the full node's children, caching the newly hashed subtrees
		collapsed, cached := node.copy(), node.copy()

		for i := 0; i < 16; i++ {
			if node.Children[i] != nil {
				collapsed.Children[i], cached.Children[i], err = h.hash(node.Children[i], db, false)
				if err != nil {
					return original, original, err
				}
			}
		}
		cached.Children[16] = node.Children[16]
		return collapsed, cached, nil

	default:
		// Value and hash nodes don't have children so they're left as were
		return node, original, nil
	}
}

// store hashes the node n and if we have a storage layer specified, it writes
// the key/value pair to it and tracks any node->child references as well as any
// node->external trie references.
func (h *hasher) store(node Node, db *db.RocksDB, force bool) (Node, error) {

	// Don't store hashes or empty nodes.
	if _, isHash := node.(hashNode); node == nil || isHash {
		return node, nil
	}

	// Generate the RLP encoding of the node
	h.tmp.Reset()
	if err := rlp.Encode(&h.tmp, node); err != nil {
		panic("encode error: " + err.Error())
	}
	if len(h.tmp) < 32 && !force {
		return node, nil // Nodes smaller than 32 bytes are stored inside their parent
	}

	// Larger nodes are replaced by their hash and stored in the database.
	hash, _ := node.cache()
	if hash == nil {
		hash = h.makeHashNode(h.tmp)
	}

	// Store node
	if db != nil {
		hash := common.BytesToHash(hash)
		blob, err := rlp.EncodeToBytes(node)
		if err != nil {
			return nil, err
		}
		db.Put(hash[:], blob)
	}

	return hash, nil
}

func (h *hasher) makeHashNode(data []byte) hashNode {
	n := make(hashNode, h.sha.Size())
	h.sha.Reset()
	h.sha.Write(data)
	h.sha.Read(n)
	return n
}
