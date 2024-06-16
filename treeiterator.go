package btree

import (
	"cmp"
	"github.com/eurozulu/btree/utils"
)

type TreeIterator[K cmp.Ordered, V any] interface {
	HasNext() bool
	Next() []Entry[K, V]
	Depth() int
}

type treeIterator[K cmp.Ordered, V any] struct {
	next      []Entry[K, V]
	nodes     utils.StackSlice[*node[K, V]]
	linkEntry Entry[K, V]
}

func (it *treeIterator[K, V]) HasNext() bool {
	if it.next == nil {
		it.next = it.getNext()
	}
	return it.next != nil
}

func (it *treeIterator[K, V]) Next() []Entry[K, V] {
	if !it.HasNext() {
		return nil
	}
	next := it.next
	it.next = nil
	return next
}

func (it treeIterator[K, V]) Depth() int {
	d := len(it.nodes)
	if it.linkEntry != nil {
		d--
	}
	return d
}

func (it *treeIterator[K, V]) getNext() []Entry[K, V] {
	if it.linkEntry != nil {
		le := []Entry[K, V]{it.linkEntry}
		it.linkEntry = nil
		return le
	}
	// ensure current node is a leaf
	n, ok := it.nodes.Peek()
	if !ok {
		return nil
	}
	if !n.IsLeaf() {
		it.skipToFirstLeaf(n)
	}
	leaf, ok := it.nodes.Pop()
	if !ok {
		return nil
	}
	it.linkEntry = it.skipToNextNode(leaf)
	return leaf.Entries()
}

// skipToNextNode positions the nodes stack on the next available leaf node.
// and returns the sperator entry dividing the two sibling nodes.
// If the current leaf is the last sibling in the parent, recursively calls itself with the parent node
// If no more nodes are available, nil is returned.
func (it *treeIterator[K, V]) skipToNextNode(child *node[K, V]) Entry[K, V] {
	parent, ok := it.nodes.Peek()
	if !ok {
		return nil
	}
	nextindex := utils.IndexOf(child, parent.children) + 1
	if nextindex == 0 {
		// not a child of given parent!!!
		return nil
	}
	if nextindex >= len(parent.children) {
		// no more siblings in given parent, recursive call with parent as 'child'
		parent, _ = it.nodes.Pop()
		return it.skipToNextNode(parent)
	}
	it.nodes.Push(parent.children[nextindex])
	return parent.entries[nextindex-1]
}

func (it *treeIterator[K, V]) skipToFirstLeaf(n *node[K, V]) *node[K, V] {
	for {
		if n.IsLeaf() {
			break
		}
		n = n.children[0]
		it.nodes.Push(n)
	}
	return n
}

func newTreeIterator[K cmp.Ordered, V any](root *node[K, V]) TreeIterator[K, V] {
	nodes := utils.StackSlice[*node[K, V]]{}
	if root != nil && len(root.entries) > 0 {
		nodes.Push(root)
	}

	it := &treeIterator[K, V]{nodes: nodes}
	if root != nil {
		it.skipToFirstLeaf(root)
	}
	return it
}
