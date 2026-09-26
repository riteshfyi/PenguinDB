package storage

import "bytes"

func treeGet(tree *BTree, node BNode, key []byte) ([]byte, bool) {
	//get the current leaf type
	//just verify the key first okay ?
	assert(0 <= len(key) && len(key) <= BTREE_MAX_KEY_SIZE)

	idx := nodeLookupLE(node, key)
	switch node.btype() {
	case BNODE_LEAF:
		if bytes.Equal(key, node.getKey(idx)) {
			return node.getVal(idx), true
		}
		return []byte{}, false
	case BNODE_NODE:
		knode := tree.get(node.getPtr(idx))
		return treeGet(tree, knode, key)
	default:
		panic("Invalid Node Type")
	}
}

func (tree *BTree) Get(key []byte) ([]byte, bool) {
	assert(0 <= len(key) && len(key) <= BTREE_MAX_KEY_SIZE)
	node := tree.get(tree.root)
	val, found := treeGet(tree, node, key)
	return val, found
}
