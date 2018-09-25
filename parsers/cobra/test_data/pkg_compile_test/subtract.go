// Copyright 2018 Parker Heindl. All rights reserved.
// Licensed under the MIT License. See LICENSE.md in the
// project root for information.
//
package pkg_compile_test

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var subtractCmd = &cobra.Command{
	Use:   "subtract float64 float64 [...float64]",
	Short: "subtract numbers",
	Long:  `A long description of subtracting numbers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		b, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		fmt.Println(a - b)
		return nil
	},
}

func init() {
	subtractCmd.Flags().BoolP("absolute", "a", false, "Round the result to the absolute value")
	RootCmd.AddCommand(subtractCmd)
}
