// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package license

import (
	"path/filepath"
	"time"

	"github.com/heindl/goethe/utilities"
	"github.com/pkg/errors"
	"github.com/ryanuber/go-license"
)

//go:generate parcello -m license.md

// Render returns the relevant README.md section.
func Render(modInfo *utilities.ModuleInfo) ([]byte, error) {

	lcns, err := license.NewFromDir(modInfo.FilePath())
	if err != nil || !lcns.Recognized() {
		utilities.MakeRecommendation(`
			There is no license file in the root of your module.
			Without one, anyone using your code is at risk of a lawsuit. 
			Read more here: https://choosealicense.com and consider adding a LICENSE.md file to your module root.
		`)
		return nil, nil
	}

	var data = struct {
		Type       string
		File       string
		AuthorName string
		Year       int
	}{
		Type: lcns.Type,
		Year: time.Now().Year(),
	}

	data.File, err = filepath.Rel(modInfo.FilePath(), lcns.File)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get relative path between %s and %s", modInfo.FilePath(), lcns.File)
	}

	data.AuthorName, err = utilities.GitUserName(modInfo.FilePath())
	if err != nil {
		return nil, err
	}

	return utilities.ExecuteTemplate("license.md", data)

}
