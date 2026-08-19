// Copyright 2026 Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
/*
 * Created: 19th August 2026
 * Updated: 19th August 2026
 */

package main

import (
	"fmt"

	"github.com/synesissoftware/p99.Go"
	ver2go "github.com/synesissoftware/ver2go"
)

func main() {
	fmt.Printf("p99 v%s\n", ver2go.CalcVersionString(p99.VersionMajor, p99.VersionMinor, p99.VersionPatch, p99.VersionAB))
	fmt.Printf("ver2go v%s\n", ver2go.CalcVersionString(ver2go.VersionMajor, ver2go.VersionMinor, ver2go.VersionPatch, ver2go.VersionAB))
}
