package storage

const (
	BTREE_MAX_KEY_SIZE = 1000
	BTREE_MAX_VAL_SIZE = 3000
	BTREE_PAGE_SIZE    = 4096
)

type BTree struct {
	//pointer to the tree's node
	root uint64
	get  func(uint64) BNode //get a page from the disk
	new  func(BNode) uint64 //allocate a page to the object
	del  func(uint64)       //free the corresponding page
}
