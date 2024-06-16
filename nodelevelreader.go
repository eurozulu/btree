package btree

import (
	"cmp"
)

type NodeLevelReader[K cmp.Ordered, V any] interface {
	NodesAtDepth(depth int) []Node[K, V]
	Depth() int
}

type nodeLevelReader[K cmp.Ordered, V any] struct {
	root Node[K, V]
}

func (it nodeLevelReader[K, V]) Depth() int {
	d := 0
	n := it.root
	for !n.IsLeaf() {
		n = n.Children()[0]
		d++
	}
	return d
}

func (it *nodeLevelReader[K, V]) NodesAtDepth(depth int) []Node[K, V] {
	ns := []Node[K, V]{it.root}
	d := 0
	for d < depth {
		ns = it.childNodesOf(ns)
		d++
	}
	return ns
}

func (it *nodeLevelReader[K, V]) childNodesOf(n []Node[K, V]) []Node[K, V] {
	nds := make([]Node[K, V], 0, len(n))
	for _, parent := range n {
		nds = append(nds, parent.Children()...)
	}
	return nds
}

func newNodeIterator[K cmp.Ordered, V any](root *node[K, V]) NodeLevelReader[K, V] {
	return &nodeLevelReader[K, V]{root: root}
}
