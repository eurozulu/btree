package btree

import (
	"cmp"
	"io"
)

var TreeDegree int = 5
var minimumKeyCount = (TreeDegree - 1) / 2
var minimumSplitCount = (minimumKeyCount * 2) + 1

// BTree implements a simple btree of K keys, mapped to the V values.
type BTree[K cmp.Ordered, V any] interface {
	Put(key K, value V)
	Get(key K) (V, bool)
	Remove(key K) bool
	Depth() int
	Iterate() TreeIterator[K, V]
}

// btree wraps a single root node, which 'grows' into new root nodes as its child nodes split.
type btree[K cmp.Ordered, V any] struct {
	rootNode RootNode[K, V]
}

func (b *btree[K, V]) Put(key K, value V) {
	if root := b.rootNode.Insert(key, value); root != nil {
		b.rootNode = root
	}
}

func (b *btree[K, V]) Get(key K) (V, bool) {
	return b.rootNode.Value(key)
}

func (b *btree[K, V]) Remove(key K) bool {
	root, ok := b.rootNode.Remove(key)
	if !ok {
		return false
	}
	if root != nil {
		b.rootNode = root
	}
	return true
}

func (b *btree[K, V]) Depth() int {
	nlr := &nodeLevelReader[K, V]{root: b.rootNode}
	return nlr.Depth()
}

func (b *btree[K, V]) Iterate() TreeIterator[K, V] {
	return &treeIterator[K, V]{
		nodes:     nil,
		linkEntry: nil,
	}
	panic("implement me")
}

func ShowTree[K cmp.Ordered, V any](tree BTree[K, V], out io.Writer) error {
	//root := tree.(*btree[K, V]).rootNode
	//walk := model.NewTreeWalker(root)
	//var nodes [][]string
	//for walk.HasNext() {
	//	next := walk.Next()
	//	node := fmt.Sprintf("--%v--", next.Entries())
	//	depth := next.Depth()
	//	for len(nodes) < depth {
	//		nodes = append(nodes, []string{})
	//	}
	//	nodes[depth-1] = append(nodes[depth-1], node)
	//}
	//for _, row := range nodes {
	//	fmt.Fprintf(out, "%v", row)
	//	fmt.Fprintln(out)
	//}
	return nil
}

func NewTree[K cmp.Ordered, V any]() BTree[K, V] {
	return &btree[K, V]{
		rootNode: NewRootNode[K, V](),
	}
}
