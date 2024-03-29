// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package cobra

import (
	"path"

	"github.com/heindl/goethe/utilities"
)

func Parse(mod *utilities.ModuleInfo) ([]*CommandData, error) {

	cmdVars, err := search(mod)
	if err != nil {
		return nil, err
	}

	res := []*CommandData{}
	for varImportPath := range cmdVars {
		cmd, err := compile(mod.AbsPathFromImportPath(path.Dir(varImportPath)), path.Base(varImportPath))
		if err != nil {
			return nil, err
		}
		res = append(res, cmd)
	}

	return res, nil
}
