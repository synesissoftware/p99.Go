// Copyright 2026 Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 6th July 2026
 * Updated: 6th July 2026
 */

// Package p99 provides a low-cost, zero-allocation histogram for recording
// event durations in nanoseconds and querying performance percentiles (p50,
// p90, p99, p99.9, and so on).
//
// [Histogram] uses 64 logarithmic power-of-two buckets and performs linear
// interpolation within buckets to approximate percentile values. It is
// designed for high-frequency latency measurement with minimal overhead.
package p99

import (
	"fmt"
	"math"
	"math/bits"
	"strings"
	"time"
)

// Specifies the number of logarithmic buckets in a [Histogram].
const BucketCount = 64

// A low-cost, zero-allocation, fixed-size structure for recording event
// durations in nanoseconds and querying high-resolution percentiles.
type Histogram struct {
	eventCount     uint64
	eventTimeTotal uint64
	hasOverflowed  bool
	minEventTime   uint64
	maxEventTime   uint64
	hasMinMax      bool
	buckets        [BucketCount]uint64
}

// Returns a zero-initialized histogram.
func New() *Histogram {
	return &Histogram{}
}

// Resets the histogram to the equivalent of a newly constructed instance.
func (h *Histogram) Clear() {
	*h = Histogram{}
}

// Records an event with the given duration.
// The nanosecond value is truncated to uint64.
func (h *Histogram) PushEventDuration(d time.Duration) bool {
	return h.PushEventTimeNs(uint64(d.Nanoseconds()))
}

// Records an event with the given number of nanoseconds. Returns false if
// overflow has occurred or the running total would overflow.
func (h *Histogram) PushEventTimeNs(timeInNs uint64) bool {
	if !h.tryAddNsToTotalAndUpdateMinMax(timeInNs) {
		return false
	}

	h.eventCount++
	h.buckets[BucketIndex(timeInNs)]++
	return true
}

// Records an event with the given number of microseconds.
func (h *Histogram) PushEventTimeUs(timeInUs uint64) bool {
	hi, lo := bits.Mul64(timeInUs, 1_000)
	if hi != 0 {
		h.hasOverflowed = true
		return false
	}
	return h.PushEventTimeNs(lo)
}

// Records an event with the given number of milliseconds.
func (h *Histogram) PushEventTimeMs(timeInMs uint64) bool {
	hi, lo := bits.Mul64(timeInMs, 1_000_000)
	if hi != 0 {
		h.hasOverflowed = true
		return false
	}
	return h.PushEventTimeNs(lo)
}

// Records an event with the given number of seconds.
func (h *Histogram) PushEventTimeS(timeInS uint64) bool {
	hi, lo := bits.Mul64(timeInS, 1_000_000_000)
	if hi != 0 {
		h.hasOverflowed = true
		return false
	}
	return h.PushEventTimeNs(lo)
}

// Attempts to obtain the count of events in the bucket at index.
func (h *Histogram) BucketValue(index int) (bool, uint64) {
	if index < 0 || index >= BucketCount {
		return false, 0
	}
	return true, h.buckets[index]
}

// Returns a copy of all bucket counts.
func (h *Histogram) Buckets() [BucketCount]uint64 {
	return h.buckets
}

// Returns the number of recorded events.
func (h *Histogram) EventCount() uint64 {
	return h.eventCount
}

// Attempts to obtain the total event time in nanoseconds when no overflow
// has occurred.
func (h *Histogram) EventTimeTotal() (bool, uint64) {
	if h.hasOverflowed {
		return false, 0
	}
	return true, h.eventTimeTotal
}

// Returns the total event time in nanoseconds regardless of whether
// overflow has occurred.
func (h *Histogram) EventTimeTotalRaw() uint64 {
	return h.eventTimeTotal
}

// Reports whether an overflow has occurred.
func (h *Histogram) HasOverflowed() bool {
	return h.hasOverflowed
}

// Attempts to obtain the minimum event time observed.
func (h *Histogram) MinEventTime() (bool, uint64) {
	if !h.hasMinMax {
		return false, 0
	}
	return true, h.minEventTime
}

// Attempts to obtain the maximum event time observed.
func (h *Histogram) MaxEventTime() (bool, uint64) {
	if !h.hasMinMax {
		return false, 0
	}
	return true, h.maxEventTime
}

