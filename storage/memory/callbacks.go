package memory

import (
	"PenguinDB/storage/btree"
)

// callback for BTree, derefernece a pointer
func (db *KV) pageGet(ptr uint64) btree.BNode {
	start := uint64(0)
	for _, chunk := range db.mmap.chunks {
		end := start + uint64(len(chunk))/BTREE_PAGE_SIZE
		if ptr < end {
			offset := BTREE_PAGE_SIZE * (ptr - start)
			return btree.BNode{chunk[offset : offset+BTREE_PAGE_SIZE]}
		}
		start = end
	}
	panic("bad ptr")
}

func (db *KV) pageNew(node btree.BNode) uint64 {
	//TODO : reuse deallocated pages
	assert(len(node.data) <= BTREE_PAGE_SIZE)
	ptr := db.page.flushed + uint64(len(db.page.temp))
	db.page.temp = append(db.page.temp, node.data)
	return ptr
}

func (db *KV) pageDel(ptr uint64) {
	total := db.page.flushed + uint64(len(db.page.temp))
	assert(ptr < total)

	//WRONG_IMPLEMENTATION
	// if ptr < db.page.flushed {
	//remove from the mmap
	// offset := ptr
	// db.mmap.chunks = append(db.mmap.chunks[:offset], db.mmap.chunks[offset+1:]...) //remove the ith map in memory
	//no need to delete data, just mark it to the free list
	//TODO: add this node to freelist, once the freelist is created.
	// } else {
	//remove form the temp
	// offset := ptr - db.page.flushed
	// db.page.temp = append(db.page.temp[:offset], db.page.temp[offset+1:]...)
	// }

	//mark this to freelist to be reused later.

}
