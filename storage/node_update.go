package storage

import (
	"encoding/binary"
)

func leafInsert(new BNode, old BNode, idx uint16, key []byte, val []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys()+1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}

func leafUpdate(new BNode, old BNode, idx uint16, key []byte, val []byte) {
	//key already exists, so just update the val. make sure to update the offsets
	new.setHeader(BNODE_LEAF, old.nkeys())
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx+1, old.nkeys()-idx-1)

	// pos := old.kvPos(idx)
	// storedKey := old.getKey(idx)

	// if bytes.Equal(storedKey, key) != true {
	// 	panic("invalid Node Update, Key Doesn't Exist")
	// }
	// klen := binary.LittleEndian.Uint16(old.data[pos:])
	// // vlen := binary.LittleEndian.Uint16(old.data[pos+2:])
	// vlen := uint16(len(val))
	// copy(new.data[pos+4+klen:], val)
	// // copy(new.data[pos+2:], vlen)
	// binary.LittleEndian.PutUint16(new.data[pos+2:], vlen)
}

func nodeAppendRange(new BNode, old BNode, dstNew uint16, srcOld uint16, n uint16) {
	assert(srcOld+n <= old.nkeys())
	assert(dstNew+n <= new.nkeys())

	if n == 0 {
		return
	}

	for i := uint16(0); i < n; i++ {
		new.setPtr(dstNew+i, old.getPtr(srcOld+i))
	}

	dstBegin := new.getOffset(dstNew)
	srcBegin := old.getOffset(srcOld)

	for i := uint16(1); i <= n; i++ {
		offset := dstBegin + old.getOffset(srcOld+i) - srcBegin
		new.setOffset(dstNew+i, offset)
	}

	begin := old.kvPos(srcOld)
	end := old.kvPos(srcOld + n)

	copy(new.data[new.kvPos(dstNew):], old.data[begin:end])
}

func nodeAppendKV(new BNode, idx uint16, ptr uint64, key []byte, val []byte) {
	new.setPtr(idx, ptr)
	pos := new.kvPos(idx)
	binary.LittleEndian.PutUint16(new.data[pos:], uint16(len(key)))
	binary.LittleEndian.PutUint16(new.data[pos+2:], uint16(len(val)))
	copy(new.data[pos+4:], key)
	copy(new.data[pos+4+uint16(len(key)):], val)
	new.setOffset(idx+1, new.getOffset(idx)+4+uint16((len(key)+len(val))))
}

func nodeInsert(tree *BTree, new BNode, node BNode, idx uint16, key []byte, val []byte) {
	kptr := node.getPtr(idx)
	knode := tree.get(kptr)
	tree.del(kptr)
	knode = treeInsert(tree, knode, key, val)
	nsplit, splited := nodeSplit3(knode)
	nodeReplaceKidN(tree, new, node, idx, splited[:nsplit]...)
}

func nodeSplit2(left BNode, right BNode, old BNode) {
	ncount := old.nkeys()
	copy(left.data[:2], old.data[:2])
	copy(right.data[:2], old.data[:2])

	found := uint16(ncount / 2)

	leftKeyCount := found
	rightKeyCount := ncount - found
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, leftKeyCount)
	copy(left.data[2:], buf)
	binary.LittleEndian.PutUint16(buf, rightKeyCount)
	copy(right.data[2:], buf)

	nodeAppendRange(left, old, 0, 0, leftKeyCount)
	nodeAppendRange(right, old, 0, leftKeyCount, rightKeyCount)
}

// divides node , if size of node becomes more than page size
func nodeSplit3(old BNode) (uint16, [3]BNode) {
	if old.nbytes() <= BTREE_PAGE_SIZE {
		old.data = old.data[:BTREE_PAGE_SIZE] //allocator expects 4096 bytes
		return 1, [3]BNode{old}
	}

	left := BNode{make([]byte, 2*BTREE_PAGE_SIZE)} //extra space for the operations
	right := BNode{make([]byte, BTREE_PAGE_SIZE)}
	nodeSplit2(left, right, old)

	if left.nbytes() <= BTREE_PAGE_SIZE {
		left.data = left.data[:BTREE_PAGE_SIZE]
		return 2, [3]BNode{left, right}
	}

	leftleft := BNode{make([]byte, BTREE_PAGE_SIZE)}
	middle := BNode{make([]byte, BTREE_PAGE_SIZE)}

	nodeSplit2(leftleft, middle, left)

	assert(leftleft.nbytes() <= BTREE_PAGE_SIZE)
	return 3, [3]BNode{leftleft, middle, right}
}

func nodeReplaceKidN(
	tree *BTree, new BNode, old BNode, idx uint16, kids ...BNode,
) {
	inc := uint16(len(kids))
	new.setHeader(BNODE_NODE, old.nkeys()+inc-1)
	nodeAppendRange(new, old, 0, 0, idx)

	for i, node := range kids {
		nodeAppendKV(new, idx+uint16(i), tree.new(node), node.getKey(0), nil)
	}
	nodeAppendRange(new, old, idx+inc, idx+1, old.nkeys()-(idx+1))
}
