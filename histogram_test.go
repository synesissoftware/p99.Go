// Copyright 2026 Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 6th July 2026
 * Updated: 6th July 2026
 */

package p99_test

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/synesissoftware/p99.Go"
)

func Test_Histogram_String(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		h := p99.New()
		require.Equal(t, "Histogram{n: 0, ∑: 0, ∞: false, ↓: nil, ↑: nil, b: {}}", h.String())
	})

	t.Run("populated", func(t *testing.T) {
		h := p99.New()
		require.True(t, h.PushEventTimeNs(100))
		require.True(t, h.PushEventTimeNs(200))
		require.True(t, h.PushEventTimeNs(200))

		want := "Histogram{n: 3, ∑: 500, ∞: false, ↓: 100, ↑: 200, b: {6: 1, 7: 2}}"
		require.Equal(t, want, h.String())
	})
}

func Test_Histogram_GoString(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		h := p99.New()
		want := "&p99.Histogram{eventCount: 0, eventTimeTotal: 0, hasOverflowed: false, minEventTime: nil, maxEventTime: nil, buckets: {}}"
		require.Equal(t, want, h.GoString())
	})

	t.Run("populated", func(t *testing.T) {
		h := p99.New()
		require.True(t, h.PushEventTimeNs(100))
		require.True(t, h.PushEventTimeNs(10_000))
		require.True(t, h.PushEventTimeNs(10_001))

		want := "&p99.Histogram{eventCount: 3, eventTimeTotal: 20101, hasOverflowed: false, minEventTime: 100, maxEventTime: 10001, buckets: {2^6: 1, 2^13: 2}}"
		require.Equal(t, want, h.GoString())
	})
}

func Test_Histogram_New(t *testing.T) {
	h := p99.New()

	require.Equal(t, uint64(0), h.EventCount())
	ok, total := h.EventTimeTotal()
	require.True(t, ok)
	require.Equal(t, uint64(0), total)
	require.Equal(t, uint64(0), h.EventTimeTotalRaw())
	require.False(t, h.HasOverflowed())

	ok, _ = h.MinEventTime()
	require.False(t, ok)
	ok, _ = h.MaxEventTime()
	require.False(t, ok)

	buckets := h.Buckets()
	for i := 0; i < p99.BucketCount; i++ {
		require.Equal(t, uint64(0), buckets[i])
		ok, v := h.BucketValue(i)
		require.True(t, ok)
		require.Equal(t, uint64(0), v)
	}

	ok, _ = h.BucketValue(p99.BucketCount)
	require.False(t, ok)
}

func Test_Histogram_BucketIndex(t *testing.T) {
	cases := []struct {
		in   uint64
		want int
	}{
		{0, 0},
		{1, 0},
		{2, 1},
		{3, 1},
		{4, 2},
		{7, 2},
		{8, 3},
		{15, 3},
		{16, 4},
		{31, 4},
		{1024, 10},
		{2047, 10},
		{1 << 63, 63},
		{math.MaxUint64, 63},
	}
	for _, tc := range cases {
		require.Equal(t, tc.want, p99.BucketIndex(tc.in))
	}
}

func Test_Histogram_BucketRange(t *testing.T) {
	cases := []struct {
		index     int
		wantOK    bool
		wantLower uint64
		wantUpper uint64
	}{
		{0, true, 0, 1},
		{1, true, 2, 3},
		{2, true, 4, 7},
		{3, true, 8, 15},
		{4, true, 16, 31},
		{10, true, 1024, 2047},
		{63, true, 1 << 63, math.MaxUint64},
		{64, false, 0, 0},
	}
	for _, tc := range cases {
		ok, lower, upper := p99.BucketRange(tc.index)
		require.Equal(t, tc.wantOK, ok)
		require.Equal(t, tc.wantLower, lower)
		require.Equal(t, tc.wantUpper, upper)
	}
}

