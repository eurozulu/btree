package btree

import (
	"cmp"
	"fmt"
	"github.com/eurozulu/btree/utils"
	"log"
)

// Node represents a single entity within a tree.
// Nodes contains Key/Value pairs, known as entries, and, optionally two or more child nodes.
// Depending on the degree size of the tree (N), A Node MUST contain:
// at least (N - 1) / 2 Keys. i.e. degree with a of 4, (4 - 1) / 2 = 1. (rounded down)
// So a degree of 3 or 4, dictates a minimum of one key in a node.
// A node must also never contain more than N - 1 Keys.
type Node[K cmp.Ordered, V any] interface {
	// Keys returns the ordered list of keys in this node
	Keys() []K

	// Value returns the Value associated with the given key.
	// Value performs a recursive search into child Nodes when this node doesn't contain the key.
	Value(key K) (V, bool)

	// IsLeaf indicates if this node has children or not.  When true this node has no child nodes.
	IsLeaf() bool

	Entries() []Entry[K, V]

	// GetChild returns the child node for the given key.
	// The child is determined by the given keys position in this nodes Keys.
	// If the given key is smaller than the smallest key in this node, the first child is returned.
	// if the given key is larger than the smallest key, but smaller than other keys in this node, the second child is returned.
	// and so on unless the given key is larger than any key in this node, and so the last child is returned.
	// If the given key is in this nodes Keys, then nil is returned.
	GetChild(key K) Node[K, V]
	Children() []Node[K, V]
}

type RootNode[K cmp.Ordered, V any] interface {
	Node[K, V]
	Insert(key K, value V) RootNode[K, V]
	Remove(key K) (RootNode[K, V], bool)
}

type node[K cmp.Ordered, V any] struct {
	entries  []*entry[K, V]
	children []*node[K, V]
}

func (n node[K, V]) String() string {
	return fmt.Sprintf("[%v]", n.entries)
}

func (n node[K, V]) Keys() []K {
	if len(n.entries) == 0 {
		return nil
	}
	kz := make([]K, len(n.entries))
	for i, entry := range n.entries {
		kz[i] = entry.key
	}
	return kz
}

func (n node[K, V]) Value(key K) (V, bool) {
	i, ok := n.entryPosition(key)
	if ok {
		return n.entries[i].value, ok
	}
	if !n.IsLeaf() {
		return n.children[i].Value(key)
	}
	var v V
	return v, false
}

func (n node[K, V]) IsLeaf() bool {
	return len(n.children) == 0
}

func (n node[K, V]) Entries() []Entry[K, V] {
	es := make([]Entry[K, V], len(n.entries))
	for i, entry := range n.entries {
		es[i] = entry
	}
	return es
}

func (n node[K, V]) Children() []Node[K, V] {
	nds := make([]Node[K, V], len(n.children))
	for i, child := range n.children {
		nds[i] = child
	}
	return nds
}

func (n node[K, V]) GetChild(key K) Node[K, V] {
	if n.IsLeaf() {
		return nil
	}
	i, ok := n.entryPosition(key)
	if ok {
		return nil
	}
	return n.children[i]
}

func (n *node[K, V]) Insert(key K, value V) RootNode[K, V] {
	if root := n.insert(&entry[K, V]{
		key:   key,
		value: value,
	}); root != nil {
		return root
	}
	return nil
}

func (n *node[K, V]) Remove(key K) (RootNode[K, V], bool) {
	nn, ok := n.remove(key)
	if !ok {
		return nil, false
	}
	if nn == nil || n.IsLeaf() {
		return nil, true
	}
	if len(n.children) != 1 {
		log.Panicf("unexpected child count of %d on new root after remove of key %v", len(n.children), key)
	}
	return n.children[0], true
}

func (n *node[K, V]) insert(e *entry[K, V]) *node[K, V] {
	i, exists := n.entryPosition(e.key)
	if exists {
		// already exists, update value
		n.entries[i].value = e.value
		return nil
	}
	var child *node[K, V]
	if !n.IsLeaf() {
		// pass entry into child. If child splits, insert the resulting entry/child into this node
		splitChild := n.children[i].insert(e)
		if splitChild == nil {
			// inserted without splitting up into this node.
			return nil
		}
		// child has split, insert the resulting 'splitChild' into this node
		e = splitChild.entries[0]
		// note: splitChild.children[1] already exists as a child of this node at n.children[i], so discared in return.
		child = splitChild.children[0]
		i, _ = n.entryPosition(e.key)
	}
	n.entries = utils.InsertAtIndex(e, n.entries, i)
	if child != nil {
		n.children = utils.InsertAtIndex(child, n.children, i)
	}

	// ensure this nodes entries do not exceed limits
	if len(n.entries) < TreeDegree {
		return nil
	}

	// this node needs to split
	e, nc := n.splitNode()
	// return new parent node with single entry (e) and 2 children, new child (nc) and this node (n)
	return &node[K, V]{
		entries:  []*entry[K, V]{e},
		children: []*node[K, V]{nc, n},
	}
}

