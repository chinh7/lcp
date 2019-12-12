package trie

import (
	"fmt"
)

// PopFirst return the first key of trie
func (tree *Trie) PopFirst() ([]byte, error) {
	value, newRoot, reachedHashNode, err := tree.getFirst(tree.root)
	if err == nil && reachedHashNode {
		tree.root = newRoot
	}
	return value, err
}

func (tree *Trie) getFirst(currentNode Node) (value []byte, newNode Node, reachedHashNode bool, err error) {
	switch node := currentNode.(type) {
	case *shortNode:
		value, newNode, reachedHashNode, err = tree.getFirst(node.Value)
		if err == nil && reachedHashNode {
			node = node.copy()
			node.Value = newNode
		}
		return value, node, reachedHashNode, err

	case *branchNode:
		for _, child := range node.Children {
			value, newNode, reachedHashNode, err = tree.getFirst(child)
			if err == nil && reachedHashNode {
				node = node.copy()
				child = newNode
			}
			return value, node, reachedHashNode, err
		}
		return nil, nil, false, nil
	case nil:
		return nil, nil, false, nil

	case valueNode:
		return node, node, false, nil

	case hashNode:
		loadedNode, err := tree.loadNode(node)
		if err != nil {
			return nil, node, true, err
		}
		value, newNode, _, err := tree.getFirst(loadedNode)
		return value, newNode, true, err

	default:
		panic(fmt.Sprintf("%T: invalid node: %v", currentNode, currentNode))
	}
}
