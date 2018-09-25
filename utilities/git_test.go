// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package utilities

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGit(t *testing.T) {

	modInfo, err := ReadModule(".")
	assert.NoError(t, err)

	tag, err := GitLatestTag(modInfo.FilePath())
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(tag, "v"))
	assert.Len(t, strings.Split(tag, "."), 3)

	user, err := GitUserName(modInfo.FilePath())
	assert.NoError(t, err)
	assert.Equal(t, "Parker Heindl", user)
}
