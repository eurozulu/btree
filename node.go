package btree

import (
	"cmp"
	"fmt"
	"log"
)

// NodeEntry represents the container for each Key entry in the Node.
type NodeEntry[K cmp.Ordered, V any] struct {
	Key   K
	Value *V
}

// node is the container of ordered NodeEntries.
// Node must have at least one entry.
// node may have child nodes.  When present, number of child nodes must be equal to the
// number of entries, plus one.
// When no child nodes present, node is known as a leaf node. #IsLeaf returns true.
type node[K cmp.Ordered, V any] struct {
	Entries  []NodeEntry[K, V]
	Children []node[K, V]
}

func (n node[K, V]) IsLeaf() bool {
	return len(n.Children) == 0
}

// Get returns the NodeEntry for the given key if it is present in the node or its children.
// If the key is not found, nil is returned.
func (n *node[K, V]) Get(key K) *NodeEntry[K, V] {
	i, e := n.KeyIndex(key)
	if e != nil {
		return e
	}
	if n.IsLeaf() {
		return nil
	}
	if i < 0 {
		i = len(n.Entries)
	}
	return n.Children[i].Get(key)
}

// Insert the given key/value pair into this leaf node.
// If node is not a leaf node panics.
func (n *node[K, V]) Insert(key K, value *V) {
	if !n.IsLeaf() {
		log.Panicf("Can not insert %v into a non leaf node", key)
	}
	i, e := n.KeyIndex(key)
	if i < 0 {
		// no existing key > new key, append to the end
		i = len(n.Entries)
	}
	if e == nil {
		n.Entries = InsertAtIndex(NodeEntry[K, V]{}, n.Entries, i)
		e = &n.Entries[i]
	}
	e.Key = key
	e.Value = value
}

// Delete the given key from this node
func (n *node[K, V]) Delete(key K) error {
	i, e := n.KeyIndex(key)
	if e == nil {
		return fmt.Errorf("key %v is unknown", key)
	}
	n.Entries = RemoveAtIndex(n.Entries, i)
	return nil
}

// Split this node into two child nodes with the median entry a single entry parent node.
func (n *node[K, V]) Split() *node[K, V] {
	l := len(n.Entries)
	if l < 3 {
		log.Fatalf("node too small to split. onlt %d entries found", len(n.Entries))
	}
	m := l / 2
	child1 := &node[K, V]{Entries: n.Entries[:m]}
	child2 := &node[K, V]{Entries: n.Entries[m+1:]}
	if !n.IsLeaf() {
		child1.Children = append(child1.Children, n.Children[:m+1]...)
		child2.Children = append(child2.Children, n.Children[m+1:]...)
	}
	return &node[K, V]{
		Entries:  []NodeEntry[K, V]{n.Entries[m]},
		Children: []node[K, V]{*child1, *child2},
	}
}

func (n node[K, V]) LastEntry() *NodeEntry[K, V] {
	if len(n.Entries) == 0 {
		return nil
	}
	return &n.Entries[len(n.Entries)-1]
}

func (n node[K, V]) LastChild() *node[K, V] {
	if n.IsLeaf() {
		return nil
	}
	return &n.Children[len(n.Children)-1]
}

func (n node[K, V]) String() string {
	if len(n.Children) > 0 {
		return fmt.Sprintf("{Entries: %v, Children: %v}", n.Entries, n.Children)
	}
	return fmt.Sprintf("{Entries: %v, Leaf}", n.Entries)
}

func (n *node[K, V]) mergeChild(childIndex int) int {
	entryIndex := childIndex
	if childIndex > 0 {
		// not the first child, perform backmerge
		n.backMergeChild(childIndex)
		entryIndex--
	} else {
		// first child, perform forwardMerge
		n.forwardMergeChild(childIndex)
	}
	// remove entry, now merged into child and also remove now empty child.
	n.Entries = RemoveAtIndex(n.Entries, entryIndex)
	n.Children = RemoveAtIndex(n.Children, childIndex)
	return entryIndex
}

func (n *node[K, V]) forwardMergeChild(childIndex int) {
	peer := &n.Children[childIndex+1]
	child := &n.Children[childIndex]
	// Add parent entry to end of child entries before added peers entries
	// (child.Entries may be empty)
	entries := append(child.Entries, n.Entries[childIndex])
	peer.Entries = append(entries, peer.Entries...)
	if len(child.Children) > 0 {
		peer.Children = append(child.Children, peer.Children...)
	}
}

func (n *node[K, V]) backMergeChild(childIndex int) {
	peer := &n.Children[childIndex-1]
	child := &n.Children[childIndex]
	// add parent entry to end of (back) peer before adding child entries
	entries := append(peer.Entries, n.Entries[childIndex-1])
	peer.Entries = append(entries, child.Entries...)
	if len(child.Children) > 0 {
		peer.Children = append(peer.Children, child.Children...)
	}
}

func (n *node[K, V]) getPreceeedingNode(nd *node[K, V]) *node[K, V] {
	for !nd.IsLeaf() {
		nd = nd.LastChild()
	}
	return nd
}

// keyIndex searches the nodes Entries for a matching key.
// If the key is found, the index in the Entries slice and the Entry iteself are returned.
// If the key is not found, but a key in this node is greater than the given key, tha index of the larger key is returned with a nil NodeEntry.
// If the given key is not in the Entries AND greater than all those keys, -1 and nil are returned.
func (n *node[K, V]) KeyIndex(key K) (int, *NodeEntry[K, V]) {
	for i, entry := range n.Entries {
		c := cmp.Compare(key, entry.Key)
		if c == 0 {
			return i, &entry
		}
		if c < 0 {
			return i, nil
		}
	}
	return -1, nil
}
