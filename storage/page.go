package storage

import "fmt"

func writePages(db *KV) error {
	freed := []uint64{}

	for ptr, page := range db.page.updates {
		if page == nil {
			freed = append(freed, ptr)
		}
	}
	db.free.Update(db.page.nfree, freed)

	npages := int(db.page.flushed) + (db.page.nappend)

	if err := extendFile(db, npages); err != nil {
		return err
	}

	if err := extendMmap(db, npages); err != nil {
		return err
	}

	for ptr, page := range db.page.updates {
		// ptr := db.page.flushed + uint64(i)
		if page != nil {
			copy(pageGetMapped(db, ptr).data, page)
		}
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
	db.page.flushed += uint64(db.page.nappend)
	clear(db.page.updates)
	db.page.nfree = 0
	db.page.nappend = 0

	if err := masterStore(db); err != nil {
		return err
	}

	if err := db.fp.Sync(); err != nil {
		return fmt.Errorf("fsync: %w", err)
	}

	return nil
}
