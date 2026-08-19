// Copyright 2026 Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 6th July 2026
 * Updated: 6th July 2026
 */

// Demonstrates recording event durations and querying percentiles.
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/synesissoftware/p99.Go"
)

type simpleRng struct {
	state uint64
}

func newSimpleRng(seed uint64) *simpleRng {
	return &simpleRng{state: seed}
}

func (r *simpleRng) next() uint64 {
	r.state = r.state*6_364_136_223_846_793_005 + 1
	return r.state
}

func main() {
	tries := 100
	if v := os.Getenv("P99_TRIES"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse P99_TRIES value %q, defaulting to 100\n", v)
		} else {
			tries = n
		}
	}

	fmt.Printf("Running Histogram example with %d tries...\n", tries)

	histogram := p99.New()
	rng := newSimpleRng(12_345)

	for i := 0; i < tries; i++ {
		delayUs := (rng.next() % 1_000) + 1
		start := time.Now()
		time.Sleep(time.Duration(delayUs) * time.Microsecond)
		elapsed := time.Since(start)
		histogram.PushEventDuration(elapsed)
	}

	fmt.Println("\nHistogram printed via GoString():")
	fmt.Println(histogram.GoString())

	fmt.Println("\nPercentiles (approximated):")
	printPercentileFn("p50 (f64):", func() (bool, uint64) { return histogram.ValueAtPercentile(50.0) })
	printPercentileFn("p50 (integer):", histogram.ValueAtP50)
	printPercentileFn("p75 (integer):", histogram.ValueAtP75)
	printPercentileFn("p90 (integer):", histogram.ValueAtP90)
	printPercentileFn("p95 (integer):", histogram.ValueAtP95)
	printPercentileFn("p99 (integer):", histogram.ValueAtP99)
	printPercentileFn("p99.5 (integer):", histogram.ValueAtP99_5)
	printPercentileFn("p99.9 (integer):", histogram.ValueAtP99_9)
	printPercentileFn("p99.99 (integer):", histogram.ValueAtP99_99)
}

func printPercentileFn(label string, fn func() (bool, uint64)) {
	ok, v := fn()
	printPercentile(label, ok, v)
}

func printPercentile(label string, ok bool, v uint64) {
	if ok {
		fmt.Printf("  %-18s %d ns\n", label, v)
	} else {
		fmt.Printf("  %-18s nil ns\n", label)
	}
}
