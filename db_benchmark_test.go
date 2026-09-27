package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"PenguinDB/storage"
)

const (
	// This is the number of key/value records loaded before read benchmarks.
	// It is not the number of physical B-tree nodes/pages.
	benchmarkRecordCount = 1000
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
	}

	// Validate the completed fixture once. The previous per-insert validation
	// was useful for locating the split bug but added unnecessary setup work.
	for i, key := range keys {
		got, ok := db.Get(key)
		if !ok {
			b.Fatalf("validate record %d (%q): key was not found", i, key)
		}
		if !bytes.Equal(got, value) {
			b.Fatalf("validate record %d (%q): unexpected value", i, key)
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
	b.ReportMetric(benchmarkRecordCount, "records")
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
	b.ReportMetric(benchmarkRecordCount, "records")
}

// BenchmarkKVDurableUpdate measures updates distributed across the same
// 1,000-record database used by the read benchmarks. Set currently flushes
// database pages and syncs the database file, so this measures durable updates
// rather than B-tree-only updates.
func BenchmarkKVDurableUpdate(b *testing.B) {
	db := openBenchmarkDB(b)
	keys, _ := populateBenchmarkDB(b, db)
	valueA := bytes.Repeat([]byte{'A'}, benchmarkValueSize)
	valueB := bytes.Repeat([]byte{'B'}, benchmarkValueSize)

	b.ReportAllocs()
	b.SetBytes(int64(len(keys[0]) + len(valueA)))

	iteration := 0
	var lastKey []byte
	var lastValue []byte

	for b.Loop() {
		keyIndex := iteration % len(keys)
		generation := iteration / len(keys)
		value := valueA
		if generation%2 == 1 {
			value = valueB
		}

		if err := db.Set(keys[keyIndex], value); err != nil {
			b.Fatalf("update key: %v", err)
		}
		lastKey = keys[keyIndex]
		lastValue = value
		iteration++
	}

	got, ok := db.Get(lastKey)
	if !ok {
		b.Fatal("last updated key was not found")
	}
	if !bytes.Equal(got, lastValue) {
		b.Fatal("last updated key has an unexpected value")
	}
	b.ReportMetric(benchmarkRecordCount, "records")
}
