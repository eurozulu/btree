package btree

import (
	"fmt"
	"testing"
)

func TestNode_IsLeaf(t *testing.T) {
	tn := &node[int, string]{}
	if !tn.IsLeaf() {
		t.Errorf("IsLeaf expected true but got false")
	}
	tn = testNodeWithChildren()
	if tn.IsLeaf() {
		t.Errorf("IsLeaf expected false but got true")
	}
}

func TestNode_Remove(t *testing.T) {
	tn := &node[int, string]{}
	root, ok := tn.Remove(0)
	if ok {
		t.Errorf("Remove on empty node expected false but got true")
	}
	if root != nil {
		t.Errorf("Remove on empty node expected nil root but got %v", root)
	}

	tn.Insert(0, "zero")
	tn.Insert(1, "one")
	tn.Insert(2, "two")
	if len(tn.entries) != 3 {
		t.Errorf("unexpected entries after insert")
	}
	root, ok = tn.Remove(1)
	if !ok {
		t.Errorf("Remove with key '1' expected true but got false")
	}
	if root != nil {
		t.Errorf("Remove with key '1'expected nil root but got %v", root)
	}
	if len(tn.entries) != 2 {
		t.Errorf("unexpected entries after insert")
	}
	if tn.entries[1].key != 2 {
		t.Errorf("After remove with key '1' expected key 2 at position 2, found %v", tn.entries[1].key)
	}
	root, ok = tn.Remove(1)
	if ok {
		t.Errorf("after Remove with key '1', duplicate remove expected false but got true")
	}
	if _, ok := tn.Value(1); ok {
		t.Errorf("expected false on Value of removed key '1' but got true")
	}
	root, ok = tn.Remove(2)
	if !ok {
		t.Errorf("Remove with key '2' expected true but got false")
	}
	if root != nil {
		t.Errorf("Remove with key '2' expected nil root but got %v", root)
	}
	if len(tn.entries) != 1 {
		t.Errorf("unexpected entries after insert")
	}

	root, ok = tn.Remove(0)
	if !ok {
		t.Errorf("Remove with key '0' expected true but got false")
	}
	if root != nil {
		t.Errorf("Remove with key '0' expected nil root but got %v", root)
	}
	if len(tn.entries) != 0 {
		t.Errorf("unexpected entries after insert")
	}

	tn = testThreeNode()
	root, ok = tn.Remove(30)
	if !ok {
		t.Errorf("Remove with key '30' expected true but got false")
	}
	if root == nil {
		t.Errorf("Remove with key '30' expected non-nil root but got nil")
	}
	if !root.IsLeaf() {
		t.Errorf("Remove with key '30' expected root to be leaf")
	}
	keys := root.Keys()
	if len(keys) != 2 {
		t.Errorf("unexpected 2 keys after delete, found %v", len(keys))
	}
	if keys[0] != 10 {
		t.Errorf("keys[0] expected 10 but got %v", keys[0])
	}
	if keys[1] != 20 {
		t.Errorf("keys[1] expected 20 but got %v", keys[0])
	}

	tn = testThreeNode()
	root, ok = tn.Remove(10)
	if !ok {
		t.Errorf("Remove with key '10' expected true but got false")
	}
	if root == nil {
		t.Errorf("Remove with key '10' expected non-nil root but got nil")
	}
	if !root.IsLeaf() {
		t.Errorf("Remove with key '10' expected root to be leaf")
	}
	keys = root.Keys()
	if len(keys) != 2 {
		t.Errorf("unexpected 2 keys after delete, found %v", len(keys))
	}
	if keys[0] != 20 {
		t.Errorf("keys[0] expected 20 but got %v", keys[0])
	}
	if keys[1] != 30 {
		t.Errorf("keys[1] expected 30 but got %v", keys[1])
	}

	tn = testThreeNode()
	root, ok = tn.Remove(20)
	if !ok {
		t.Errorf("Remove with key '20' expected true but got false")
	}
	if root == nil {
		t.Errorf("Remove with key '20' expected non-nil root but got nil")
	}
	if !root.IsLeaf() {
		t.Errorf("Remove with key '20' expected root to be leaf")
	}
	keys = root.Keys()
	if len(keys) != 2 {
		t.Errorf("unexpected 2 keys after delete, found %v", len(keys))
	}
	if keys[0] != 10 {
		t.Errorf("keys[0] expected 10 but got %v", keys[0])
	}
	if keys[1] != 30 {
		t.Errorf("keys[1] expected 30 but got %v", keys[1])
	}

}

