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

	"github.com/stretchr/testify/require"
	"github.com/synesissoftware/p99.Go"
)

func Test_VersionString(t *testing.T) {
	require.Equal(t, "0.1.0-alpha1", p99.VersionString())
}
