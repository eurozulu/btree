package btree

import (
	"cmp"
	"context"
	"log"
)

type BTree[K cmp.Ordered, V any] interface {
	Degree() int
	Depth() int
	IsEmpty() bool
	Keys(ctx context.Context) <-chan K
	Get(key K) *V
	Add(key K, value *V) error
	Remove(key K) error
}

type bTree[K cmp.Ordered, V any] struct {
	rootnode *node[K, V]
	degree   int
}

func (b bTree[K, V]) Degree() int {
	return b.degree
}

func (b bTree[K, V]) Depth() int {
	d := 0
	n := b.rootnode
	for !n.IsLeaf() {
		d++
		n = &n.Children[0]
	}
	return d
}

func (b bTree[K, V]) Count() int {
	count := 0
	it := newTreeIterator(b.rootnode)
	for it.HasNext() {
		nodes := it.Next()
		count += len(nodes)
	}
	return count
}

func (b bTree[K, V]) Keys(ctx context.Context) <-chan K {
	ch := make(chan K)
	go func(ch chan<- K) {
		defer close(ch)
		it := newTreeIterator(b.rootnode)
		for it.HasNext() {
			nodes := it.Next()
			for _, n := range nodes {
				select {
				case <-ctx.Done():
					return
				case ch <- n.Key:
				}
			}
		}
	}(ch)
	return ch
}

func (b bTree[K, V]) IsEmpty() bool {
	return b.rootnode.IsEmpty()
}

func (b bTree[K, V]) Get(key K) *V {
	ne := b.rootnode.Get(key)
	if ne == nil {
		return nil
	}
	return ne.Value
}

func (b *bTree[K, V]) Add(key K, value *V) error {
	nn := b.add(key, value, b.rootnode)
	if nn != nil {
		// a new root pushed up
		b.rootnode = nn
	}
	return nil
}

func (b *bTree[K, V]) Remove(key K) error {
	nn, err := b.remove(key, b.rootnode)
	if err != nil {
		return err
	}
	if nn != nil {
		b.rootnode = nn
	}
	return nil
}

func (b *bTree[K, V]) add(key K, value *V, nd *node[K, V]) *node[K, V] {
	if nd.IsLeaf() {
		nd.Insert(key, value)
	} else {
		nd = b.addToChild(key, value, nd)
	}
	if nd == nil || len(nd.Entries) < b.degree {
		// node size within bounds, all done
		return nil
	}
	return nd.Split()
}

func (b *bTree[K, V]) addToChild(key K, value *V, nd *node[K, V]) *node[K, V] {
	i, e := nd.keyIndex(key)
	if e != nil {
		// already exists, update value
		e.Value = value
		return nd
	}
	if i < 0 {
		// new key greater than existing keys, use last child.
		i = len(nd.Entries)
	}
	nn := b.add(key, value, &nd.Children[i])
	if nn == nil {
		// no new node due to split, all done
		return nil
	}
	// Child has split, merge nn into parent (nd) node
	nd.Entries = InsertAtIndex(nn.Entries[0], nd.Entries, i)
	nd.Children[i] = nn.Children[1]
	nd.Children = InsertAtIndex(nn.Children[0], nd.Children, i)
	return nd
}

func (b *bTree[K, V]) remove(key K, nd *node[K, V]) (*node[K, V], error) {
	if nd.IsLeaf() {
		// leaf node simply deletes key and lets parent node balance entries. (Except root node, with no parent)
		if err := nd.Delete(key); err != nil {
			return nil, err
		}
		return nil, nil //TODO review if return node required
	}
	// non leaf / parent node
	i, e := nd.keyIndex(key)
	if i < 0 {
		i = len(nd.Entries)
	}
	if e == nil {
		// not in this node, remove from child
		return b.removeFromChild(key, i, nd)
	}
	// a parent node containing key to remove

	// replace entry to be removed with the preceding entry, from the right most leaf of the child
	pn := nd.getPreceeedingNode(&nd.Children[i])
	lastE := pn.LastEntry()
	nd.Entries[i] = *lastE

	// perform a Remove of the copied entry, to remove from the leaf we stole it from and rebalnce tree
	return b.removeFromChild(lastE.Key, i, nd)
}

func (b *bTree[K, V]) removeFromChild(key K, childIndex int, nd *node[K, V]) (*node[K, V], error) {
	child := &nd.Children[childIndex]
	if _, err := b.remove(key, child); err != nil {
		return nil, err
	}
	if len(child.Entries) > 0 {
		// child still has enough entries
		return nil, nil
	}
	// Child now empty, Merge into one of its peers and include the entry from this node which "bridges" the merge childre,
	entryIndex := nd.mergeChild(childIndex)
	mergedChild := nd.Children[entryIndex]
	// Ensure merged child is not too big
	if len(mergedChild.Entries) >= b.degree {
		// merges node now too big, perform split
		nn := mergedChild.Split()
		nd.Entries = InsertAtIndex(nn.Entries[0], nd.Entries, entryIndex)
		nd.Children[entryIndex] = nn.Children[0]
		nd.Children = InsertAtIndex(nn.Children[1], nd.Children, childIndex)
	}

	if len(nd.Entries) == 0 {
		// If this parent now empty, pass up it's first child
		return &nd.Children[0], nil
	}
	return nil, nil
}

func NewBTree[K cmp.Ordered, V any](degree int) BTree[K, V] {
	if degree < 2 {
		log.Fatalf("degree must be >= 2")
	}

	return &bTree[K, V]{
		rootnode: &node[K, V]{},
		degree:   degree,
	}
}
