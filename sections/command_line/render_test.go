// Copyright 2018 Parker Heindl. All rights reserved.
// Licensed under the MIT License. See LICENSE.md in the
// project root for information.
//
package command_line

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

	expected := "\n### Command-line\n### Usage\n```bash\ngoethe [command_directory]\n```\n##### Flags\n```bash\n-p, --print   Print the template data to standard out.\n```\n"

	assert.Equal(t, expected, string(b))
}
