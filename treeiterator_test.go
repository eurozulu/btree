package btree

import (
	"testing"
)

func TestTreeIterator_Depth(t *testing.T) {
	degree := 3
	bt := createTestTree(degree, 0)
	walker := newTreeIterator(bt.rootnode)
	depth := walker.Depth()
	if depth != 0 {
		t.Errorf("Expected 0 on Depth with nil root, got %v", depth)
	}

	testCount := 1
	bt = createTestTree(degree, testCount)
	walker = newTreeIterator(bt.rootnode)
	depth = walker.Depth()
	if depth != 1 {
		t.Errorf("Expected 1 on Depth with one node, got %v", depth)
	}

	testCount = 3
	bt = createTestTree(degree, testCount)
	walker = newTreeIterator(bt.rootnode)
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
	degree := 3
	testCount := 0
	bt := createTestTree(degree, testCount)
	walker := newTreeIterator(bt.rootnode)
	// check nil and empty tree
	next := walker.Next()
	if next != nil {
		t.Errorf("Expected nil on Next with nil root, got %v", next)
	}

	testCount = 1
	bt = createTestTree(degree, testCount)
	walker = newTreeIterator(bt.rootnode)
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

	testCount = 15
	bt = createTestTree(degree, testCount)
	walker = newTreeIterator(bt.rootnode)
	var total []nodeEntry[int, string]
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

func createTestTree(degree, count int) *bTree[int, string] {
	bt := NewBTree[int, string](degree)
	if err := fillTree(bt, count); err != nil {
		panic(err)
	}
	return bt.(*bTree[int, string])
}
