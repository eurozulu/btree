package btree

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestShowTree(t *testing.T) {
	TreeDegree = 4
	tree := NewTree[int, string]().(*btree[int, string])
	tree.rootNode = BuildTestTreeRoot()

	buf := bytes.NewBuffer(nil)
	if err := ShowTree[int, string](tree, buf); err != nil {
		t.Fatal(err)
	}
	fmt.Println(buf.String())
}

func TestBtree_Put(t *testing.T) {
	TreeDegree = 100
	tree := NewTree[int, string]().(*btree[int, string])
	for i := 0; i < 1495598; i++ {
		tree.Put(i, "-"+strconv.Itoa(i)+"-")
	}
	expectDepth := 3
	depth := tree.Depth()
	if depth != expectDepth {
		t.Errorf("Depth expected %d, got %d", expectDepth, depth)
	}
	nlr := &nodeLevelReader[int, string]{root: tree.rootNode}
	expectCount := 29325
	nodes := nlr.NodesAtDepth(expectDepth)
	if len(nodes) != expectCount {
		t.Errorf("expected %d Nodes at depth %d, got %d", expectCount, expectDepth, len(nodes))
	}
	expectDepth--
	expectCount = 575
	nodes = nlr.NodesAtDepth(expectDepth)
	if len(nodes) != expectCount {
		t.Errorf("expected %d Nodes at depth %d, got %d", expectCount, expectDepth, len(nodes))
	}
	expectDepth--
	expectCount = 11
	nodes = nlr.NodesAtDepth(expectDepth)
	if len(nodes) != expectCount {
		t.Errorf("expected %d Nodes at depth %d, got %d", expectCount, expectDepth, len(nodes))
	}

	expectDepth--
	expectCount = 1
	nodes = nlr.NodesAtDepth(expectDepth)
	if len(nodes) != expectCount {
		t.Errorf("expected %d Nodes at depth %d, got %d", expectCount, expectDepth, len(nodes))
	}

}

func TestBtree_Get(t *testing.T) {
	TreeDegree = 100
	tree := NewTree[int, string]().(*btree[int, string])
	for i := 0; i < 1495600; i++ {
		tree.Put(i, "-"+strconv.Itoa(i)+"-")
	}

	if err := testGet(tree, 1495599); err != nil {
		t.Fatal(err)
	}

	if err := testGet(tree, 0); err != nil {
		t.Fatal(err)
	}

	if err := testGet(tree, 1495600/2); err != nil {
		t.Fatal(err)
	}

}

func testGet(tree BTree[int, string], key int) error {
	tm := time.Now()
	s, ok := tree.Get(key)
	fmt.Printf("took: %s\n", time.Since(tm))
	if !ok {
		return fmt.Errorf("Get returned false on key '%v' expected to be present", key)
	}
	expect := strings.Join([]string{"-", strconv.Itoa(key), "-"}, "")
	if s != expect {
		return fmt.Errorf("Get expected %s value for key %v, got %s", expect, key, s)
	}
	return nil
}
