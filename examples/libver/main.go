// Copyright 2026 Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
/*
 * Created: 19th August 2026
 * Updated: 20th August 2026
 */

package main

import (
	"fmt"

	"github.com/synesissoftware/p99.Go"
	ver2go "github.com/synesissoftware/ver2go"
)

func main() {
	fmt.Printf("p99 v%s\n", p99.VersionString())
	fmt.Printf("ver2go v%s\n", ver2go.VersionString())
}
