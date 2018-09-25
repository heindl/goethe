// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package pkg_compile_test

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add float64 float64 [...float64]",
	Short: "add numbers",
	Long:  "A long description of adding numbers together.",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		b, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		fmt.Fprint(os.Stdout, a+b)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
