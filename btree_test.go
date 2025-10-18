package btree

import (
	"cmp"
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"testing"
	"time"
)

const testKey0 = "zero"

func TestBTree_Add_NoSplit(t *testing.T) {
	bt := NewBTree[int, string](3)
	count := 2
	if err := fillTree(bt, count); err != nil {
		t.Error(err)
	}
	if bt.Depth() != 0 {
		t.Errorf("Unexpected depth after adding %d items to tree of degree 3.  Expected %d depth, found %d", count, 0, bt.Depth())
	}
	if err := checkContains(bt, count); err != nil {
		t.Error(err)
	}

}

func TestBTree_Add_withSplit(t *testing.T) {
	bt := NewBTree[int, string](3)
	count := 3
	if err := fillTree(bt, count); err != nil {
		t.Error(err)
	}
	if bt.Depth() != 1 {
		t.Errorf("Unexpected depth after adding %d items to tree of degree 3.  Expected %d depth, found %d", count, 1, bt.Depth())
	}
	if err := checkContains(bt, count); err != nil {
		t.Error(err)
	}
}

func TestBTree_Add_threeLayers(t *testing.T) {
	degree := 3
	bt := NewBTree[int, string](degree)
	count := 13
	if err := fillTree(bt, count); err != nil {
		t.Error(err)
	}
	if bt.Depth() != 2 {
		t.Errorf("Unexpected depth after adding %d items to tree of degree %d.  Expected %d depth, found %d", count, degree, 2, bt.Depth())
	}
	if err := checkContains(bt, count); err != nil {
		t.Error(err)
	}
}

func TestBTree_Remove_From_Leaf_Root(t *testing.T) {
	bt := NewBTree[int, string](3)
	if err := bt.Remove(123); err == nil {
		t.Error("Expected error removing unknown key from empty tree, got nil")
	}

	count := 2
	if err := fillTree(bt, count); err != nil {
		t.Error(err)
	}
	if err := bt.Remove(123); err == nil {
		t.Error("Expected error removing unknown key from populated tree, got nil")
	}

	root := bt.(*bTree[int, string]).rootnode
	if len(root.Entries) != 2 {
		t.Errorf("expected %d entry count, found %d", 2, len(root.Entries))
	}
	if err := bt.Remove(1); err != nil {
		t.Errorf("unexpected error removing item %d from tree.  %v", 1, err)
	}
	if len(root.Entries) != 1 {
		t.Errorf("expected %d entry count, found %d", 1, len(root.Entries))
	}
	if err := bt.Remove(0); err != nil {
		t.Errorf("unexpected error removing item %d from tree.  %v", 0, err)
	}
	if len(root.Entries) != 0 {
		t.Errorf("expected %d entry count, found %d", 0, len(root.Entries))
	}
}

func TestBTree_Remove_From_Peer(t *testing.T) {
	bt := NewBTree[int, string](3)
	count := 3
	if err := fillTree(bt, count); err != nil {
		t.Error(err)
	}

	// Test Back fill
	if err := bt.Remove(2); err != nil {
		t.Errorf("unexpected error removing item %d from tree.  %v", 2, err)
	}
	root := bt.(*bTree[int, string]).rootnode
	if len(root.Entries) != 2 {
		t.Errorf("expected entry count in root after delete, expected %d found %d", 2, len(root.Entries))
	}
	if !root.IsLeaf() {
		t.Errorf("unexpected %d children found on root, expected to be a lead", len(root.Children))
	}
	if root.Entries[0].Key != 0 {
		t.Errorf("unexpected item in root index 0, expected %d at entry zero, got %d", 0, root.Entries[0].Key)
	}

	if root.Entries[1].Key != 1 {
		t.Errorf("unexpected item in root index 1, expected %d at entry zero, got %d", 1, root.Entries[1].Key)
	}

	// Test forward fill
	bt = NewBTree[int, string](3)
	if err := fillTree(bt, count); err != nil {
		t.Error(err)
	}
	if err := bt.Remove(1); err != nil {
		t.Errorf("unexpected error removing item %d from tree.  %v", 1, err)
	}
	root = bt.(*bTree[int, string]).rootnode
	if len(root.Entries) != 2 {
		t.Errorf("expected entry count in root after delete, expected %d found %d", 2, len(root.Entries))
	}
	if !root.IsLeaf() {
		t.Errorf("unexpected %d children found on root, expected to be a lead", len(root.Children))
	}
	if root.Entries[0].Key != 0 {
		t.Errorf("unexpected key in root index 0, expected %d at entry zero, got %d", 0, root.Entries[0].Key)
	}
	if root.Entries[1].Key != 2 {
		t.Errorf("unexpected item in root index 1, expected %d at entry zero, got %d", 2, root.Entries[0].Key)
	}
}

