package btree

import (
	"testing"
)

func TestTreeIterator_HasNext(t *testing.T) {
	// empty node == zero depth
	n := &node[int, string]{}
	walker := newTreeIterator[int, string](n)
	if walker.HasNext() {
		t.Error("Expected false HasNext on empty node")
	}

	n = buildTestNode(0, "zero").(*node[int, string])
	walker = newTreeIterator[int, string](n)
	if !walker.HasNext() {
		t.Error("Expected true HasNext on single entry node before next")
	}
	_ = walker.Next()
	if walker.HasNext() {
		t.Error("Expected false HasNext on single entry node after next")
	}

	n = BuildTestTreeRoot().(*node[int, string])
	walker = newTreeIterator[int, string](n)
	if !walker.HasNext() {
		t.Error("Expected HasNext true with test tree before first next")
	}
	_ = walker.Next()
	if !walker.HasNext() {
		t.Error("Expected HasNext true after first next")
	}
}

func TestTreeIterator_Depth(t *testing.T) {
	walker := newTreeIterator[int, string](nil)
	depth := walker.Depth()
	if depth != 0 {
		t.Errorf("Expected 0 on Depth with nil root, got %v", depth)
	}

	n := buildTestNode(0, "zero").(*node[int, string])
	walker = newTreeIterator[int, string](n)
	depth = walker.Depth()
	if depth != 1 {
		t.Errorf("Expected 1 on Depth with one node, got %v", depth)
	}

	walker = newTreeIterator[int, string](testThreeNode())
	depth = walker.Depth()
	if depth != 2 {
		t.Errorf("Expected 2 on Depth with three nodes, got %v", depth)
	}
	_ = walker.Next()
	depth = walker.Depth()
	if depth != 1 {
		t.Errorf("Expected 1 on Depth with three nodes, got %v", depth)
	}
	_ = walker.Next()
	depth = walker.Depth()
	if depth != 2 {
		t.Errorf("Expected 2 on Depth with three nodes, got %v", depth)
	}
	_ = walker.Next()
	depth = walker.Depth()
	if depth != 0 {
		t.Errorf("Expected 0 on Depth with three nodes, got %v", depth)
	}

}

func TestTreeIterator_Next(t *testing.T) {
	// check nil and empty roots
	walker := newTreeIterator[int, string](nil)
	next := walker.Next()
	if next != nil {
		t.Errorf("Expected nil on Next with nil root, got %v", next)
	}

	walker = newTreeIterator[int, string](&node[int, string]{})
	next = walker.Next()
	if next != nil {
		t.Errorf("Expected nil on Next with empty node, got %v", next)
	}

	n := buildTestNode(0, "zero").(*node[int, string])
	walker = newTreeIterator[int, string](n)
	next = walker.Next()
	if next == nil {
		t.Errorf("Expected non nil result on Next on single entry node")
	}
	if len(next) != 1 {
		t.Error("Expected 1 element returned, got ", len(next))
	}
	if next[0].Key() != 0 {
		t.Error("Expected returned key 0, got ", next[0].Key())
	}
	if next[0].Value() != "zero" {
		t.Error("Expected returned value 'zero', got ", next[0].Value())
	}
	next = walker.Next()
	if next != nil {
		t.Errorf("Expected nil on Next with empty node, got %v", next)
	}
	if walker.HasNext() {
		t.Error("Expected false HasNext on single entry node after next")
	}

	n = BuildTestTreeRoot().(*node[int, string])
	walker = newTreeIterator[int, string](n)
	var total []Entry[int, string]
	for walker.HasNext() {
		next = walker.Next()
		if next == nil {
			t.Error("Expected non nil result on Next at this depth")
		}
		if len(next) < 1 {
			t.Error("Expected 4 element returned, got ", len(next))
		}
		total = append(total, next...)
		next = walker.Next()
		if next == nil {
			break
		}
		if len(next) != 1 {
			t.Error("Expected 1 element returned, got ", len(next))
		}
		total = append(total, next...)
	}
	if len(total) != 35 {
		t.Error("Expected 35 elements returned from test tree, got ", len(total))
	}
	var lastKey int
	for _, entry := range total {
		if entry.Key() <= lastKey {
			t.Errorf("unexpected entry out of order. previous key was %d and found key %d following", lastKey, entry.Key())
		}
		lastKey = entry.Key()
	}
}
