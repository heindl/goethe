// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package pkg_compile_test

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is an exported root command variable to test execution from an external package.
var RootCmd = &cobra.Command{
	Use:   "go-maths",
	Short: "run some maths",
	Long:  `A longer description of an example math function.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

var cfgFile string

func init() {

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goethe-test.yaml)")
	RootCmd.PersistentFlags().BoolP("round", "r", false, "Round the result")
}
