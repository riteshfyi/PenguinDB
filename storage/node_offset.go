package storage

import "encoding/binary"

func offsetPos(node BNode, idx uint16) uint16 {
	assert(1 <= idx && idx <= node.nkeys())
	return HEADER + 8*node.nkeys() + 2*(idx-1)
}

func (node BNode) getOffset(idx uint16) uint16 {

	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node.data[offsetPos(node, idx):]) //starting offset from the staert of the kvpart
}

func (node BNode) setOffset(idx uint16, offset uint16) {
	// assert(idx < node.nkeys())
	binary.LittleEndian.PutUint16(node.data[offsetPos(node, idx):], offset)
}