func TestBTree_Remove_From_Parent(t *testing.T) {
	bt := NewBTree[int, string](3)
	count := 3
	if err := fillTree(bt, count); err != nil {
		t.Error(err)
	}

	if err := bt.Remove(1); err != nil {
		t.Errorf("unexpected error removing key %d from tree.  %v", 1, err)
	}
	root := bt.(*bTree[int, string]).rootnode
	if len(root.Entries) != 2 {
		t.Errorf("expected entry count in root after delete, expected %d found %d", 2, len(root.Entries))
	}
	if !root.IsLeaf() {
		t.Errorf("unexpected %d children found on root, expected to be a lead", len(root.Children))
	}
	if root.Entries[0].Key != 0 {
		t.Errorf("unexpected key in root index 0, expected %d at entry zero, got %d", 0, root.Entries[0].Key)
	}
	if root.Entries[1].Key != 2 {
		t.Errorf("unexpected item in root index 1, expected %d at entry zero, got %d", 2, root.Entries[0].Key)
	}
}

func TestBTree_Remove_From_Root_3layer(t *testing.T) {
	bt := NewBTree[int, string](3)
	count := 7
	if err := fillTree(bt, count); err != nil {
		t.Error(err)
	}
	if err := validateTree(bt); err != nil {
		t.Error(err)
	}
	if bt.Depth() != 2 {
		t.Errorf("unexpected depth, expected %d, found %d", 2, bt.Depth())
	}

	if err := bt.Remove(3); err != nil {
		t.Errorf("unexpected error removing item %d from tree.  %v", 4, err)
	}
	if err := validateTree(bt); err != nil {
		t.Error(err)
	}
	root := bt.(*bTree[int, string]).rootnode
	if err := checkContainsEntries(root, 2, 5); err != nil {
		t.Error(err)
	}
	if err := checkContainsEntries(&root.Children[0], 0, 1); err != nil {
		t.Error(err)
	}
	if err := checkContainsEntries(&root.Children[1], 4); err != nil {
		t.Error(err)
	}
	if err := checkContainsEntries(&root.Children[2], 6); err != nil {
		t.Error(err)
	}

	// root now has 2 entries/3 children (2,5)
	if err := bt.Remove(5); err != nil {
		t.Errorf("unexpected error removing item %d from tree.  %v", 6, err)
	}
	if err := validateTree(bt); err != nil {
		t.Error(err)
	}
	root = bt.(*bTree[int, string]).rootnode
	if err := checkContainsEntries(root, 1, 4); err != nil {
		t.Error(err)
	}
	if err := checkContainsEntries(&root.Children[0], 0); err != nil {
		t.Error(err)
	}
	if err := checkContainsEntries(&root.Children[1], 2); err != nil {
		t.Error(err)
	}
	if err := checkContainsEntries(&root.Children[2], 6); err != nil {
		t.Error(err)
	}

}

