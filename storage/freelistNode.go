package storage

import "encoding/binary"

const BNODE_FREE_LIST = 3
const FREE_LIST_HEADER = 4 + 8 + 8
const FREE_LIST_CAP = (BTREE_PAGE_SIZE - FREE_LIST_HEADER) / 8

func flnSize(node BNode) int {
	assert(FREE_LIST_HEADER <= len(node.data) && len(node.data) <= BTREE_PAGE_SIZE)

	size := binary.LittleEndian.Uint16(node.data[2:])
	return int(size)
}
func flnNext(node BNode) uint64 {
	assert(FREE_LIST_HEADER <= len(node.data) && len(node.data) <= BTREE_PAGE_SIZE)

	next := binary.LittleEndian.Uint64(node.data[8:])
	return next
}
func flnPtr(node BNode, idx int) uint64 {
	size := flnSize(node)

	assert(idx >= 0 && idx < size)
	return binary.BigEndian.Uint64(node.data[FREE_LIST_HEADER+8*idx:])
}
func flnSetPtr(node BNode, idx int, ptr uint64) {
	size := flnSize(node)
	assert(idx >= 0 && idx < size)
	binary.BigEndian.PutUint64(node.data[FREE_LIST_HEADER+8*idx:], ptr)
}
func flnSetHeader(node BNode, size uint16, next uint64) {
	assert(uint16(0) <= size && size <= uint16(FREE_LIST_CAP))
	binary.LittleEndian.PutUint16(node.data[2:], size)
	binary.LittleEndian.PutUint64(node.data[12:], next)
}
func flnSetTotal(node BNode, total uint64) {
	assert(FREE_LIST_HEADER <= len(node.data) && len(node.data) <= BTREE_PAGE_SIZE)

	binary.LittleEndian.PutUint64(node.data[4:], total)
}
