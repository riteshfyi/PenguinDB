package storage

import (
	"unsafe"
)

type C struct {
	tree  BTree
	ref   map[string]string
	pages map[uint64]BNode
}

func NewC() *C {
	pages := map[uint64]BNode{}
	return &C{
		tree: BTree{
			get: func(ptr uint64) BNode {
				node, ok := pages[ptr]
				assert(ok)
				return node
			},

			new: func(node BNode) uint64 {
				assert(node.nbytes() <= BTREE_PAGE_SIZE)
				key := uint64(uintptr(unsafe.Pointer(&node.data[0])))
				pages[key] = node
				return key
			},

			del: func(ptr uint64) {
				_, ok := pages[ptr]
				assert(ok)
				delete(pages, ptr)
			},
		},
		ref:   map[string]string{},
		pages: pages,
	}
}

func (c *C) Add(key string, val string) {
	c.tree.Insert([]byte(key), []byte(val))
	c.ref[key] = val
}
func (c *C) Del(key string) bool {
	delete(c.ref, key)
	return c.tree.Delete([]byte(key))
}

func (c *C) Get(key string) (string, bool) {
	val, ok := c.tree.Get([]byte(key))
	return string(val), ok
}

func (c *C) getRef(key string) (string, bool) {
	val, ok := c.ref[key]
	return val, ok
}
