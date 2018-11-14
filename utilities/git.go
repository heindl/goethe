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
	"sort"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func GitLatestTag(modPath string) (string, error) {

	repo, err := git.PlainOpen(modPath)
	if err != nil {
		return "", errors.Wrapf(err, "could not open git repo %s", modPath)
	}

	tagObjects, err := repo.Tags()
	if err != nil {
		return "", errors.Wrapf(err, "could not read tags %s", modPath)
	}
	defer tagObjects.Close()

	tags := []string{}

	if err := tagObjects.ForEach(func(r *plumbing.Reference) error {
		tags = append(tags, r.Name().Short())
		return nil
	}); err != nil {
		return "", errors.Wrapf(err, "could not iterate over tags %s", modPath)
	}

	if len(tags) == 0 {
		return "", nil
	}

	sort.Strings(tags)

	return tags[len(tags)-1], nil

}

var gitConfSectionLineRe = regexp.MustCompile(`^\[.*]\s*$`)
var gitConfUserSectionLineRe = regexp.MustCompile(`^\[\s*user\s*]\s*$`)

func GitUserName(modPath string) (string, error) {
	repo, err := git.PlainOpen(modPath)
	if err != nil {
		return "", errors.Wrapf(err, "could not open git repo %s", modPath)
	}
	conf, err := repo.Config()
	if err != nil || conf == nil {
		return "", errors.Wrapf(err, "could not get config %s", modPath)
	}

	// TODO: Consider organization as well as user.
	userSec := conf.Raw.Section("user")
	if userSec != nil && userSec.Option("name") != "" {
		return userSec.Option("name"), nil
	}

	// Local config doesn't have a name, so check global.
	user, err := user.Current()
	if err != nil {
		return "", errors.Wrapf(err, "could not get current user %s", modPath)
	}
	if user == nil {
		return "", nil
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
