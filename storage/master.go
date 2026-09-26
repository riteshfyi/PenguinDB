package storage

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	DB_SIG = "MiyukiShiroganeW"
)

func masterLoad(db *KV) error {
	if db.mmap.file == 0 { //create a new master page on the first write
		db.page.flushed = 1 //reserved for master page
		return nil
	}

	data := db.mmap.chunks[0]
	root := binary.LittleEndian.Uint64(data[16:])
	used := binary.LittleEndian.Uint64(data[24:])

	if !bytes.Equal([]byte(DB_SIG), data[:16]) {
		return errors.New("bag signature")
	}

	bad := !(1 <= used && used <= uint64(db.mmap.file/BTREE_PAGE_SIZE)) || !(0 <= root && root < used)

	if bad {
		return errors.New("Bad master page")
	}
	db.tree.root = root
	db.page.flushed = used
	return nil
}

// update the master page
func masterStore(db *KV) error {
	var data [32]byte
	copy(data[0:], []byte(DB_SIG))
	binary.LittleEndian.PutUint64(data[16:], db.tree.root)
	binary.LittleEndian.PutUint64(data[24:], db.page.flushed)

	_, err := db.fp.WriteAt(data[:], 0)

	if err != nil {
		return fmt.Errorf("write master page : %w", err)
	}
	return nil
}