func Test_Histogram_PUSH_EVENTS(t *testing.T) {
	h := p99.New()

	require.True(t, h.PushEventTimeNs(1))
	require.True(t, h.PushEventTimeNs(3))
	require.True(t, h.PushEventTimeUs(10))
	require.True(t, h.PushEventTimeMs(5))
	require.True(t, h.PushEventTimeS(2))
	require.True(t, h.PushEventDuration(100*time.Nanosecond))

	require.Equal(t, uint64(6), h.EventCount())
	require.False(t, h.HasOverflowed())

	ok, min := h.MinEventTime()
	require.True(t, ok)
	require.Equal(t, uint64(1), min)

	ok, max := h.MaxEventTime()
	require.True(t, ok)
	require.Equal(t, uint64(2_000_000_000), max)

	ok, total := h.EventTimeTotal()
	require.True(t, ok)
	require.Equal(t, uint64(2_005_010_104), total)

	buckets := h.Buckets()
	checks := map[int]uint64{0: 1, 1: 1, 6: 1, 13: 1, 22: 1, 30: 1}
	for idx, want := range checks {
		require.Equal(t, want, buckets[idx])
	}

	h.Clear()
	require.Equal(t, uint64(0), h.EventCount())
	ok, total = h.EventTimeTotal()
	require.True(t, ok)
	require.Equal(t, uint64(0), total)
}

func Test_Histogram_OVERFLOW(t *testing.T) {
	h := p99.New()

	require.True(t, h.PushEventTimeNs(math.MaxUint64))
	ok, total := h.EventTimeTotal()
	require.True(t, ok)
	require.Equal(t, uint64(math.MaxUint64), total)
	require.False(t, h.HasOverflowed())

	require.False(t, h.PushEventTimeNs(1))
	require.True(t, h.HasOverflowed())
	ok, _ = h.EventTimeTotal()
	require.False(t, ok)
	require.Equal(t, uint64(math.MaxUint64), h.EventTimeTotalRaw())
}

func Test_Histogram_PERCENTILES_EMPTY(t *testing.T) {
	h := p99.New()

	ok, _ := h.ValueAtPercentile(50.0)
	require.False(t, ok)
	ok, _ = h.ValueAtP50()
	require.False(t, ok)
	ok, _ = h.ValueAtP99()
	require.False(t, ok)
}

func Test_Histogram_PERCENTILES_SINGLE_EVENT(t *testing.T) {
	h := p99.New()
	require.True(t, h.PushEventTimeNs(100))

	for _, p := range []float64{0.0, 50.0, 99.0, 100.0} {
		ok, v := h.ValueAtPercentile(p)
		require.True(t, ok)
		require.Equal(t, uint64(100), v)
	}

	for _, fn := range []func() (bool, uint64){
		h.ValueAtP50,
		h.ValueAtP90,
		h.ValueAtP99,
		h.ValueAtP99_999_9,
	} {
		ok, v := fn()
		require.True(t, ok)
		require.Equal(t, uint64(100), v)
	}
}

func Test_Histogram_PERCENTILES_INTERPOLATION(t *testing.T) {
	h := p99.New()
	require.True(t, h.PushEventTimeNs(100))
	require.True(t, h.PushEventTimeNs(200))

	ok, p50 := h.ValueAtPercentile(50.0)
	require.True(t, ok)
	ok, p99 := h.ValueAtPercentile(99.0)
	require.True(t, ok)
	require.GreaterOrEqual(t, p50, uint64(100))
	require.LessOrEqual(t, p50, uint64(200))
	require.GreaterOrEqual(t, p99, uint64(100))
	require.LessOrEqual(t, p99, uint64(200))

	ok, v := h.ValueAtPercentile(0.0)
	require.True(t, ok)
	require.Equal(t, uint64(100), v)
	ok, v = h.ValueAtPercentile(100.0)
	require.True(t, ok)
	require.Equal(t, uint64(200), v)

	ok, p50i := h.ValueAtP50()
	require.True(t, ok)
	require.GreaterOrEqual(t, p50i, uint64(100))
	ok, p99i := h.ValueAtP99()
	require.True(t, ok)
	require.LessOrEqual(t, p99i, uint64(200))
}

