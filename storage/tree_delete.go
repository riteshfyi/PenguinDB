package storage

import "bytes"

func treeDelete(tree *BTree, node BNode, key []byte) BNode {
	idx := nodeLookupLE(node, key)

	switch node.btype() {
	case BNODE_LEAF:

		if !bytes.Equal(node.getKey(idx), key) {
			return BNode{}
		}
		new := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		leafDelete(new, node, idx)
		return new
	case BNODE_NODE:
		return nodeDelete(tree, node, idx, key)
	default:
		panic("INVALID NODE TYPE")
	}
}

func (tree *BTree) Delete(key []byte) bool {
	assert(len(key) != 0)
	assert(len(key) <= BTREE_MAX_KEY_SIZE)

	if tree.root == uint64(0) { //no root
		return false
	}

	updated := treeDelete(tree, tree.get(tree.root), key)

	if len(updated.data) == 0 {
		return false //empty
	}
	tree.del(tree.root)

	if updated.btype() == BNODE_NODE && updated.nkeys() == 1 {
		//remove a level
		tree.root = updated.getPtr(0)
	} else {
		tree.root = tree.new(updated)
	}
	return true
}
