// Copyright 2026 Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 6th July 2026
 * Updated: 6th July 2026
 */

package p99_test

import (
	"testing"

	"github.com/synesissoftware/p99.Go"
)

func buildSequentialHistogram() *p99.Histogram {
	h := p99.New()
	for i := uint64(1); i <= 100_000; i++ {
		_ = h.PushEventTimeNs(i * 10)
	}
	return h
}

func buildWideRangeHistogram() *p99.Histogram {
	h := p99.New()
	state := uint64(12_345)
	for i := uint64(1); i <= 100_000; i++ {
		state = state*6_364_136_223_846_793_005 + 1
		val := (state % 10_000_000_000) + 1
		_ = h.PushEventTimeNs(val)
		_ = i
	}
	return h
}

func Benchmark_BucketIndex_SMALL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkIndexSink = p99.BucketIndex(1)
	}
}

func Benchmark_BucketIndex_LARGE(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkIndexSink = p99.BucketIndex(^uint64(0))
	}
}

func Benchmark_PushEventTimeNs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := p99.New()
		_ = h.PushEventTimeNs(12_345)
	}
}

func Benchmark_Clear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := p99.New()
		_ = h.PushEventTimeNs(100)
		_ = h.PushEventTimeNs(200)
		h.Clear()
	}
}

func Benchmark_ValueAtPercentile_99_SEQUENTIAL(b *testing.B) {
	h := buildSequentialHistogram()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = h.ValueAtPercentile(99.0)
	}
}

func Benchmark_ValueAtP99_SEQUENTIAL(b *testing.B) {
	h := buildSequentialHistogram()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = h.ValueAtP99()
	}
}

func Benchmark_ValueAtPercentile_50_WIDE_RANGE(b *testing.B) {
	h := buildWideRangeHistogram()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = h.ValueAtPercentile(50.0)
	}
}

func Benchmark_ValueAtP50_WIDE_RANGE(b *testing.B) {
	h := buildWideRangeHistogram()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = h.ValueAtP50()
	}
}

func Benchmark_ValueAtPercentile_75_WIDE_RANGE(b *testing.B) {
	h := buildWideRangeHistogram()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = h.ValueAtPercentile(75.0)
	}
}

func Benchmark_ValueAtP75_WIDE_RANGE(b *testing.B) {
	h := buildWideRangeHistogram()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = h.ValueAtP75()
	}
}

func Benchmark_ValueAtPercentile_90_WIDE_RANGE(b *testing.B) {
	h := buildWideRangeHistogram()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = h.ValueAtPercentile(90.0)
	}
}

func Benchmark_ValueAtP90_WIDE_RANGE(b *testing.B) {
	h := buildWideRangeHistogram()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = h.ValueAtP90()
	}
}

func Benchmark_ValueAtPercentile_99_WIDE_RANGE(b *testing.B) {
	h := buildWideRangeHistogram()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = h.ValueAtPercentile(99.0)
	}
}

func Benchmark_ValueAtP99_WIDE_RANGE(b *testing.B) {
	h := buildWideRangeHistogram()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = h.ValueAtP99()
	}
}

func Benchmark_ValueAtPercentile_99_99_WIDE_RANGE(b *testing.B) {
	h := buildWideRangeHistogram()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = h.ValueAtPercentile(99.99)
	}
}

func Benchmark_ValueAtP99_99_WIDE_RANGE(b *testing.B) {
	h := buildWideRangeHistogram()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = h.ValueAtP99_99()
	}
}

var benchmarkIndexSink int
