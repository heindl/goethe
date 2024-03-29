// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package license

import (
	"testing"

	"github.com/heindl/goethe/utilities"
	"github.com/stretchr/testify/assert"
)

func TestRenderInstall(t *testing.T) {

	mod, err := utilities.ReadModule(".")
	assert.NoError(t, err)

	mod.PackagePathFilters = append(mod.PackagePathFilters, "_test")

	b, err := Render(mod)
	assert.NoError(t, err)

	assert.Equal(t, "\n##\nCopyright 2018 Parker Heindl. All rights reserved.\nUse of this source code is governed by the [MIT License](LICENSE.md).\n", string(b))
}
