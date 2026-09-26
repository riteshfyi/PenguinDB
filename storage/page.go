package storage

import "fmt"

func writePages(db *KV) error {
	npages := int(db.page.flushed) + len(db.page.temp)

	if err := extendFile(db, npages); err != nil {
		return err
	}

	if err := extendMmap(db, npages); err != nil {
		return err
	}

	for i, page := range db.page.temp {
		ptr := db.page.flushed + uint64(i)
		copy(db.pageGet(ptr).data, page)
	}
	return nil
}

func flushPages(db *KV) error {
	if err := writePages(db); err != nil {
		return err
	}
	return syncPages(db)
}

func syncPages(db *KV) error {
	//flysh data first, then update the master
	if err := db.fp.Sync(); err != nil {
		return fmt.Errorf("fsync: %w", err)
	}
	db.page.flushed += uint64(len(db.page.temp))
	db.page.temp = db.page.temp[:0]

	if err := masterStore(db); err != nil {
		return err
	}

	if err := db.fp.Sync(); err != nil {
		return fmt.Errorf("fsync: %w", err)
	}

	return nil
}
