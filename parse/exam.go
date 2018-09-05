package parse

import (
	"bytes"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"path/filepath"
	"strings"
	"gopkg.in/go-playground/validator.v9"
)

type Config struct{
	ModuleAbsPath  string `validate:"required"`
	CommandAbsPath string `validate:"required"`
	*moduleInfo
	*commandInfo
	GoDownloader   bool
	GoReleaser     bool
}

type commandInfo struct{
	PackageName string `validate:"required"`
	RootCommandVarName string `validate:"required"`
}

type moduleInfo struct{
	FullModuleName string `validate:"required"`
	Domain         string `validate:"required"`
	SubDomain      string `validate:"required"`
	ModuleName     string `validate:"required"`
}

var goValidator = validator.New()

func ParseCommandPackage(packagePath string, rootCommandVarName string) (*Config, error) {

	commandAbsPath, err := filepath.Abs(packagePath)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	config := &Config{
		CommandAbsPath: commandAbsPath,
	}

	eg := errgroup.Group{}
	eg.Go(func() error {

		eg.Go(func() error {
			modulePath, err := getModulePath(packagePath)
			if err != nil {
				return err
			}
			config.ModuleAbsPath = modulePath

			eg.Go(func() error {
				hasReleaser, err := directoryContainsFile(modulePath, ".goreleaser.yml")
				if err != nil {
					return err
				}
				config.GoReleaser = hasReleaser
				return nil
			})

			eg.Go(func() error {
				hasDownloader, err := directoryContainsFile(modulePath, ".godownloader.yml")
				if err != nil {
					return err
				}
				config.GoDownloader = hasDownloader
				return nil
			})

			eg.Go(func() error {
				moduleInfo, err := getModuleInfo(modulePath)
				if err != nil {
					return err
				}
				config.moduleInfo = moduleInfo
				return nil
			})

			return nil

		})

		eg.Go(func() error {
			commandInfo, err := getCommandPackageInfo(packagePath, rootCommandVarName)
			if err != nil {
				return err
			}
			config.commandInfo = commandInfo
			return nil
		})

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	if err := goValidator.Struct(config); err != nil {
		return nil, err
	}

	return config, nil
}

func getCommandPackageInfo(commandPath, rootCommandVar string) (*commandInfo, error) {
	commandPath, err := filepath.Abs(commandPath)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, commandPath, nil, parser.AllErrors)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	if len(pkgs) != 1 {
		return nil, errors.New("Expected one package in the command directory.")
	}
	info := &commandInfo{
		RootCommandVarName: rootCommandVar,
	}

	for _, pkg := range pkgs {
		info.PackageName = pkg.Name
	}

	// Verify the existence of the root command.
	rootCommandFile := ""
	for _, pkg := range pkgs {
		for fileName, file := range pkg.Files {
			for _, decl := range file.Decls {

				fn, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}
				if fn.Tok != token.VAR {
					continue
				}
				for _, spec := range fn.Specs {
					v, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for _, n := range v.Names {
						if n.Name == rootCommandVar {
							rootCommandFile = fileName
						}
					}
				}
			}
		}
	}
	if rootCommandFile == "" {
		return nil, errors.New(fmt.Sprintf("Cobra root variable [%s] not found", rootCommandFile))
	}

	return info, nil

}

// getModulePath recursively searches upward from given path for a go.mod file, and errors out if not found.
func getModulePath(commandPath string) (string, error) {

	abs, err := filepath.Abs(commandPath)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}

	for {
		rel, err := filepath.Rel("/", abs)
		if err != nil {
			return "", errors.Wrap(err, 0)
		}
		if rel == "." {
			return "", errors.New(fmt.Sprintf(`
				The given command path [%s] appears to not be within a Go module, which is required. 
				No 'go.mod' file was found in the directory, or in an parent directory.
				Please read https://github.com/golang/go/wiki/Modules for more information about setting one up.
			`, rel))
		}
		rel = filepath.Join("/", rel)

		hasFile, err := directoryContainsFile(rel, "go.mod")
		if err != nil {
			return "", errors.Wrap(err, 0)
		}
		if hasFile {
			return rel, nil
		}
		abs += "/.."
	}

}