// remove attempts to remove the entry with the given key from this node or one of its child nodes.
// If the key is found and removed, it returns true, otherwise false.
// If the resulting removal results in the node having too few entries, the node returns itself, "for correction".
// If the removal can be accommodated and the node still retains the minimum entries, it returns nil.
func (n *node[K, V]) remove(key K) (*node[K, V], bool) {
	i, exists := n.entryPosition(key)
	if n.IsLeaf() {
		if !exists {
			// not found
			return nil, false
		}
		n.entries = utils.RemoveAtIndex(n.entries, i)
		// if still has enough left, return nil and true
		if len(n.entries) >= minimumKeyCount {
			return nil, true
		}
		// not enough left, return self for 'correction'
		return n, true
	}
	// Parent Node
	if !exists {
		// key not in this node, recurse down to child
		return n.removeFromChild(i, key)
	}

	// this parent node contains the key to delete.
	// merge entries children and try to get a new key.
	if n.mergeChildIntoPreviousSibling(i, false) {
		// resulting merge retained both children.  and replaced entry with new key
		return nil, true
	}
	// no key available after merge, leaving this node invalid.
	// return self for correction
	return n, true
}

// removeFromChild performs recursive call to remove, on the child at the given index.
// It checks the resulting child after remove to ensure the node is still valid.
// When not valid, it performs a merge of its own children to correct it.
func (n *node[K, V]) removeFromChild(childIndex int, key K) (*node[K, V], bool) {
	child := n.children[childIndex]
	nn, ok := child.remove(key)
	if !ok {
		return nil, false
	}
	if nn == nil || n.mergeChildIntoPreviousSibling(childIndex, true) {
		// nothing to correct from child or successfully merged children.
		return nil, true
	}
	// failed to merge children, leaving this node invalid.
	return n, true
}

func (n *node[K, V]) mergeChildIntoPreviousSibling(childIndex int, includeEntry bool) bool {
	i := childIndex
	if i == 0 {
		i++
	}
	child := n.children[i]
	sibling := n.children[i-1]
	es := sibling.entries
	if includeEntry {
		es = append(es, n.entries[i-1])
	}
	if len(child.entries) > 0 {
		es = append(es, child.entries...)
	}
	if len(child.children) > 0 {
		sibling.children = append(sibling.children, child.children...)
	}
	sibling.entries = es
	// Check if sibling can be split
	if len(sibling.entries) < minimumSplitCount {
		// not enough entries to spare for the parent key.  Return nil, making parent invalid.
		// Ensure, now redundant child is removed, as its contents are now copied into the sibling.
		n.entries = utils.RemoveAtIndex(n.entries, i-1)
		n.children = utils.RemoveAtIndex(n.children, i)
		return false
	}
	// sibling can be split across itself and child, returning (possibly) new parent key
	e, nn := sibling.splitNode()
	n.entries[childIndex-1] = e
	n.children[i] = sibling
	n.children[i-1] = nn
	return true
}

// splitNode creates a new node and moves the lower elements of this mode into the new node.
// the 'centre' or mean element is also removed from this node and returned
// any child nodes this node has are also moved.  The number of child nodes moved is equal
// to the number of entries moved, plus one.  i.e. children are split equally between the new node and this node.
func (n *node[K, V]) splitNode() (*entry[K, V], *node[K, V]) {
	eIndex := int(len(n.entries) / 2)
	e := n.entries[eIndex]
	nn := &node[K, V]{entries: n.entries[:eIndex]}
	n.entries = n.entries[eIndex+1:]
	if !n.IsLeaf() {
		nn.children = n.children[:eIndex+1]
		n.children = n.children[eIndex+1:]
	}
	return e, nn
}

// EntryPosition gets the index of the entries slice where the given key is or should appear.
// If the given key already exists, the index and true are returned
// If the given key does not exist, the index of where it would appear and false are returned.
// e.g. entries: [2,5,7]
// key:2 returns (0, true)
// key:1 returns (0, false)
// key:7 returns (2, true)
// key:6 returns (2, false)
// key:9 returns (3, false)
// key:100 returns (3, false)
func (n node[K, V]) entryPosition(key K) (int, bool) {
	for i, e := range n.entries {
		if key == e.key {
			return i, true
		}
		if key < e.key {
			return i, false
		}
	}
	return len(n.entries), false
}

func NewRootNode[K cmp.Ordered, V any]() RootNode[K, V] {
	return &node[K, V]{}
}
