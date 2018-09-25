// Copyright 2018 Parker Heindl. All rights reserved.
// Licensed under the MIT License. See LICENSE.md in the
// project root for information.
//
package utilities

import (
	"io/ioutil"

	"github.com/pkg/errors"
)

func DirectoryContainsFile(directoryPath, fileName string) (bool, error) {

	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return false, errors.Wrapf(err, "could not read directory %s", directoryPath)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		areEqual := f.Name() == fileName
		if areEqual { // || (!caseSensitive && strings.EqualFold(f.Name(), fileName)) {
			return true, nil
		}
	}
	return false, nil
}
