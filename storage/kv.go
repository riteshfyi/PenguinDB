package storage

import (
	"os"
)

// represents a db file
type KV struct {
	Path string
	fp   *os.File
	tree BTree
	free Freelist
	mmap struct {
		file   int      //file size
		total  int      //mmap size
		chunks [][]byte //multiple mmaps, can be non-continous
	}
	page struct {
		flushed uint64            //database size in number of pages
		nfree   int               //total freelist nodes used
		nappend int               //extra appended nodes in the memory pages
		updates map[uint64][]byte //newloy allocated pages
	}
}
