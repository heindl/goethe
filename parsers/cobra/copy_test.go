// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package cobra

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/phogolabs/parcello"
	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {

	// Test binary package existence.
	f, err := parcello.Dir(".").Open("parse.go")
	assert.NoError(t, err)
	assert.NoError(t, f.Close())
	f, err = parcello.Open("parse.go")
	assert.NoError(t, err)
	assert.NoError(t, f.Close())

	cmdPath, err := filepath.Abs("./test_data/pkg_compile_test")
	assert.NoError(t, err)

	dest, teardown, err := copyCommandDirectory(cmdPath, "RootCmd")
	assert.NoError(t, err)

	defer func() {
		err := teardown()
		assert.NoError(t, err)
		_, err = ioutil.ReadDir(dest)
		assert.Error(t, err)
	}()

	assert.True(t, strings.Contains(dest, tempFileSuffix))

	createdFiles, err := ioutil.ReadDir(dest)
	assert.NoError(t, err)

	foundFileNames := []string{}
	for _, f := range createdFiles {
		foundFileNames = append(foundFileNames, f.Name())
	}
	sort.Strings(foundFileNames)
	assert.Equal(t, foundFileNames, []string{"add.go", "gen_parse_cobra_readme.go", "main.go", "root.go", "subtract.go"})

}
