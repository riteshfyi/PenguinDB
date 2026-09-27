package storage

import (
	"fmt"
	"os"
	"syscall"
)

func (db *KV) Open() error {
	fp, err := os.OpenFile(db.Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("OpenFile: %w", err)
	}

	db.fp = fp

	sz, chunk, err := mmapInit(db.fp)
	if err != nil {
		goto fail
	}

	db.mmap.file = sz
	db.mmap.total = len(chunk)
	db.mmap.chunks = [][]byte{chunk}

	db.tree.get = db.pageGet
	db.tree.new = db.pageNew
	db.tree.del = db.pageDel

	db.free.get = db.pageGet
	db.free.new = db.pageAppend
	db.free.use = db.pageUse

	err = masterLoad(db)

	if err != nil {
		goto fail
	}

	return nil

fail:
	db.Close()
	return fmt.Errorf("KV.Open: %w", err)
}

func (db *KV) Close() {
	for _, chunk := range db.mmap.chunks {
		err := syscall.Munmap(chunk)
		assert(err == nil)
	}

	_ = db.fp.Close()
}

func (db *KV) Get(key []byte) ([]byte, bool) {
	return db.tree.Get(key)
}

func (db *KV) Set(key []byte, val []byte) error {
	db.tree.Insert(key, val)
	return flushPages(db)
}

func (db *KV) Del(key []byte) (bool, error) {
	deleted := db.tree.Delete(key)
	return deleted, flushPages(db)
}
