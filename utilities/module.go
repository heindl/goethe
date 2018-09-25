// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package utilities

import (
	"bytes"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type ModuleInfo struct {
	domain             string
	subDomain          string
	name               string
	localFilePath      string
	packages           []Pkg
	PackagePathFilters []string
}

var reTestPkgFilter = regexp.MustCompile(`.*_test`)
var reGitFilter = regexp.MustCompile(`.*.git`)

func (Ω *ModuleInfo) regExpPackageFilters() ([]*regexp.Regexp, error) {

	y := []*regexp.Regexp{
		reGitFilter,
	}
	for _, fltr := range Ω.PackagePathFilters {
		re, err := regexp.Compile(fltr)
		if err != nil {
			return nil, errors.Wrapf(err, "could not compile filter regexp `%s`", re)
		}
		y = append(y, re)
	}

	// if is not a testing environment, ignore directories that have _test suffix.
	if flag.Lookup("test.v") == nil {
		y = append(y, reTestPkgFilter)
	}
	return y, nil

}

func (Ω *ModuleInfo) AbsPathFromImportPath(impPath string) string {
	return filepath.Join(Ω.localFilePath, strings.TrimPrefix(impPath, Ω.RemotePath()))
}

func (Ω *ModuleInfo) Domain() string {
	return Ω.domain
}

func (Ω *ModuleInfo) FilePath() string {
	return Ω.localFilePath
}

func (Ω *ModuleInfo) SubDomain() string {
	return Ω.subDomain
}

func (Ω *ModuleInfo) Name() string {
	return Ω.name
}

func (Ω *ModuleInfo) RemotePath() string {
	s := []string{
		Ω.domain,
	}
	if Ω.subDomain != "" {
		s = append(s, Ω.subDomain)
	}
	s = append(s, Ω.name)
	return strings.Join(s, "/")
}

type Pkg struct {
	AstPkg     *ast.Package
	AbsPath    string
	ImportPath string
}

// Packages examines a module recursively for packages.
// If not run in a `go test` environment, it will ignore directories with `_test` suffix.
func (Ω *ModuleInfo) Packages() ([]Pkg, error) {

	reFilters, err := Ω.regExpPackageFilters()
	if err != nil {
		return nil, err
	}

	pkgCh := make(chan Pkg)
	defer close(pkgCh)
	go func() {
		for pkg := range pkgCh {
			Ω.packages = append(Ω.packages, pkg)
		}
	}()
	eg := errgroup.Group{}
	if err := filepath.Walk(Ω.FilePath(), func(_path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		for _, re := range reFilters {
			if re.MatchString(_path) {
				return nil
			}
		}

		rel, err := filepath.Rel(Ω.FilePath(), _path)
		if err != nil {
			return errors.Wrapf(err, "could not get relative path [%s, %s]", Ω.FilePath(), _path)
		}

		eg.Go(func() error {
			ts := token.NewFileSet()
			pkgs, err := parser.ParseDir(ts, _path, nil, parser.AllErrors)
			if err != nil {
				return errors.Wrapf(err, "could not parse %s", _path)
			}
			for _, pkg := range pkgs {
				pkgCh <- Pkg{
					ImportPath: path.Join(Ω.RemotePath(), rel),
					AstPkg:     pkg,
					AbsPath:    _path,
				}
			}
			return nil
		})
		return nil
	}); err != nil {
		return nil, err
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return Ω.packages, nil
}

func ReadModule(path string) (*ModuleInfo, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get absolute path %s", absPath)
	}

	if _, err := os.Stat(filepath.Join(absPath, "go.mod")); err == nil {
		return parseModFile(filepath.Join(absPath, "go.mod"))
	}
	// Mod file doesn't exist, so check to see if we are in a package directory.
	modulePath, err := getModulePath(absPath)
	if err != nil {
		return nil, err
	}
	return parseModFile(filepath.Join(modulePath, "go.mod"))
}

func parseModFile(modFilePath string) (*ModuleInfo, error) {

	modFileBytes, err := ioutil.ReadFile(modFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read go.mod file %s", modFilePath)
	}

	firstLineIndex := bytes.Index(modFileBytes, []byte("\n"))

	if firstLineIndex == -1 {
		return nil, errors.Errorf("unexpected mod.go format %s", string(modFilePath))
	}

	pathSections := bytes.Split(bytes.TrimPrefix(bytes.TrimSpace(modFileBytes[0:firstLineIndex]), []byte("module ")), []byte("/"))

	if len(pathSections[0]) == 0 || len(pathSections) > 3 {
		return nil, errors.Errorf("unexpected module name in go.mod file %s", modFilePath)
	}

	mod := &ModuleInfo{
		domain:        string(bytes.ToLower(pathSections[0])),
		localFilePath: filepath.Dir(modFilePath),
	}
	if len(pathSections) == 3 {
		mod.subDomain = string(bytes.ToLower(pathSections[1]))
		mod.name = string(bytes.ToLower(pathSections[2]))
	} else {
		mod.name = string(bytes.ToLower(pathSections[1]))
	}

	return mod, nil
}

// getModulePath recursively searches upward from given path for a go.mod file, and errors out if not found.
func getModulePath(path string) (string, error) {

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", errors.Wrap(err, path)
	}

	for {
		rel, err := filepath.Rel("/", abs)
		if err != nil {
			return "", errors.Wrap(err, abs)
		}
		if rel == "." {
			logrus.Errorf(`The given command path [%s] appears to not be within a Go module, which is necessary.
				No 'go.mod' file was found in the directory, or in an parent directory.
				Please read https://github.com/golang/go/wiki/Modules for more information about setting one up.
			`, rel)
			return "", errors.New("given path is not within a module")
		}
		rel = filepath.Join("/", rel)

		hasFile, err := DirectoryContainsFile(rel, "go.mod")
		if err != nil {
			return "", err
		}
		if hasFile {
			return rel, nil
		}
		abs += "/.."
	}

}
