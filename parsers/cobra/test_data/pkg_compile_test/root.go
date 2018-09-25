// Copyright 2018 Parker Heindl. All rights reserved.
// Licensed under the MIT License. See LICENSE.md in the
// project root for information.
//
package pkg_compile_test

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Test that more comments can be added below
var RootCmd = &cobra.Command{
	Use:   "go-maths",
	Short: "run some maths",
	Long:  `A longer description of an example math function.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var cfgFile string

func init() {

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goethe-test.yaml)")
	RootCmd.PersistentFlags().BoolP("round", "r", false, "Round the result")
}
