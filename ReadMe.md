# PenguinDB

## Read/write benchmark

Measured on an Apple M4 Pro (`darwin/arm64`) with 1,000 records, 12-byte keys, 128-byte values, and a single goroutine. Results are the average of five runs.

| Workload | Latency | Throughput | Allocations |
| --- | ---: | ---: | ---: |
| Warm cached read | 225.88 ns/op | ~4.43M reads/s | 0 B/op, 0 allocs/op |
| Durable update (two `fsync`s) | 6.628 ms/op | ~151 writes/s | 20,480 B/op, 3 allocs/op |

These numbers are a baseline for the current implementation; they do not measure cold reads, concurrency, batching, WAL, or snapshots.
