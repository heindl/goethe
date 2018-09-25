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
