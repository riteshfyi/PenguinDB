package storage

func leafDelete(new BNode, old BNode, idx uint16) {
	new.setHeader(BNODE_LEAF, old.nkeys()-1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendRange(new, old, idx, idx+1, old.nkeys()-(idx+1))
}

func nodeDelete(tree *BTree, node BNode, idx uint16, key []byte) BNode {
	kptr := node.getPtr(idx)
	knode := tree.get(kptr)
	updated := treeDelete(tree, knode, key)

	if len(updated.data) == 0 {
		return BNode{}
	}
	tree.del(kptr)

	new := BNode{data: make([]byte, BTREE_PAGE_SIZE)}

	mergeDir, sibling := shouldMerge(tree, node, idx, updated)

	switch {

	case mergeDir < 0: //left
		merged := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		nodeMerge(merged, sibling, updated)
		tree.del(node.getPtr(idx - 1))
		nodeReplace2Kid(new, node, idx-1, tree.new(merged), merged.getKey(0))

	case mergeDir > 0:
		merged := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		nodeMerge(merged, updated, sibling)
		tree.del(node.getPtr(idx + 1))
		nodeReplace2Kid(new, node, idx, tree.new(merged), merged.getKey(0))

	case mergeDir == 0:
		assert(updated.nkeys() > 0)
		nodeReplaceKidN(tree, new, node, idx, updated)
	}

	return new
}

func shouldMerge(tree *BTree, node BNode, idx uint16, updated BNode) (int, BNode) {
	/*
		The conditions for merging are:
		1. The node is smaller than 1/4 of a page (this is arbitrary).
		2. Has a sibling and the merged result does not exceed one page.
	*/
	if updated.nbytes() > BTREE_PAGE_SIZE/4 {
		return 0, BNode{}
	}

	if idx > 0 {
		sibling := tree.get(node.getPtr(idx - 1))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return -1, sibling
		}
	}

	if idx+1 < node.nkeys() {
		sibling := tree.get(node.getPtr(idx + 1))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return +1, sibling
		}
	}

	return 0, BNode{}
}

func nodeMerge(new BNode, left BNode, right BNode) {
	new.setHeader(left.btype(), left.nkeys()+right.nkeys())
	nodeAppendRange(new, left, 0, 0, left.nkeys())
	nodeAppendRange(new, right, left.nkeys(), 0, right.nkeys())
}

// update the links, remove the one deleted node
func nodeReplace2Kid(new BNode, old BNode, idx uint16, merged uint64, key []byte) {
	//remove idx + 1, key, val, ptr, etc. ptr already dereferenced
	//set new key as the key given , ptr as merged,
	new.setHeader(old.btype(), old.nkeys()-1)
	nodeAppendRange(new, old, idx, idx+1, old.nkeys()-idx-1)
	copy(old.data[old.kvPos(idx):], key)
	new.setPtr(idx, merged)
}
