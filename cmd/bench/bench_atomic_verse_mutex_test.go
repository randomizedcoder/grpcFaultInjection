package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkAtomicUint64(b *testing.B) {
	var counter atomic.Uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Add(1)
		}
	})
}

func BenchmarkMutexUint64(b *testing.B) {
	var counter uint64
	var mu sync.Mutex
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})
}
