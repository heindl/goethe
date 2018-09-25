// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package utilities

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	git "gopkg.in/src-d/go-git.v4"
)

func GitLatestTag(modPath string) (string, error) {

	repo, err := git.PlainOpen(modPath)
	if err != nil {
		return "", errors.Wrapf(err, "could not open git repo %s", modPath)
	}

	tags, err := repo.Tags()
	if err != nil {
		return "", err
	}
	defer tags.Close()

	tag, err := tags.Next()
	if err != nil {
		return "", err
	}

	return tag.Name().Short(), nil
}

var gitConfSectionLineRe = regexp.MustCompile(`^\[.*]\s*$`)
var gitConfUserSectionLineRe = regexp.MustCompile(`^\[\s*user\s*]\s*$`)

func GitUserName(modPath string) (string, error) {
	repo, err := git.PlainOpen(modPath)
	if err != nil {
		return "", errors.Wrapf(err, "could not open git repo %s", modPath)
	}
	conf, err := repo.Config()
	if err != nil {
		return "", errors.Wrapf(err, "could not get config %s", modPath)
	}

	// TODO: Consider organization as well as user.
	userSec := conf.Raw.Section("user")
	if userSec != nil && userSec.Option("name") != "" {
		return userSec.Option("name"), nil
	}

	// Local config doesn't have a name, so check global.
	user, err := user.Current()
	if user == nil || err != nil {
		return "", err
	}

	gitConfFile := filepath.Join(user.HomeDir, ".gitconfig")
	b, err := ioutil.ReadFile(gitConfFile)
	if err != nil {
		return "", errors.Wrapf(err, "could not read .gitconfig file %s", gitConfFile)
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	inUserSection := false
	for scanner.Scan() {
		txt := scanner.Text()
		isSection := gitConfSectionLineRe.MatchString(txt)
		isUserSection := gitConfUserSectionLineRe.MatchString(txt)
		inUserSection = (inUserSection && !isSection) || (isSection && isUserSection)
		if inUserSection {
			values := strings.Split(txt, "=")
			if len(values) != 2 {
				continue
			}
			if strings.TrimSpace(values[0]) != "name" {
				continue
			}
			return strings.TrimSpace(values[1]), nil
		}
	}
	return "", nil
}
