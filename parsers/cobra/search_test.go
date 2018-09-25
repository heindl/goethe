// Copyright 2018 Parker Heindl. All rights reserved.
// Licensed under the MIT License. See LICENSE.md in the
// project root for information.
//
package cobra

import (
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heindl/goethe/utilities"
	"github.com/stretchr/testify/assert"
)

func TestSearcher(t *testing.T) {

	const pkgPathRelativeToModule = "/parsers/cobra/test_data/pkg_search_test"
	const modImportPath = "github.com/heindl/goethe"

	absPkgPath, err := filepath.Abs("./test_data/pkg_search_test")
	assert.NoError(t, err)

	_ = strings.TrimSuffix(absPkgPath, pkgPathRelativeToModule)

	mod, err := utilities.ReadModule(".")
	assert.NoError(t, err)
	mod.PackagePathFilters = append(mod.PackagePathFilters, "/cmd", "/pkg_compile_test")

	cmds, err := search(mod)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cmds))

	assert.Len(t, cmds[path.Join(modImportPath, pkgPathRelativeToModule, "localRootCmd")], 1)
	assert.Len(t, cmds[path.Join(modImportPath, pkgPathRelativeToModule, "pkg_search_inner_test", "ExportedTestCmd")], 1)
}

func TestSearcherInMain(t *testing.T) {

	mod, err := utilities.ReadModule(".")
	assert.NoError(t, err)

	mod.PackagePathFilters = append(mod.PackagePathFilters, "_test")

	cmds, err := search(mod)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(cmds))

	assert.Equal(t, cmds["github.com/heindl/goethe/cmd/goetheRootCmd"][0], "github.com/heindl/goethe/cmd")
}
