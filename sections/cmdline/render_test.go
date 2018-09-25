// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package cmdline

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

	expected := "\n## Command-line\n\n\n```bash\ngoethe [go_module_path]\n```\n```bash\n-p, --print   Print the template data to standard out.\n```\n"

	assert.Equal(t, expected, string(b))
}