// Attempts to obtain the approximated duration in nanoseconds at the given
// percentile. Percentile is clamped to [0, 100].
func (h *Histogram) ValueAtPercentile(percentile float64) (bool, uint64) {
	if h.eventCount == 0 {
		return false, 0
	}

	p := math.Max(0, math.Min(100, percentile))

	if p <= 0 {
		return true, h.minEventTime
	}
	if p >= 100 {
		return true, h.maxEventTime
	}

	targetRank := float64(h.eventCount) * (p / 100.0)
	var accumulated uint64

	for i := 0; i < BucketCount; i++ {
		count := h.buckets[i]
		if count == 0 {
			continue
		}

		prevAccumulated := accumulated
		accumulated += count

		if float64(accumulated) >= targetRank {
			lower, upper := bucketRange(i)
			targetOffset := targetRank - float64(prevAccumulated)

			var rangeWidth float64
			if i == BucketCount-1 {
				rangeWidth = float64(math.MaxUint64 - lower)
			} else {
				rangeWidth = float64(upper - lower)
			}

			fraction := targetOffset / float64(count)
			interpolated := float64(lower) + (rangeWidth * fraction)
			value := uint64(math.Round(interpolated))
			return true, h.clampToMinMax(value)
		}
	}

	return true, h.maxEventTime
}

// Attempts to obtain the approximated duration at p50 (50th percentile).
func (h *Histogram) ValueAtP50() (bool, uint64) {
	return h.valueAtTargetRank(u64MulDiv(h.eventCount, 1, 2))
}

// Attempts to obtain the approximated duration at p75 (75th percentile).
func (h *Histogram) ValueAtP75() (bool, uint64) {
	return h.valueAtTargetRank(u64MulDiv(h.eventCount, 3, 4))
}

// Attempts to obtain the approximated duration at p90 (90th percentile).
func (h *Histogram) ValueAtP90() (bool, uint64) {
	return h.valueAtTargetRank(u64MulDiv(h.eventCount, 90, 100))
}

// Attempts to obtain the approximated duration at p95 (95th percentile).
func (h *Histogram) ValueAtP95() (bool, uint64) {
	return h.valueAtTargetRank(u64MulDiv(h.eventCount, 95, 100))
}

// Attempts to obtain the approximated duration at p99 (99th percentile).
func (h *Histogram) ValueAtP99() (bool, uint64) {
	return h.valueAtTargetRank(u64MulDiv(h.eventCount, 99, 100))
}

// Attempts to obtain the approximated duration at p99.5 (99.5th
// percentile).
func (h *Histogram) ValueAtP99_5() (bool, uint64) {
	return h.valueAtTargetRank(u64MulDiv(h.eventCount, 995, 1_000))
}

// Attempts to obtain the approximated duration at p99.9 (99.9th
// percentile).
func (h *Histogram) ValueAtP99_9() (bool, uint64) {
	return h.valueAtTargetRank(u64MulDiv(h.eventCount, 999, 1_000))
}

// Attempts to obtain the approximated duration at p99.99 (99.99th
// percentile).
func (h *Histogram) ValueAtP99_99() (bool, uint64) {
	return h.valueAtTargetRank(u64MulDiv(h.eventCount, 9_999, 10_000))
}

// Attempts to obtain the approximated duration at p99.999 (99.999th
// percentile).
func (h *Histogram) ValueAtP99_999() (bool, uint64) {
	return h.valueAtTargetRank(u64MulDiv(h.eventCount, 99_999, 100_000))
}

// Attempts to obtain the approximated duration at p99.9999 (99.9999th
// percentile).
func (h *Histogram) ValueAtP99_999_9() (bool, uint64) {
	return h.valueAtTargetRank(u64MulDiv(h.eventCount, 999_999, 1_000_000))
}

// Calculates the bucket index for a duration in nanoseconds.
func BucketIndex(timeInNs uint64) int {
	if timeInNs <= 1 {
		return 0
	}
	return bits.Len64(timeInNs) - 1
}

// Attempts to obtain the inclusive nanosecond range for the given bucket
// index.
func BucketRange(index int) (bool, uint64, uint64) {
	if index < 0 || index >= BucketCount {
		return false, 0, 0
	}
	lower, upper := bucketRange(index)
	return true, lower, upper
}

func bucketRange(index int) (lower, upper uint64) {
	if index == 0 {
		return 0, 1
	}
	lower = uint64(1) << index
	if index == BucketCount-1 {
		return lower, math.MaxUint64
	}
	return lower, (uint64(1) << (index + 1)) - 1
}