func TestNode_Insert(t *testing.T) {
	tn := &node[int, string]{}
	tn.Insert(0, "zero")
	if len(tn.entries) != 1 {
		t.Errorf("Insert() expected 1 element but got %d", len(tn.entries))
	}
	tn.Insert(0, "zeroooo")
	if len(tn.entries) != 1 {
		t.Errorf("Insert() expected 1 element but got %d", len(tn.entries))
	}
	tn.Insert(10, "ten")
	if len(tn.entries) != 2 {
		t.Errorf("Insert() expected 2 element but got %d", len(tn.entries))
	}
	if tn.entries[0].key != 0 {
		t.Errorf("Insert() expected first entry with key '0' but got %d", tn.entries[0].key)
	}
	if tn.entries[1].key != 10 {
		t.Errorf("Insert() expected second entry with key '10' but got %d", tn.entries[1].key)
	}
	tn.Insert(20, "twenty")
	if len(tn.entries) != 3 {
		t.Errorf("Insert() expected 3 element but got %d", len(tn.entries))
	}
	if tn.entries[2].key != 20 {
		t.Errorf("Insert() expected third entry with key '20' but got %d", tn.entries[2].key)
	}
	tn.Insert(20, "two tens")
	if len(tn.entries) != 3 {
		t.Errorf("Insert() expected 3 element but got %d", len(tn.entries))
	}
}

func TestNode_Value(t *testing.T) {
	tn := &node[int, string]{}
	if err := testNodeHasValue(tn, -1, false, ""); err != nil {
		t.Error(err)
	}

	tn.Insert(0, "zero")
	if err := testNodeHasValue(tn, -1, false, ""); err != nil {
		t.Error(err)
	}
	if err := testNodeHasValue(tn, 0, true, "zero"); err != nil {
		t.Error(err)
	}
	tn.Insert(0, "zerooo")
	if err := testNodeHasValue(tn, 0, true, "zerooo"); err != nil {
		t.Error(err)
	}

	tn = BuildTestTreeRoot().(*node[int, string])
	if err := testNodeHasValue(tn, 0, false, ""); err != nil {
		t.Error(err)
	}
	if err := testNodeHasValue(tn, 1, true, "one"); err != nil {
		t.Error(err)
	}
	if err := testNodeHasValue(tn, 90, true, "ninty"); err != nil {
		t.Error(err)
	}
}

func TestNode_GetChild(t *testing.T) {
	TreeDegree = 3
	tn := &node[int, string]{}
	rn := tn.GetChild(10)
	if rn != nil {
		t.Errorf("GetChild(0) expected nil on empty node but got %v", rn)
	}
	tn.Insert(10, "ten")
	rn = tn.GetChild(10) // should be nil as this node contains that key
	if rn != nil {
		t.Errorf("GetChild(0) expected nil on empty node but got %v", rn)
	}
	tn.Insert(20, "twenty")
	tn = tn.Insert(30, "thirty").(*node[int, string])
	if tn == nil {
		t.Errorf("new root expected non-nil new root node but got nil")
	}
	if tn.IsLeaf() {
		t.Errorf("new root expected non leaf node but found leaf")
	}
	rn = tn.GetChild(10)
	if rn == nil {
		t.Errorf("expected non-nil for key 10 but got nil")
	}
	keys := rn.Keys()
	if len(keys) != 1 {
		t.Errorf("expected 1 key but got %v", len(keys))
	}
	if keys[0] != 10 {
		t.Errorf("keys[0] expected 10 but got %v", keys[0])
	}

}

func TestNode_Keys(t *testing.T) {
	tn := &node[int, string]{}
	if tn.Keys() != nil {
		t.Errorf("Keys() expected nil on empty node but got %v", tn.Keys())
	}
	tn.Insert(0, "zero")
	keys := tn.Keys()
	if len(keys) != 1 {
		t.Errorf("Keys() expected 1 element but got %d", len(keys))
	}
	if keys[0] != 0 {
		t.Errorf("Keys() expected element '0' but got %d", keys[0])
	}
	tn.Insert(60, "ten")
	tn.Insert(40, "ten")
	tn.Insert(10, "ten")
	keys = tn.Keys()
	if len(keys) != 4 {
		t.Errorf("Keys() expected 4 element but got %d", len(keys))
	}
	if keys[0] != 0 {
		t.Errorf("Keys() expected element '0' but got %d", keys[0])
	}
	if keys[1] != 10 {
		t.Errorf("Keys() expected element '10' but got %d", keys[0])
	}
	if keys[2] != 40 {
		t.Errorf("Keys() expected element '40' but got %d", keys[0])
	}
	if keys[3] != 60 {
		t.Errorf("Keys() expected element '40' but got %d", keys[0])
	}
}

func testNodeHasValue(tn *node[int, string], key int, expectOk bool, expectvalue string) error {
	v, ok := tn.Value(key)
	if ok != expectOk {
		return fmt.Errorf("Value() expected %v with key %d, but got %v", expectOk, key, ok)
	}
	if v != expectvalue {
		return fmt.Errorf("Value() expected '%s' with key %d, but got '%s'", expectvalue, key, v)
	}
	return nil
}
