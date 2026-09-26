package storage

import (
	"fmt"
)

// extend file to atleast `npages`.
func extendFile(db *KV, npages int) error {
	filePages := db.mmap.file / BTREE_PAGE_SIZE
	if filePages >= npages {
		return nil
	}
	//increase file by 12.5%, instead of 100%, we save space
	for filePages < npages {
		inc := filePages / 8

		if inc < 1 {
			inc = 1
		}

		filePages += inc
	}

	fileSize := filePages * BTREE_PAGE_SIZE
	err := db.fp.Truncate(int64(fileSize))

	if err != nil {
		return fmt.Errorf("fallocate : %w", err)
	}

	db.mmap.file = fileSize
	return nil
}
