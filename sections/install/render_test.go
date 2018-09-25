// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package install

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

	assert.Equal(t, "\n## Install\n\n```bash\ngo get -u github.com/heindl/goethe\n```\n", string(b))
}
