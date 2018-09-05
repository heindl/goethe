package parse

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)


func TestPackageExam(t *testing.T) {

	wd, err := os.Getwd()
	assert.NoError(t, err)

	thisPath, err := filepath.Abs(wd + "/..")
	assert.NoError(t, err)

	for _, pkg := range []string{"main_test_pkg", "cmd_test_pkg"} {

		testDirPath := filepath.Join(wd, pkg)

		conf, err := ParseCommandPackage(testDirPath, "rootCmd")
		assert.NoError(t, err)

		assert.False(t, conf.GoDownloader)
		assert.False(t, conf.GoReleaser)

		assert.Equal(t, conf.ModuleAbsPath, thisPath)

		assert.Equal(t, conf.Domain, "github.com")
		assert.Equal(t, conf.FullModuleName, "github.com/heindl/cobrareadme")
		assert.Equal(t, conf.SubDomain, "heindl")
		assert.Equal(t, conf.ModuleName, "cobrareadme")

		assert.Equal(t, conf.PackageName, strings.Split(pkg, "_")[0])
		assert.Equal(t, conf.RootCommandVarName, "rootCmd")
	}

	_, err = ParseCommandPackage("main_test_pkg", "brokenCmd")
	assert.Error(t, err)

}