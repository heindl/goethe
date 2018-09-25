// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package header

import (
	"testing"

	"github.com/heindl/goethe/utilities"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {

	mod, err := utilities.ReadModule(".")
	assert.NoError(t, err)

	mod.PackagePathFilters = append(mod.PackagePathFilters, "_test")

	b, err := Render(mod)
	assert.NoError(t, err)

	expected := "\n[//]: # (Code generated by goethe DO NOT EDIT)\n# goethe\n\n*v0.0.2* | *September 25, 2018*\n#### Generate a github flavored README.md file from a Go module.\n\nGoethe parses your module directory for license, installation and command data, and generates a README.md file.\n\nIt is **very unstable** right now.\n\nGo1.11 modules are required, and this version only scans for a [Cobra Command](https://github.com/spf13/cobra) for documentation.\n\nFuture versions will support godoc, other command line helpers, and distribution and deployment tools. It will hopefully save time for both the reader and writer of open source code.\n"
	assert.Equal(t, expected, string(b))
}
