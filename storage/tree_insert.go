package storage

import "bytes"

func treeInsert(tree *BTree, node BNode, key []byte, val []byte) BNode {
	new := BNode{data: make([]byte, 2*BTREE_PAGE_SIZE)}

	idx := nodeLookupLE(node, key)

	switch node.btype() {
	case BNODE_LEAF:
		if bytes.Equal(key, node.getKey(idx)) {
			leafUpdate(new, node, idx, key, val)
		} else {
			leafInsert(new, node, idx+1, key, val)
		}

	case BNODE_NODE:
		nodeInsert(tree, new, node, idx, key, val)

	default:
		panic("bad node!")
	}
	return new
}

func (tree *BTree) Insert(key []byte, val []byte) {
	assert(len(key) != 0)
	assert(len(key) <= BTREE_MAX_KEY_SIZE)
	assert(len(val) <= BTREE_MAX_VAL_SIZE)

	if tree.root == uint64(0) {
		//create a root node.
		root := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		root.setHeader(BNODE_LEAF, uint16(2))
		nodeAppendKV(root, 0, 0, nil, nil) //only first left node has 1 extra nil, nil KV pair
		nodeAppendKV(root, 1, 0, key, val)
		tree.root = tree.new(root)
		return
	}

	node := tree.get(tree.root)
	tree.del(tree.root)

	node = treeInsert(tree, node, key, val)
	nsplit, splitted := nodeSplit3(node)

	if nsplit > 1 {
		root := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		root.setHeader(BNODE_NODE, nsplit)

		for i, knode := range splitted[:nsplit] {
			ptr, key := tree.new(knode), node.getKey(0)
			nodeAppendKV(root, uint16(i), ptr, key, nil)
		}
		tree.root = tree.new(root)
	} else {
		tree.root = tree.new(splitted[0])
	}

}
