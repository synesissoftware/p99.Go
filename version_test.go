// Copyright 2026 Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 6th July 2026
 * Updated: 20th August 2026
 */

package p99_test

import (
	"github.com/synesissoftware/p99.Go"

	"github.com/stretchr/testify/require"

	"testing"
)

const (
	Expected_VersionMajor uint16 = 0
	Expected_VersionMinor uint16 = 2
	Expected_VersionPatch uint16 = 1
	Expected_VersionAB    uint16 = 0xFFFF
)

func Test_Version_Elements(t *testing.T) {
	require.Equal(t, Expected_VersionMajor, p99.VersionMajor)
	require.Equal(t, Expected_VersionMinor, p99.VersionMinor)
	require.Equal(t, Expected_VersionPatch, p99.VersionPatch)
	require.Equal(t, Expected_VersionAB, p99.VersionAB)
}

func Test_Version(t *testing.T) {
	require.Equal(t, uint64(0x0000_0002_0001_FFFF), p99.Version())
}

func Test_VersionString(t *testing.T) {
	require.Equal(t, "0.2.1", p99.VersionString())
}