func (h *Histogram) tryAddNsToTotalAndUpdateMinMax(timeInNs uint64) bool {
	if h.hasOverflowed {
		return false
	}

	newTotal, carry := bits.Add64(h.eventTimeTotal, timeInNs, 0)
	if carry != 0 {
		h.hasOverflowed = true
		return false
	}
	h.eventTimeTotal = newTotal

	if !h.hasMinMax {
		h.minEventTime = timeInNs
		h.maxEventTime = timeInNs
		h.hasMinMax = true
		return true
	}

	if timeInNs < h.minEventTime {
		h.minEventTime = timeInNs
	}
	if timeInNs > h.maxEventTime {
		h.maxEventTime = timeInNs
	}
	return true
}

func (h *Histogram) valueAtTargetRank(targetRank uint64) (bool, uint64) {
	if h.eventCount == 0 {
		return false, 0
	}

	var accumulated uint64

	for i := 0; i < BucketCount; i++ {
		count := h.buckets[i]
		if count == 0 {
			continue
		}

		prevAccumulated := accumulated
		accumulated += count

		if accumulated >= targetRank {
			lower, upper := bucketRange(i)
			targetOffset := targetRank - prevAccumulated

			var value uint64
			if targetOffset == 0 {
				value = lower
			} else {
				var rangeWidth uint64
				if i == BucketCount-1 {
					rangeWidth = math.MaxUint64 - lower
				} else {
					rangeWidth = upper - lower
				}

				value = lower + u64MulDiv(rangeWidth, targetOffset, count)
			}

			return true, h.clampToMinMax(value)
		}
	}

	return true, h.maxEventTime
}

func (h *Histogram) clampToMinMax(value uint64) uint64 {
	if h.hasMinMax {
		if value < h.minEventTime {
			return h.minEventTime
		}
		if value > h.maxEventTime {
			return h.maxEventTime
		}
	}
	return value
}

func u64MulDiv(a, b, divisor uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	q, _ := bits.Div64(hi, lo, divisor)
	return q
}

// A compact debug representation of the histogram.
func (h *Histogram) String() string {
	return h.formatString(false)
}

// A verbose debug representation of the histogram.
func (h *Histogram) GoString() string {
	return h.formatString(true)
}

func (h *Histogram) formatString(verbose bool) string {
	var b strings.Builder

	hasTotal, total := h.EventTimeTotal()
	if verbose {
		fmt.Fprintf(&b, "&p99.Histogram{eventCount: %d, eventTimeTotal: ", h.eventCount)
		if hasTotal {
			fmt.Fprintf(&b, "%d", total)
		} else {
			b.WriteString("nil")
		}
		fmt.Fprintf(&b, ", hasOverflowed: %t", h.hasOverflowed)
		if h.hasMinMax {
			fmt.Fprintf(&b, ", minEventTime: %d, maxEventTime: %d", h.minEventTime, h.maxEventTime)
		} else {
			b.WriteString(", minEventTime: nil, maxEventTime: nil")
		}
		b.WriteString(", buckets: {")
		first := true
		for i := 0; i < BucketCount; i++ {
			if h.buckets[i] == 0 {
				continue
			}
			if !first {
				b.WriteString(", ")
			}
			first = false
			if verbose {
				fmt.Fprintf(&b, "2^%d: %d", i, h.buckets[i])
			} else {
				fmt.Fprintf(&b, "%d: %d", i, h.buckets[i])
			}
		}
		b.WriteString("}}")
		return b.String()
	}

	fmt.Fprintf(&b, "Histogram{n: %d, ∑: ", h.eventCount)
	if hasTotal {
		fmt.Fprintf(&b, "%d", total)
	} else {
		b.WriteString("nil")
	}
	fmt.Fprintf(&b, ", ∞: %t", h.hasOverflowed)
	if h.hasMinMax {
		fmt.Fprintf(&b, ", ↓: %d, ↑: %d", h.minEventTime, h.maxEventTime)
	} else {
		b.WriteString(", ↓: nil, ↑: nil")
	}
	b.WriteString(", b: {")
	first := true
	for i := 0; i < BucketCount; i++ {
		if h.buckets[i] == 0 {
			continue
		}
		if !first {
			b.WriteString(", ")
		}
		first = false
		fmt.Fprintf(&b, "%d: %d", i, h.buckets[i])
	}
	b.WriteString("}}")
	return b.String()
}
