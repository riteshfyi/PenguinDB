package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"PenguinDB/storage"
)

const (
	// This intentionally crosses page-split boundaries so the current
	// multi-page insertion bug is reproducible during benchmark setup.
	benchmarkRecordCount = 1024
	benchmarkValueSize   = 128
)

// openBenchmarkDB creates an isolated database for one benchmark and arranges
// for it to be closed when the benchmark finishes.
func openBenchmarkDB(b *testing.B) *storage.KV {
	b.Helper()

	db := &storage.KV{
		Path: filepath.Join(b.TempDir(), "benchmark.db"),
	}
	if err := db.Open(); err != nil {
		b.Fatalf("open benchmark database: %v", err)
	}

	b.Cleanup(db.Close)
	return db
}

// populateBenchmarkDB creates a repeatable set of records before timing begins.
func populateBenchmarkDB(b *testing.B, db *storage.KV) ([][]byte, []byte) {
	b.Helper()

	keys := make([][]byte, benchmarkRecordCount)
	value := bytes.Repeat([]byte{'v'}, benchmarkValueSize)

	for i := range keys {
		keys[i] = []byte(fmt.Sprintf("key-%08d", i))
		if err := db.Set(keys[i], value); err != nil {
			b.Fatalf("populate record %d: %v", i, err)
		}

		// Validate after every insertion so a failure reports the operation that
		// first made an older record unreachable, rather than only the final state.
		for previous := 0; previous <= i; previous++ {
			got, ok := db.Get(keys[previous])
			if !ok {
				b.Fatalf(
					"after inserting %d records, record %d (%q) disappeared",
					i+1,
					previous,
					keys[previous],
				)
			}
			if !bytes.Equal(got, value) {
				b.Fatalf(
					"after inserting %d records, record %d (%q) has an unexpected value",
					i+1,
					previous,
					keys[previous],
				)
			}
		}
	}

	return keys, value
}

// BenchmarkKVGetWarm measures successful lookups while rotating through a
// database whose pages have already been touched during setup.
func BenchmarkKVGetWarm(b *testing.B) {
	db := openBenchmarkDB(b)
	keys, expectedValue := populateBenchmarkDB(b, db)

	if _, ok := db.Get(keys[0]); !ok {
		b.Fatal("preloaded key was not found")
	}

	b.ReportAllocs()
	b.SetBytes(int64(len(keys[0]) + len(expectedValue)))

	iteration := 0
	var result []byte
	var found bool

	for b.Loop() {
		key := keys[iteration%len(keys)]
		result, found = db.Get(key)
		iteration++
	}

	if !found {
		b.Fatal("benchmark key was not found")
	}
	if !bytes.Equal(result, expectedValue) {
		b.Fatal("benchmark returned an unexpected value")
	}
}

// BenchmarkKVGetMissing measures the cost of proving that a key is absent
// from an already-populated, warm database.
func BenchmarkKVGetMissing(b *testing.B) {
	db := openBenchmarkDB(b)
	_, _ = populateBenchmarkDB(b, db)
	missingKey := []byte("key-that-does-not-exist")

	if _, ok := db.Get(missingKey); ok {
		b.Fatal("missing key unexpectedly exists")
	}

	b.ReportAllocs()
	b.SetBytes(int64(len(missingKey)))

	var found bool
	for b.Loop() {
		_, found = db.Get(missingKey)
	}

	if found {
		b.Fatal("missing key unexpectedly found")
	}
}

// BenchmarkKVDurableUpdate measures an update through the public Set API.
// Set currently flushes database pages and syncs the database file, so this is
// a durable-update benchmark rather than a B-tree-only benchmark.
func BenchmarkKVDurableUpdate(b *testing.B) {
	db := openBenchmarkDB(b)
	key := []byte("benchmark-key")
	valueA := bytes.Repeat([]byte{'A'}, benchmarkValueSize)
	valueB := bytes.Repeat([]byte{'B'}, benchmarkValueSize)

	if err := db.Set(key, valueA); err != nil {
		b.Fatalf("create initial key: %v", err)
	}

	b.ReportAllocs()
	b.SetBytes(int64(len(key) + len(valueA)))

	iteration := 0
	for b.Loop() {
		value := valueA
		if iteration%2 == 1 {
			value = valueB
		}

		if err := db.Set(key, value); err != nil {
			b.Fatalf("update key: %v", err)
		}
		iteration++
	}
}
