package storage

// callback for BTree, derefernece a pointer
func (db *KV) pageGet(ptr uint64) BNode {
	page, ok := db.page.updates[ptr]

	if ok {
		assert(page != nil)
		return BNode{page}
	}

	return pageGetMapped(db, ptr)
}

func pageGetMapped(db *KV, ptr uint64) BNode {
	start := uint64(0)
	for _, chunk := range db.mmap.chunks {
		end := start + uint64(len(chunk))/BTREE_PAGE_SIZE
		if ptr < end {
			offset := BTREE_PAGE_SIZE * (ptr - start)
			return BNode{data: chunk[offset : offset+BTREE_PAGE_SIZE]}
		}
		start = end
	}
	panic("bad ptr")
}

func (db *KV) pageNew(node BNode) uint64 {
	//TODO : reuse deallocated pages
	assert(len(node.data) <= BTREE_PAGE_SIZE)
	// ptr := db.page.flushed + uint64(len(db.page.temp))
	// db.page.temp = append(db.page.temp, node.data)
	ptr := uint64(0)
	if db.page.nfree < db.free.Total() {
		ptr = db.free.Get(db.page.nfree)
		db.page.nfree++
	} else {
		ptr = db.page.flushed + uint64(db.page.nappend)
		db.page.nappend++
	}
	db.page.updates[ptr] = node.data
	return ptr
}

func (db *KV) pageDel(ptr uint64) {
	// total := db.page.flushed + uint64(len(db.page.temp))
	// assert(ptr < total)
	total := db.page.flushed + uint64(db.page.nappend)

	assert(ptr < total)
	db.page.updates[ptr] = nil

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

func (db *KV) pageAppend(node BNode) uint64 {
	assert(len(node.data) <= BTREE_PAGE_SIZE)
	ptr := db.page.flushed + uint64(db.page.nappend)
	db.page.nappend++
	db.page.updates[ptr] = node.data
	return ptr
}

// callback for FreeList, reuse a page.
func (db *KV) pageUse(ptr uint64, node BNode) {
	db.page.updates[ptr] = node.data
}
