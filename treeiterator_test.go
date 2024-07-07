package btree

import (
	"testing"
)

func TestTreeIterator_HasNext(t *testing.T) {
	degree := 3
	bt := NewBTree[int, string](degree)

	// empty node == zero depth
	walker := newTreeIterator[int, string](bt)
	if walker.HasNext() {
		t.Error("Expected false HasNext on empty node")
	}
	if walker.Next() != nil {
		t.Error("Expected Next as nil on empty node")
	}

	fillTree(bt, 1)
	walker = newTreeIterator[int, string](bt)
	if !walker.HasNext() {
		t.Error("Expected true HasNext on single entry node before next")
	}
	_ = walker.Next()
	if walker.HasNext() {
		t.Error("Expected false HasNext on single entry node after next")
	}

	// test tree node tree
	bt = NewBTree[int, string](degree)
	testCount := 3
	fillTree(bt, testCount)
	walker = newTreeIterator[int, string](bt)
	for i := 0; i < testCount; i++ {
		if !walker.HasNext() {
			t.Errorf("Expected HasNext true with test tree before %d iteration", i+1)
		}
		_ = walker.Next()
	}
	_ = walker.Next()
	if walker.HasNext() {
		t.Error("Expected HasNext true after first next")
	}

}

func TestTreeIterator_Depth(t *testing.T) {
	walker := newTreeIterator[int, string](nil)
	depth := walker.Depth()
	if depth != 0 {
		t.Errorf("Expected 0 on Depth with nil root, got %v", depth)
	}

	degree := 3
	testCount := 1
	bt := NewBTree[int, string](degree)
	fillTree(bt, testCount)
	walker = newTreeIterator[int, string](bt)
	depth = walker.Depth()
	if depth != 1 {
		t.Errorf("Expected 1 on Depth with one node, got %v", depth)
	}

	fillTree(bt, 3)
	walker = newTreeIterator[int, string](bt)
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

	degree := 3
	bt := NewBTree[int, string](degree)
	walker = newTreeIterator[int, string](bt)
	next = walker.Next()
	if next != nil {
		t.Errorf("Expected nil on Next with empty node, got %v", next)
	}

	fillTree(bt, 1)
	walker = newTreeIterator[int, string](bt)
	next = walker.Next()
	if next == nil {
		t.Errorf("Expected non nil result on Next on single entry node")
	}
	if len(next) != 1 {
		t.Errorf("Expected 1 element returned, got %d", len(next))
	}
	if next[0].Key != 0 {
		t.Errorf("Expected returned key 0, got %d", next[0].Key)
	}
	if next[0].Value == nil {
		t.Errorf("Expected non nil value for key %v", next[0].Key)
	}
	expect := "-0-"
	if *next[0].Value != expect {
		t.Errorf("Expected returned value '%s', got '%s'", expect, *next[0].Value)
	}
	next = walker.Next()
	if next != nil {
		t.Errorf("Expected nil on Next with empty node, got %v", next)
	}
	if walker.HasNext() {
		t.Error("Expected false HasNext on single entry node after next")
	}

	bt = NewBTree[int, string](degree)
	testCount := 15
	fillTree(bt, testCount)
	walker = newTreeIterator[int, string](bt)
	var total []NodeEntry[int, string]
	for walker.HasNext() {
		next = walker.Next()
		if len(next) < 1 {
			t.Error("Expected non empty entires for next iteration")
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
	if len(total) != testCount {
		t.Errorf("Expected %d elements returned from test tree, got %d", testCount, len(total))
	}
	lastKey := -1
	for _, entry := range total {
		if entry.Key <= lastKey {
			t.Errorf("unexpected entry out of order. previous key was %d and found key %d following", lastKey, entry.Key)
		}
		lastKey = entry.Key
	}
}