func TestBTree_Remove_From_Root_4layer(t *testing.T) {
	bt := NewBTree[int, string](3)
	count := 15
	if err := fillTree(bt, count); err != nil {
		t.Error(err)
	}
	if err := validateTree(bt); err != nil {
		t.Error(err)
	}
	if bt.Depth() != 3 {
		t.Errorf("unexpected depth, expected %d, found %d", 2, bt.Depth())
	}
	root := bt.(*bTree[int, string]).rootnode
	if err := checkContainsEntries(root, 7); err != nil {
		t.Error(err)
	}

	// 7 is only root entry
	if err := bt.Remove(7); err != nil {
		t.Errorf("unexpected error removing item %d from tree.  %v", 10, err)
	}
	if err := validateTree(bt); err != nil {
		t.Error(err)
	}
	if bt.Depth() != 2 {
		t.Errorf("unexpected depth, expected %d, found %d", 2, bt.Depth())
	}
	root = bt.(*bTree[int, string]).rootnode
	if err := checkContainsEntries(root, 6, 11); err != nil {
		t.Error(err)
	}
	if err := checkContainsEntries(&root.Children[0], 1, 3); err != nil {
		t.Error(err)
	}
	if err := checkContainsEntries(&root.Children[1], 9); err != nil {
		t.Error(err)
	}
	if err := checkContainsEntries(&root.Children[2], 13); err != nil {
		t.Error(err)
	}

	if err := checkContainsEntries(&root.Children[0].Children[0], 0); err != nil {
		t.Error(err)
	}
	if err := checkContainsEntries(&root.Children[0].Children[1], 2); err != nil {
		t.Error(err)
	}
	if err := checkContainsEntries(&root.Children[0].Children[2], 4, 5); err != nil {
		t.Error(err)
	}
	if err := checkContainsEntries(&root.Children[2].Children[0], 12); err != nil {
		t.Error(err)
	}
	if err := checkContainsEntries(&root.Children[2].Children[1], 14); err != nil {
		t.Error(err)
	}

}

func checkContains(bt BTree[int, string], count int) error {
	for i := 1; i < count; i++ {
		v := bt.Get(i)
		if v == nil {
			return fmt.Errorf("Expected value for key '%v' but got nil", i)
		}
		expect := "-" + strconv.Itoa(i) + "-"
		if (*v) != expect {
			return fmt.Errorf("Expected value for key '%d' of %s but got %s", i, expect, *v)
		}
	}
	return nil
}

func checkContainsEntries(n *node[int, string], entries ...int) error {
	if len(n.Entries) < len(entries) {
		return fmt.Errorf("unexpected entry count, expected %d entries, found %d", len(entries), len(n.Entries))
	}
	sort.Ints(entries)
	for i, entry := range entries {
		if n.Entries[i].Key != entry {
			return fmt.Errorf("Expected entry key '%v', found '%v'", entry, n.Entries[i].Key)
		}
	}
	return nil
}

func fillTree(bt BTree[int, string], count int) error {
	for i := 0; i < count; i++ {
		v := "-" + strconv.Itoa(i) + "-"
		if err := bt.Add(i, &v); err != nil {
			return err
		}
	}
	return nil
}

