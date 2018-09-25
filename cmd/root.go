// Copyright 2018 Parker Heindl. All rights reserved.
// Licensed under the MIT License. See LICENSE.md in the
// project root for information.
//
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/heindl/goethe/render"
	"github.com/phogolabs/parcello"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//go:generate parcello -m long.md
// rootCmd represents the base command when called without any subcommands
var goetheRootCmd = &cobra.Command{
	Use:               "goethe [command_directory]",
	Short:             "Statically generate a github flavored README.md from a Go module.",
	DisableAutoGenTag: true,
	// TODO: Include Version from static parsing.
	RunE: func(cmd *cobra.Command, args []string) error {
		directoryPath := "."
		if len(args) > 0 {
			directoryPath = args[0]
		}
		if cmd.Flag("print") != nil && cmd.Flag("print").Value.String() == "true" {
			return render.Render(directoryPath, os.Stdout)
		}
		return writeReadme(directoryPath)
	},
}

func init() {
	descriptionFile, err := parcello.Open("long.md")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(descriptionFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	goetheRootCmd.Long = string(b)

	goetheRootCmd.PersistentFlags().BoolP("print", "p", false, "Print the template data to standard out.")
}

func writeReadme(cmdPath string) error {

	b := bytes.NewBuffer([]byte{})
	writer := bufio.NewWriter(b)
	if err := render.Render(cmdPath, writer); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return errors.WithStack(err)
	}

	readmeFilePath := filepath.Join(cmdPath, "README.md")

	// Important to not override custom data.
	if err := checkReadmeWriteSafety(readmeFilePath); err != nil {
		return err
	}

	if err := ioutil.WriteFile(readmeFilePath, b.Bytes(), 0700); err != nil {
		return errors.Wrapf(err, "could not write file %s", readmeFilePath)
	}
	return nil
}

// Note that this code is altered from traditional generated regexp
// search `^// Code generated .* DO NOT EDIT\.$` to account for markdown
// comments.
var goGeneratedRegexp = regexp.MustCompile(`(?m)^.*//.*Code generated .* DO NOT EDIT.*$`)

func checkReadmeWriteSafety(readmeFilePath string) error {

	if _, err := os.Stat(readmeFilePath); os.IsNotExist(err) {
		return nil
	}
	readmeBytes, err := ioutil.ReadFile(readmeFilePath)
	if err != nil {
		return errors.Wrapf(err, "could not read existing README.md file [%s]", readmeFilePath)
	}
	if !goGeneratedRegexp.Match(readmeBytes) {
		return errors.Errorf("a README.md file [%s] already exists that wasn't automatically generated. please move it to in order to use Goethe.", readmeFilePath)
	}
	return nil
}

// FindRootCommandVars adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := goetheRootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
