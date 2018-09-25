// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package cmdline

import (
	"github.com/heindl/goethe/parsers/cobra"
	"github.com/heindl/goethe/utilities"
	"github.com/pkg/errors"
)

//go:generate parcello -m command.md

// Render returns the relevant README.md section.
func Render(mod *utilities.ModuleInfo) ([]byte, error) {

	commands, err := cobra.Parse(mod)
	if err != nil {
		return nil, err
	}

	if len(commands) == 0 {
		return nil, nil
	}

	if len(commands) > 1 {
		return nil, errors.Errorf("found %d root commands, though the system can only handle one right now. post a github issue and i'll work on this", len(commands))
	}

	return utilities.ExecuteTemplate("command.md", commands[0])
}
