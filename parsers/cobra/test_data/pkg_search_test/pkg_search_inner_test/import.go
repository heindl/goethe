// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package pkg_search_inner_test

import (
	"github.com/spf13/cobra"
)

var ExportedTestCmd = cobra.Command{
	Use: "Testing 1 2 3",
}

func Execute() {
	ExportedTestCmd.Execute()
}
