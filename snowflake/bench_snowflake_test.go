package snowflake

import (
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkGenerate(b *testing.B) {
	b.ReportAllocs()

	node, _ := NewNode(1, 1580601600000, 10, 12)
	start := time.Now()
	b.ResetTimer()
	var x int64
	var n int64
	for i := 0; i < b.N; i++ {
		x = node.Generate()
		n++
	}
	b.StopTimer()
	_ = x

	dur := time.Since(start).Microseconds()
	b.ReportMetric(float64(n*1000000/dur), "ids/s")
}

func BenchmarkGenerateBatch(b *testing.B) {
	b.ReportAllocs()
	node, _ := NewNode(1, 1580601600000, 10, 12)

	start := time.Now()
	b.ResetTimer()
	count := uint16(100)
	var x []int64
	var n int64
	for i := 0; i < b.N; i++ {
		x = node.GenerateBatch(count)
		n += int64(count)
	}
	b.StopTimer()
	_ = x
	dur := time.Since(start).Microseconds()
	b.ReportMetric(float64(n*1000000/dur), "ids/s")
}

func BenchmarkGenerate_Parallel(b *testing.B) {
	b.ReportAllocs()

	node, _ := NewNode(1, 1580601600000, 10, 12)

	start := time.Now()
	b.ResetTimer()
	var x int64
	var n int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x = node.Generate()
			atomic.AddInt64(&n, 1)
		}
	})
	b.StopTimer()
	_ = x

	dur := time.Since(start).Microseconds()
	b.ReportMetric(float64(n*1000000/dur), "ids/s")
}

func BenchmarkGenerateBatch_Parallel(b *testing.B) {
	b.ReportAllocs()
	node, _ := NewNode(1, 1580601600000, 10, 12)

	start := time.Now()
	b.ResetTimer()
	count := uint16(100)
	var x []int64
	var n int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x = node.GenerateBatch(count)
			atomic.AddInt64(&n, int64(count))
		}
	})
	b.StopTimer()
	_ = x
	dur := time.Since(start).Microseconds()
	b.ReportMetric(float64(n*1000000/dur), "ids/s")
}