func getModuleInfo(modulePath string) (*moduleInfo, error) {
	b, err := ioutil.ReadFile(filepath.Join(modulePath, "go.mod"))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	a := bytes.Split(b, []byte("\n"))
	if len(a) == 0 || strings.TrimSpace(string(a[0])) == "" {
		return nil, errors.New("Golang module `go.mod` file is not formatted in expected way")
	}

	res := &moduleInfo{}

	res.FullModuleName = strings.TrimSpace(strings.Replace(string(a[0]), "module", "", 1))
	parts := strings.Split(res.FullModuleName, "/")
	if parts[0] == "" || len(parts) > 3 {
		return nil, errors.New(fmt.Sprintf("Package name [%s] in unexpected format", res.FullModuleName))
	}
	res.Domain = strings.ToLower(parts[0])
	if len(parts) > 1 {
		res.SubDomain = strings.ToLower(parts[1])
	}
	if len(parts) > 2 {
		res.ModuleName = strings.ToLower(parts[2])
	}

	logrus.WithField("moduleData", res).Infof("Module parsed for config data.")
	return res, nil
}

func directoryContainsFile(directoryPath, fileName string) (bool, error) {

	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return false, errors.WrapPrefix(err, fmt.Sprintf("Could not read module path [%s]", directoryPath), 0)
	}
	for _, f := range files {
		if fileName == f.Name() {
			return true, nil
		}
	}
	return false, nil
}


//func executeCommandToParse(modulePath string) error {
//
//	tempPath := path.Join(os.TempDir(), "cobrareadme-" + uuid.New().String())
//
//	// Copy to temp directory to ensure we do not affect
//	if err := shutil.CopyTree(modulePath, tempPath, &shutil.CopyTreeOptions{
//		Symlinks: false,
//		Ignore: func (_ string, files []os.FileInfo) []string {
//			res := []string{}
//			for _, file := range files {
//				if !isNecessaryFile(file.Name()) {
//					res = append(res, file.Name())
//				}
//			}
//			return res
//		},
//	}); err != nil {
//		return errors.Wrap(err, 0)
//	}
//
//	// Replace main file.
//	files, err := ioutil.ReadDir(path.Join(tempPath, "cmd"))
//	if err != nil {
//		return errors.Wrap(err, 0)
//	}
//
//	for _, file := range files {
//		if file.IsDir() {
//			continue
//		}
//		b, err := ioutil.ReadFile(path.Join(tempPath, "cmd", file.Name()))
//		if err != nil {
//			return errors.Wrap(err, 0)
//		}
//		if bytes.Index(b, []byte("func main("))
//	}
//
//	if err := ioutil.WriteFile(
//		path.Join(tempPath, "cobraReadmeStaticMainReplacement.go"),
//		staticassets.Bytes("main.go.txt"),
//		os.ModePerm,
//	); err != nil {
//		return errors.Wrap(err, 0)
//	}
//
//
//
//
//}

//func safeClose(c io.Closer, err *error) {
//	if closeErr := c.Close(); closeErr != nil && *err == nil {
//				*err = closeErr
//		}
//	}


//func isNecessaryFile(fileName string) bool {
//	for _, suffix := range []string{".go", ".mod", ".sum"} {
//		if strings.HasSuffix(fileName, suffix) {
//			return true
//		}
//	}
//	return false
//}
//
//func Parse(modulePath string) (resErr error) {
//
//	cfg, err := parseModule(modulePath)
//	if err != nil {
//		return err
//	}
//
//
//
//
//
//	options.UseGoreleaserWithGithub = options.UseGoreleaserWithGithub && strings.HasPrefix(options.PackageName, "github/")
//
//	//rootBuf := bytes.Buffer{}
//	//if err := doc.GenMarkdown(rootCmd, &rootBuf); err != nil {
//	//	return errors.Wrap(err, 0)
//	//}
//	//rootMkd := rootBuf.String()
//	//
//	//rootMkd = rootMkd[strings.LastIndex(rootMkd, "######"):]
//
//	buf := new(bytes.Buffer)
//	buf.WriteString("# " + rootCmd.Name() + "\n\n")
//
//	buf.WriteString("### " + rootCmd.Short + "\n\n")
//
//	buf.WriteString("#### " + rootCmd.Long + "\n\n")
//
//	buf.WriteString("### Installation \n")
//	if options.UsesGoreleaser && strings.HasPrefix(options.PackageName, "github/") {
//
//		buf.WriteString(
//			fmt.Sprintf(`
//				bash <(wget -s https://raw.githubusercontent.com/%s/master/godownloader.sh)`, strings.TrimPrefix(options.PackageName, "github/")),
//		)
//	}
//
//
//	name := rootCmd.CommandPath()
//
//	for _, subCmd := range rootCmd.Commands() {
//		subCmdBuf := bytes.Buffer{}
//		if err := doc.GenMarkdown(subCmd, &subCmdBuf); err != nil {
//			return errors.Wrap(err, 0)
//		}
//		subCmdMkd := subCmdBuf.String()
//		fmt.Println(subCmdMkd)
//
//	}
//	return nil
//}