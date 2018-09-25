// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package cobra

import (
	"encoding/json"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func compile(packageFilePath, rootCommandVarName string) (command *CommandData, responseError error) {
	if err := validateSystemGo(); err != nil {
		return nil, err
	}

	dest, teardownFunc, err := copyCommandDirectory(packageFilePath, rootCommandVarName)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := teardownFunc(); closeErr != nil && responseError == nil {
			responseError = closeErr
		}
	}()

	if err := copyStaticParseFile(rootCommandVarName, dest); err != nil {
		return nil, err
	}

	// FindRootCommandVars test with local go and capture json output.
	cmd := exec.Command("go", "run", path.Join(dest, "."))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "could not execute duplicated command package [%s]: %s", packageFilePath, string(out))
	}

	res := CommandData{}
	if err := json.Unmarshal(out, &res); err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal command from %s: %s", packageFilePath, string(out))
	}
	return &res, nil
}

func validateSystemGo() error {
	version, err := getSystemGoVersion()
	if err != nil {
		return err
	}
	if version < 1.11 {
		return errors.Errorf("Program requires system go version of at least 1.11, found %f", version)
	}
	return nil
}

func getSystemGoVersion() (float32, error) {
	cmd := exec.Command("go", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, errors.Wrap(err, "could not execute go command")
	}
	if len(out) == 0 {
		return 0, errors.New("Could not get Go version")
	}
	v, err := strconv.ParseFloat(strings.TrimLeft(strings.Split(string(out), " ")[2], "go"), 32)
	if err != nil {
		return 0, errors.Wrap(err, "could not parse `go version` output")
	}
	return float32(v), nil
}