func Test_Histogram_PERCENTILES_WIDE_RANGE(t *testing.T) {
	h := p99.New()
	values := []uint64{
		1,
		10,
		100,
		1_000,
		10_000,
		100_000,
		1_000_000,
		10_000_000,
		100_000_000,
		1_000_000_000,
		10_000_000_000,
	}
	for _, v := range values {
		require.True(t, h.PushEventTimeNs(v))
	}

	require.Equal(t, uint64(len(values)), h.EventCount())

	ok, min := h.MinEventTime()
	require.True(t, ok)
	require.Equal(t, uint64(1), min)
	ok, max := h.MaxEventTime()
	require.True(t, ok)
	require.Equal(t, uint64(10_000_000_000), max)

	ok, p50 := h.ValueAtP50()
	require.True(t, ok)
	ok, p75 := h.ValueAtP75()
	require.True(t, ok)
	ok, p90 := h.ValueAtP90()
	require.True(t, ok)
	ok, p95 := h.ValueAtP95()
	require.True(t, ok)
	ok, p99 := h.ValueAtP99()
	require.True(t, ok)
	ok, p99_5 := h.ValueAtP99_5()
	require.True(t, ok)
	ok, p99_9 := h.ValueAtP99_9()
	require.True(t, ok)
	ok, p99_99 := h.ValueAtP99_99()
	require.True(t, ok)
	ok, p99_999 := h.ValueAtP99_999()
	require.True(t, ok)
	ok, p99_999_9 := h.ValueAtP99_999_9()
	require.True(t, ok)

	ordered := []uint64{p50, p75, p90, p95, p99, p99_5, p99_9, p99_99, p99_999, p99_999_9}
	for i := 1; i < len(ordered); i++ {
		require.LessOrEqual(t, ordered[i-1], ordered[i])
	}
	require.GreaterOrEqual(t, p50, uint64(1))
	require.LessOrEqual(t, p99_999_9, uint64(10_000_000_000))
}

func Test_Histogram_PERCENTILES_MANY_EVENTS(t *testing.T) {
	h := p99.New()
	const count = 100_000

	for i := uint64(1); i <= count; i++ {
		require.True(t, h.PushEventTimeNs(i))
	}

	require.Equal(t, uint64(count), h.EventCount())

	ok, min := h.MinEventTime()
	require.True(t, ok)
	require.Equal(t, uint64(1), min)
	ok, max := h.MaxEventTime()
	require.True(t, ok)
	require.Equal(t, uint64(count), max)

	ok, p50 := h.ValueAtP50()
	require.True(t, ok)
	ok, p90 := h.ValueAtP90()
	require.True(t, ok)
	ok, p99 := h.ValueAtP99()
	require.True(t, ok)
	ok, p99_9 := h.ValueAtP99_9()
	require.True(t, ok)

	require.Equal(t, uint64(50_000), p50)
	require.Equal(t, uint64(100_000), p90)
	require.Equal(t, uint64(100_000), p99)
	require.Equal(t, uint64(100_000), p99_9)

	ok, p75 := h.ValueAtP75()
	require.True(t, ok)
	ok, p95 := h.ValueAtP95()
	require.True(t, ok)
	ok, p99_5 := h.ValueAtP99_5()
	require.True(t, ok)
	ok, p99_99 := h.ValueAtP99_99()
	require.True(t, ok)
	ok, p99_999 := h.ValueAtP99_999()
	require.True(t, ok)
	ok, p99_999_9 := h.ValueAtP99_999_9()
	require.True(t, ok)

	ordered := []uint64{p50, p75, p90, p95, p99, p99_5, p99_9, p99_99, p99_999, p99_999_9}
	for i := 1; i < len(ordered); i++ {
		require.LessOrEqual(t, ordered[i-1], ordered[i])
	}
}

func Test_Histogram_COMPARE_FLOAT_AND_INT_PERCENTILES(t *testing.T) {
	h := p99.New()
	for i := 1; i <= 10_000; i++ {
		val := uint64((i * i) % 1_000_000)
		require.True(t, h.PushEventTimeNs(val))
	}

	pairs := []struct {
		p    float64
		intF func() (bool, uint64)
	}{
		{50.0, h.ValueAtP50},
		{75.0, h.ValueAtP75},
		{90.0, h.ValueAtP90},
		{95.0, h.ValueAtP95},
		{99.0, h.ValueAtP99},
		{99.5, h.ValueAtP99_5},
		{99.9, h.ValueAtP99_9},
		{99.99, h.ValueAtP99_99},
		{99.999, h.ValueAtP99_999},
		{99.9999, h.ValueAtP99_999_9},
	}

	for _, pair := range pairs {
		ok, floatVal := h.ValueAtPercentile(pair.p)
		require.True(t, ok)
		ok, intVal := pair.intF()
		require.True(t, ok)
		require.True(t, approxEqual(floatVal, intVal, 0.01), "percentile %v: float=%d int=%d", pair.p, floatVal, intVal)
	}
}

func approxEqual(a, b uint64, tolerance float64) bool {
	if a == b {
		return true
	}
	var larger, smaller uint64
	if a > b {
		larger, smaller = a, b
	} else {
		larger, smaller = b, a
	}
	if smaller == 0 {
		return larger == 0
	}
	diff := float64(larger-smaller) / float64(smaller)
	return diff <= tolerance
}
