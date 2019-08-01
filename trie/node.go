package trie

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"
)

// Node is the unit of trie
type Node interface {
	cache() (node hashNode, isDirty bool)
	fstring(string) string
}

type (
	branchNode struct {
		Children [17]Node
		flags    nodeFlag
	}
	shortNode struct {
		Key   []byte
		Value Node
		flags nodeFlag
	}
	valueNode []byte
	hashNode  []byte
)

func (n *branchNode) copy() *branchNode { copy := *n; return &copy }
func (n *shortNode) copy() *shortNode   { copy := *n; return &copy }

var nilValueNode = valueNode(nil)

// EncodeRLP encodes a full node into the consensus RLP format.
func (n *branchNode) EncodeRLP(w io.Writer) error {
	var nodes [17]Node

	for i, child := range &n.Children {
		if child != nil {
			nodes[i] = child
		} else {
			nodes[i] = nilValueNode
		}
	}
	return rlp.Encode(w, nodes)
}

var indices = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "[17]"}

func (n *branchNode) String() string { return n.fstring("") }
func (n *shortNode) String() string  { return n.fstring("") }
func (n hashNode) String() string    { return n.fstring("") }
func (n valueNode) String() string   { return n.fstring("") }

func (n *branchNode) fstring(ind string) string {
	resp := fmt.Sprintf("[\n%s  ", ind)
	for i, node := range &n.Children {
		if node == nil {
			resp += fmt.Sprintf("%s: <nil> ", indices[i])
		} else {
			resp += fmt.Sprintf("%s: %v", indices[i], node.fstring(ind+"  "))
		}
	}
	return resp + fmt.Sprintf("\n%s] ", ind)
}
func (n *shortNode) fstring(ind string) string {
	return fmt.Sprintf("{%x: %v} ", n.Key, n.Value.fstring(ind+"  "))
}
func (n hashNode) fstring(ind string) string {
	return fmt.Sprintf("<%x> ", []byte(n))
}
func (n valueNode) fstring(ind string) string {
	return fmt.Sprintf("%x ", []byte(n))
}

// flags contains caching-related metadata about a node.
type nodeFlag struct {
	hash  []byte // cached hash of the node (may be nil)
	dirty bool
}

func (n *branchNode) cache() (hashNode, bool) { return n.flags.hash, n.flags.dirty }
func (n *shortNode) cache() (hashNode, bool)  { return n.flags.hash, n.flags.dirty }
func (n hashNode) cache() (hashNode, bool)    { return nil, true }
func (n valueNode) cache() (hashNode, bool)   { return nil, true }

func mustDecodeNode(hash, buf []byte) Node {
	if s, _, err := rlp.SplitString(buf); err == nil && len(s) == 0 {
		// Handle nil node
		return nil
	}
	node, err := newNode(hash, buf)
	if err != nil {
		panic(fmt.Sprintf("node %x: %v", hash, err))
	}
	return node
}

// newNode parses the RLP encoding of a trie node.
func newNode(hash, buf []byte) (Node, error) {
	if len(buf) == 0 {
		return nil, errors.New("Unexpected end of buffer")
	}
	elements, _, err := rlp.SplitList(buf)
	if err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}
	switch count, _ := rlp.CountValues(elements); count {
	case 0:
		return valueNode(nil), nil
	case 2:
		node, err := newShortNode(hash, elements)
		return node, wrapError(err, "short")
	case 17:
		node, err := newBranchNode(hash, elements)
		return node, wrapError(err, "full")
	default:
		return nil, fmt.Errorf("Node elements count invalid: %v", count)
	}
}

func newShortNode(hash, elements []byte) (Node, error) {
	keyByte, rest, err := rlp.SplitString(elements)
	if err != nil {
		return nil, err
	}
	key := compactToHex(keyByte)
	flag := nodeFlag{hash: hash}

	// Leaf node
	if hasTerm(key) {
		value, _, err := rlp.SplitString(rest)
		if err != nil {
			return nil, fmt.Errorf("invalid value node: %v", err)
		}
		return &shortNode{
			Key:   key,
			Value: append(valueNode{}, value...),
			flags: flag,
		}, nil
	}

	// Extension node
	node, _, err := newRef(rest)
	if err != nil {
		return nil, wrapError(err, "val")
	}
	return &shortNode{
		Key:   key,
		Value: node,
		flags: flag,
	}, nil
}

func newBranchNode(hash, elements []byte) (*branchNode, error) {
	node := &branchNode{flags: nodeFlag{hash: hash}}
	for i := 0; i < 16; i++ {
		child, rest, err := newRef(elements)
		if err != nil {
			return node, wrapError(err, fmt.Sprintf("[%d]", i))
		}
		node.Children[i] = child
		elements = rest
	}
	value, _, err := rlp.SplitString(elements)
	if err != nil {
		return node, err
	}
	if len(value) > 0 {
		node.Children[16] = append(valueNode{}, value...)
	}
	return node, nil
}

const hashLen = len(Hash{})

func newRef(buf []byte) (Node, []byte, error) {
	kind, value, rest, err := rlp.Split(buf)
	if err != nil {
		return nil, buf, err
	}
	switch {
	case kind == rlp.List:
		// 'embedded' node reference. The encoding must be smaller
		// than a hash in order to be valid.
		size := len(buf) - len(rest)

		if size > hashLen {
			err := fmt.Errorf("oversized embedded node (size is %d bytes, want size < %d)", size, hashLen)
			return nil, buf, err
		}
		n, err := newNode(nil, buf)
		return n, rest, err
	case kind == rlp.String && len(value) == 0:
		// empty node
		return nil, rest, nil
	case kind == rlp.String && len(value) == 32:
		return append(hashNode{}, value...), rest, nil
	default:
		return nil, nil, fmt.Errorf("invalid RLP string size %d (want 0 or 32)", len(value))
	}
}

// wraps a decoding error with information about the path to the
// invalid child node (for debugging encoding issues).
type decodeError struct {
	what  error
	stack []string
}

func wrapError(err error, ctx string) error {
	if err == nil {
		return nil
	}
	if decErr, ok := err.(*decodeError); ok {
		decErr.stack = append(decErr.stack, ctx)
		return decErr
	}
	return &decodeError{err, []string{ctx}}
}

func (err *decodeError) Error() string {
	return fmt.Sprintf("%v (decode path: %s)", err.what, strings.Join(err.stack, "<-"))
}
