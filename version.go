// Copyright 2026 Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 6th July 2026
 * Updated: 6th July 2026
 */

package p99

import "github.com/synesissoftware/ver2go"

const (
	VersionMajor uint16 = 0
	VersionMinor uint16 = 1
	VersionPatch uint16 = 0
	VersionAB    uint16 = 0x4001
	Version      uint64 = (uint64(VersionMajor) << 48) + (uint64(VersionMinor) << 32) + (uint64(VersionPatch) << 16) + (uint64(VersionAB) << 0)
)

var (
	versionString string = ver2go.CalcVersionString(VersionMajor, VersionMinor, VersionPatch, VersionAB)
)

// Returns the string form of the library version.
func VersionString() string {
	return versionString
}
