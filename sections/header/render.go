// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package header

import (
	"time"

	"github.com/heindl/goethe/parsers/cobra"
	"github.com/heindl/goethe/utilities"
)

//go:generate parcello -m header.md

// Render returns the relevant README.md section.
func Render(modInfo *utilities.ModuleInfo) ([]byte, error) {

	versTag, err := utilities.GitLatestTag(modInfo.FilePath())
	if err != nil {
		return nil, err
	}

	var data = struct {
		ModuleName       string
		VersionTag       string
		FormattedDate    string
		ShortDescription string
		LongDescription  string
	}{
		ModuleName: modInfo.Name(),
		VersionTag: versTag,
		// TODO: Mock time in tests.
		FormattedDate: time.Now().Format("January 2, 2006"),
	}

	cmdData, err := cobra.Parse(modInfo)
	if err != nil {
		return nil, err
	}

	if len(cmdData) > 0 {
		data.ShortDescription = cmdData[0].Short
		data.LongDescription = cmdData[0].Long
	}

	return utilities.ExecuteTemplate("header.md", data)

}