func validateTree(bt BTree[int, string]) error {
	btp := bt.(*bTree[int, string])
	return validateNode(btp.rootnode, btp.degree, 0)
}
func validateNode(n *node[int, string], degree, depth int) error {
	el := len(n.Entries)

	if el >= degree {
		return fmt.Errorf("Invalid node at depth %d has %d entries when degree is %d  %v", depth, el, degree, n)
	}
	if el < degree/2 {
		return fmt.Errorf("Invalid node at depth %d only has %d entried when minimum is %d  %v", depth, el, degree/2, n)
	}
	if len(n.Entries) > 1 {
		for i := 1; i < len(n.Entries); i++ {
			if cmp.Compare(n.Entries[i-1].Key, n.Entries[i].Key) > 0 {
				return fmt.Errorf("Invalid node at depth %d entries out of order %v", depth, n)
			}
		}
	}
	if n.IsLeaf() {
		return nil
	}
	cl := len(n.Children)
	if el+1 != cl {
		return fmt.Errorf("Invalid node at depth %d has %d children, expected %d  %v", depth, cl, el+1, n)
	}

	for i, child := range n.Children {
		if err := validateNode(&child, degree, depth+1); err != nil {
			return err
		}
		if i < len(n.Entries) {
			lastChildKey := child.Entries[len(child.Entries)-1].Key
			if cmp.Compare(lastChildKey, n.Entries[i].Key) != -1 {
				return fmt.Errorf("Invalid node at depth %d entries child last key %d is not smaller than parent entry %d", depth, lastChildKey, n.Entries[i].Key)
			}
			continue
		}
		firstChildKey := child.Entries[0].Key
		if cmp.Compare(firstChildKey, n.Entries[i-1].Key) != 1 {
			return fmt.Errorf("Invalid node at depth %d entries child first key %d is not larger than parent entry %d", depth, firstChildKey, n.Entries[i].Key)
		}
	}
	return nil
}
func keyForInt(i int) string {
	return fmt.Sprintf("%x", i)
}

func TestCompareToHashmap(t *testing.T) {
	count := 10000000
	testChecks := 100000
	var memRef runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memRef)

	bt := NewBTree[int, string](100)
	tm := time.Now()
	fillTree(bt, count)
	t.Logf("took %v to fill tree with %d elements\n", time.Since(tm), count)
	tm = time.Now()
	for i := 0; i < testChecks; i++ {
		key := rand.Intn(count)
		v := bt.Get(key)
		if v == nil {
			t.Errorf("Failed to find key %v in tree", v)
			continue
		}
		expect := "-" + strconv.Itoa(key) + "-"
		if *v != expect {
			t.Errorf("Invalid value for key %d in tree, expected %s, got %s", v, expect, *v)
		}
	}
	t.Logf("took %v to read tree %d times\n", time.Since(tm), testChecks)
	showMemoryStats(&memRef)
	runtime.GC()
	runtime.ReadMemStats(&memRef)

	m := make(map[string]interface{}, count)
	tm = time.Now()
	fillMap(m, count)
	t.Logf("took %v to fill map with %d elements\n", time.Since(tm), count)
	tm = time.Now()
	for i := 0; i < testChecks; i++ {
		v := rand.Intn(count)
		key := keyForInt(v)
		nv := m[key]
		if nv == nil {
			t.Errorf("Failed to find key %s in map", key)
			continue
		}
		if nv.(int) != v {
			t.Errorf("Invalid value for key %s in map, expected %d, got %d", key, v, nv.(int))
		}
	}

	t.Logf("took %v to read map %d times\n", time.Since(tm), testChecks)
	showMemoryStats(&memRef)
}

func showMemoryStats(ref *runtime.MemStats) {
	var memNow runtime.MemStats
	runtime.ReadMemStats(&memNow)
	totalMem := memNow.TotalAlloc - ref.TotalAlloc
	totalAllocs := memNow.Mallocs - ref.Mallocs
	fmt.Printf("total memory used: %s  (%d)\n", byteString(totalMem), totalMem)
	fmt.Printf("mallocs: %d\n", totalAllocs)
}

var byteNames = []string{
	"bytes", "kilobytes", "megabytes", "gigabytes",
}

func byteString(size uint64) string {
	i := 0
	s := float64(size)
	for ; s > 1024; s /= 1024 {
		if i+1 >= len(byteNames) {
			break
		}
		i++
	}
	return fmt.Sprintf("%f %s", s, byteNames[i])
}

func fillMap(m map[string]interface{}, count int) error {
	for i := 0; i < count; i++ {
		key := keyForInt(i + 1)
		m[key] = i + 1
	}
	return nil
}
