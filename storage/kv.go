package storage

import (
	"os"
)

// represents a db file
type KV struct {
	Path string
	fp   *os.File
	tree BTree
	mmap struct {
		file   int      //file size
		total  int      //mmap size
		chunks [][]byte //multiple mmaps, can be non-continous
	}
	page struct {
		flushed uint64   //database size in number of pages
		temp    [][]byte //newloy allocated pages
	}
}
