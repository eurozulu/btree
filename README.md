# btree
## An in memory btree implementation

## Usage
The tree stores Key/Value pairs in the order of the keys.  
Keys must support the `cmp.Ordered` type.

### New Tree:  
`NewTree[cmp.Ordered, any]()`  
Generic types indicate the Key type and value type respectively.  
The key must be an `cmp.Ordered` type  
The value may be any type  
e.g.:  
`mytree := NewTree[int, string]()`  
`mytree := NewTree[string, string]`  
`mytree := NewTree[string, *mystruct]`  
  

### Add to Tree:
`myTree.Put(123, "hello world")`  
Stores the `"hello world"` string under the `123` key  

### Retrieve from Tree:
`s, ok := myTree.Get(123)`  
Retrieves the `"hello world"` string with the `123` key.  
If key exists, `ok` is true. If key not found, `ok` is false.  
`s, ok := myTree.Get(456)`  
Will return `"", false"` as `456` doesn't exist,  

### Remove from Tree:
`ok := myTree.Remove(123)`  
Removes the entry with the `123` key if found.
If found and removed, returns true.  If not found, returns false.  

### Depth of Tree:
`depth := myTree.Depth()`  
Returns the depth of the trees leaf nodes.  
A Root leaf node has a depth of zero.
A child of root has a depth of one, and so on.

### Iterate the Tree:
`walker := myTree.Iterate()`  
Returns a tree iterator to iterate over the tree keys, in order.    
Tree Iterator first returns the very first leaf node and all its entries.  
If the tree has a depth > 0, the next iteration returns a single entry node
containing the entry from the parent node, "linking" the last node to its next peer.  
Following the single parent node, the next leaf node is returned.  
This continues until all the parents entries have been returned (as signle entry nodes).
The iterator then returns a single entry node of the parents parent node,
linking the next parent peer and follows this by the first leaf of that parent.






