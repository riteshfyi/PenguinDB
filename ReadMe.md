# PenguinDB

PenguinDB is a small persistent key-value store written in Go. I built it to understand what happens below a database API: how records are indexed, stored in fixed-size pages, written to disk, and recovered when the database is opened again.

The current version supports `Set`, `Get`, and `Delete` using a copy-on-write B-tree. Pages are accessed through `mmap`, and a free list allows released pages to be reused. Writes flush changed pages before updating the database metadata.

This is a learning project inspired by *Build Your Own Database From Scratch*. It is not intended for production use.

## Current features

- Persistent, disk-backed B-tree
- Insert, lookup, update, and delete operations
- Memory-mapped file access
- Free-list page reuse
- Database close and reopen support
- Correctness tests and repeatable benchmarks

## Benchmark

Measured on an Apple M4 Pro (`darwin/arm64`) with 1,000 records, 12-byte keys, 128-byte values, and a single goroutine. Results are the average of five runs.

| Workload | Latency | Throughput | Allocations |
| --- | ---: | ---: | ---: |
| Warm cached read | 225.88 ns/op | ~4.43M reads/s | 0 B/op, 0 allocs/op |
| Durable update (two `fsync`s) | 6.628 ms/op | ~151 writes/s | 20,480 B/op, 3 allocs/op |

These results are a baseline for the current implementation. They do not measure cold reads, concurrency, batching, WAL, or snapshots.

## Project status

The persistent KV store is complete as the first milestone, so development is paused for now. Possible future work includes range queries, transactions, concurrent snapshots, and write-ahead logging.
