# p99.Go <!-- omit in toc -->

Low-cost generation of performance percentiles (p50, p90, p99, p99.9, etc.).

![Language](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![GitHub release](https://img.shields.io/github/v/release/synesissoftware/p99.Go.svg)](https://github.com/synesissoftware/p99.Go/releases/latest)
[![Last Commit](https://img.shields.io/github/last-commit/synesissoftware/p99.Go)](https://github.com/synesissoftware/p99.Go/commits/master)
[![Go](https://github.com/synesissoftware/p99.Go/actions/workflows/go.yml/badge.svg)](https://github.com/synesissoftware/p99.Go/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/synesissoftware/p99.Go.svg)](https://pkg.go.dev/github.com/synesissoftware/p99.Go)


## Table of Contents <!-- omit in toc -->

- [Introduction](#introduction)
- [How It Works](#how-it-works)
- [Performance \& Trade-offs](#performance--trade-offs)
  - [Performance Claims](#performance-claims)
  - [Trade-offs \& Sacrifices](#trade-offs--sacrifices)
- [Installation](#installation)
- [Components](#components)
  - [Constants](#constants)
  - [Functions](#functions)
  - [Structures](#structures)
    - [`Histogram`](#histogram)
      - [Minimal Example](#minimal-example)
- [Examples](#examples)
- [Project Information](#project-information)
  - [Where to get help](#where-to-get-help)
  - [Contribution guidelines](#contribution-guidelines)
  - [Dependencies](#dependencies)
    - [Development Dependencies](#development-dependencies)
  - [Related projects](#related-projects)
  - [License](#license)


## Introduction

**p99** is a lightweight, low-overhead library designed for generating real-time performance percentiles in high-frequency or latency-sensitive environments.

**p99.Go** is the **Go** implementation.


## How It Works

`Histogram` is a low-overhead, zero-allocation, fixed-size structure designed to track event durations (typically in nanoseconds) using 64 logarithmic buckets.

*   **Logarithmic Bucketing**: The bucket boundaries are spaced as powers of two:
    *   Bucket `0` represents `[0, 1]` nanoseconds;
    *   Bucket `1` represents `[2, 3]` nanoseconds;
    *   Bucket `2` represents `[4, 7]` nanoseconds;
    *   Bucket `i` represents `[2^i, 2^(i+1) - 1]` nanoseconds.
*   **Branchless Indexing**: Finding the correct bucket index for an incoming duration is extremely fast. It is computed in a few CPU instructions using `bits.Len64`.
*   **Linear Interpolation**: Percentile queries iterate through the buckets to find the target rank and perform linear interpolation within the matching bucket to approximate the exact percentile duration.


## Performance & Trade-offs

### Performance Claims

*   **Zero Allocation**: `Histogram` does not allocate memory on the heap during creation, event insertion, or percentile queries under normal operation. It is a compact structure that can reside entirely on the stack or be embedded in other structures.
*   **Ultra-Low Latency Insertion**: Recording a latency measurement (`PushEventTimeNs`) is designed for minimal overhead.
*   **Fast Queries**: Querying percentiles (such as `ValueAtP99()`) is designed to terminate early when events cluster in lower-indexed buckets.

### Trade-offs & Sacrifices

*   **Logarithmic Precision**: To achieve zero allocation and constant-time operations, `Histogram` sacrifices exact precision. It does not store individual event times.
*   **Approximation**: Percentile values are approximated using linear interpolation within the bucket boundaries. For very large values, the bucket width is wider, which leads to a wider approximation range.


## Installation

Install:

```bash
go get "github.com/synesissoftware/p99.Go"
```

Use:

```Go
import "github.com/synesissoftware/p99.Go"
```


## Components

### Constants

| Name | Value | Description |
|---|---|---|
| `BucketCount` | `64` | Number of logarithmic buckets in a `Histogram` |
| `VersionMajor` | `0` | Major version number |
| `VersionMinor` | `2` | Minor version number |
| `VersionPatch` | `0` | Patch version number |
| `VersionAB` | `ver2go.Release` (`0xFFFF`) | Final-release αβ-designator |


### Functions

| Function | Description |
|---|---|
| `New()` | Returns a zero-initialized histogram |
| `Version()` | Returns the packed 64-bit library version |
| `VersionString()` | Returns the string form of the library version |
| `BucketIndex(timeInNs uint64) int` | Calculates the bucket index for a duration |
| `BucketRange(index int) (bool, uint64, uint64)` | Returns the inclusive nanosecond range for a bucket |


### Structures

#### `Histogram`

A low-cost, zero-allocation, 64-bucket logarithmic histogram designed for recording event durations in nanoseconds and querying high-resolution percentiles.

##### Minimal Example

```Go
package main

import (
	"fmt"
	"time"

	"github.com/synesissoftware/p99.Go"
)

func main() {
	h := p99.New()

	h.PushEventTimeNs(150)
	h.PushEventTimeUs(5)
	h.PushEventTimeMs(10)
	h.PushEventDuration(250 * time.Nanosecond)

	fmt.Println("events:", h.EventCount())

	if ok, p99val := h.ValueAtP99(); ok {
		fmt.Printf("p99: %d ns\n", p99val)
	}
}
```


## Examples

Examples are provided in the `examples` directory, along with a markdown description for each. A detailed list of them is provided in [EXAMPLES.md](./EXAMPLES.md).


## Project Information

### Where to get help

[GitHub Page](https://github.com/synesissoftware/p99.Go "GitHub Page")


### Contribution guidelines

Defect reports, feature requests, and pull requests are welcome on https://github.com/synesissoftware/p99.Go.


### Dependencies

* [**ver2go**](https://github.com/synesissoftware/ver2go/)


#### Development Dependencies

* [**testify**](https://github.com/stretchr/testify/)


### Related projects

* [**p99**](https://github.com/synesissoftware/p99/)
* [**p99.Rust**](https://github.com/synesissoftware/p99.Rust/)


### License

**p99.Go** is released under the 3-clause BSD license. See [LICENSE](./LICENSE) for details.


<!-- ########################### end of file ########################### -->
