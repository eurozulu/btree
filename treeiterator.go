package btree

import (
	"cmp"
)

type TreeIterator[K cmp.Ordered, V any] interface {
	HasNext() bool
	Next() []NodeEntry[K, V]
	Depth() int
}

type treeIterator[K cmp.Ordered, V any] struct {
	next      []NodeEntry[K, V]
	nodes     StackSlice[*node[K, V]]
	linkEntry *NodeEntry[K, V]
}

func (it *treeIterator[K, V]) HasNext() bool {
	if it.next == nil {
		it.next = it.getNext()
	}
	return it.next != nil
}

func (it *treeIterator[K, V]) Next() []NodeEntry[K, V] {
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

func (it *treeIterator[K, V]) getNext() []NodeEntry[K, V] {
	if it.linkEntry != nil {
		le := []NodeEntry[K, V]{*it.linkEntry}
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
	return leaf.Entries
}

// skipToNextNode positions the nodes stack on the next available leaf node.
// and returns the seperator entry dividing the two sibling nodes.
// If the current leaf is the last sibling in the parent, recursively calls itself with the parent node
// If no more nodes are available, nil is returned.
func (it *treeIterator[K, V]) skipToNextNode(child *node[K, V]) *NodeEntry[K, V] {
	parent, ok := it.nodes.Peek()
	if !ok {
		return nil
	}
	nextindex := indexOfChild(parent, child) + 1
	if nextindex == 0 {
		// not a child of given parent!!!
		return nil
	}
	if nextindex >= len(parent.Children) {
		// no more siblings in given parent, recursive call with parent as 'child'
		parent, _ = it.nodes.Pop()
		return it.skipToNextNode(parent)
	}
	it.nodes.Push(&parent.Children[nextindex])
	return &parent.Entries[nextindex-1]
}

func (it *treeIterator[K, V]) skipToFirstLeaf(n *node[K, V]) *node[K, V] {
	for {
		if n.IsLeaf() {
			break
		}
		n = &n.Children[0]
		it.nodes.Push(n)
	}
	return n
}

func indexOfChild[K cmp.Ordered, V any](parent, child *node[K, V]) int {
	for i := range parent.Children {
		if &parent.Children[i] == child {
			return i
		}
	}
	return -1
}

func newTreeIterator[K cmp.Ordered, V any](rootNode *node[K, V]) TreeIterator[K, V] {
	it := &treeIterator[K, V]{
		nodes: StackSlice[*node[K, V]]{},
	}
	if rootNode != nil && len(rootNode.Entries) > 0 {
		it.nodes.Push(rootNode)
		it.skipToFirstLeaf(rootNode)
	}
	return it
}
