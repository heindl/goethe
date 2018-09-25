// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package cobra

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/phogolabs/parcello"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

const tempFileSuffix = "_goethe_tmp_"

const rootCmdTemplateName = "cobraReadmeRootCmd"

var cobraReadmeRootCmd = &cobra.Command{} // stub for testing

//go:generate parcello -m parse.go

func copyStaticParseFile(rootCommandVarName string, outputDirectory string) error {

	parseTemplateFile, err := parcello.Open("parse.go")
	if err != nil {
		return errors.Wrap(err, "Static compilation of parse.go is missing.")
	}

	data, err := ioutil.ReadAll(parseTemplateFile)
	if err != nil {
		return errors.Wrap(err, "Could not read parse.go as template file")
	}

	data = bytes.Replace(data, []byte(rootCmdTemplateName), []byte(rootCommandVarName), -1)
	data = bytes.Replace(data, []byte("package cobra"), []byte("package main"), 1)

	if err := ioutil.WriteFile(path.Join(outputDirectory, "gen_parse_cobra_readme.go"), data, 0700); err != nil {
		return errors.Wrap(err, "could not write parse.go template file")
	}

	return nil
}

func copyCommandDirectory(commandAbsPath, rootCommandVar string) (newDirPath string, teardownFunction func() error, responseError error) {

	dest := strings.TrimRight(commandAbsPath, "/") + tempFileSuffix + strconv.FormatInt(time.Now().UnixNano(), 10)

	entries, err := ioutil.ReadDir(commandAbsPath)
	if err != nil {
		return "", nil, errors.Wrapf(err, "Could not read command file [%s]", commandAbsPath)
	}

	if len(entries) == 0 {
		return "", nil, errors.Errorf("CommandData directory [%s] is empty", commandAbsPath)
	}

	if err := os.Mkdir(dest, 0700); err != nil {
		if strings.Contains(err.Error(), "exists") {
			return "", nil, errors.Errorf("The last execution failed to remove the generated file %s. Please delete it.", dest)
		}
		return "", nil, errors.Wrapf(err, "Could not generate build directory: %s", dest)
	}

	teardownFunc := func() error {
		if err := os.RemoveAll(dest); err != nil {
			return errors.Wrapf(err, "Could not tear down build directory: %s", dest)
		}
		return nil
	}

	defer func() {
		if responseError != nil {
			teardownFunc()
		}
	}()

	eg := errgroup.Group{}
	eg.Go(func() error {
		for _, _entry := range entries {
			entry := _entry
			eg.Go(func() error {
				// Skip symlinks.
				if entry.Mode()&os.ModeSymlink != 0 {
					return nil
				}
				return copyFile(
					filepath.Join(commandAbsPath, entry.Name()),
					filepath.Join(dest, entry.Name()),
				)
			})
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return "", nil, err
	}
	if err := copyStaticParseFile(rootCommandVar, dest); err != nil {
		return "", nil, err
	}
	return dest, teardownFunc, nil
}

var rePackage = regexp.MustCompile(`(?m)^[\s]*?package[\s]+.*$`)
var reFuncMain = regexp.MustCompile(`(?m)^[\s]*?func[\s]+main[\s]*?\(\)[\s]*?{.*$`)

func copyFile(src, dst string) error {

	if !strings.HasSuffix(src, ".go") {
		return nil
	}

	b, err := ioutil.ReadFile(src)
	if err != nil {
		return errors.Wrapf(err, "Could not read file %s", src)
	}

	b = rePackage.ReplaceAll(b, []byte("package main"))
	b = reFuncMain.ReplaceAll(b, []byte("func _main() {"))

	if err := ioutil.WriteFile(dst, b, 0700); err != nil {
		return errors.Wrapf(err, "Could not write file %s", dst)
	}

	return nil

}
