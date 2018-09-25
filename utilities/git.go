// Copyright 2018 Parker Heindl. All rights reserved.
// Licensed under the MIT License. See LICENSE.md in the
// project root for information.
//
package utilities

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
)

func GitLatestTag(modPath string) (string, error) {
	if err := validateSystemGit(); err != nil {
		return "", err
	}
	cmd := exec.Command(
		"git",
		fmt.Sprintf("--git-dir=%s", path.Join(modPath, ".git")),
		fmt.Sprintf("--work-tree=%s", modPath),
		"describe",
		"--tags")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "could not read git tags in %s", modPath)
	}
	return strings.TrimSpace(strings.Split(string(out), "\n")[0]), nil
}

func GitUserName(modPath string) (string, error) {

	if err := validateSystemGit(); err != nil {
		return "", err
	}

	cmd := exec.Command(
		"git",
		fmt.Sprintf("--git-dir=%s", path.Join(modPath, ".git")),
		fmt.Sprintf("--work-tree=%s", modPath),
		"config",
		"user.name")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "could not read git config data in %s", modPath)
	}
	return strings.TrimSpace(string(out)), nil
}

func validateSystemGit() error {
	cmd := exec.Command("git", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "could not check system git")
	}
	if len(out) == 0 {
		return errors.New("could not check system git")
	}
	return nil
}
