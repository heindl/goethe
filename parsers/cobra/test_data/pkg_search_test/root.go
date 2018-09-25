package pkg_search_test

import (
	"github.com/heindl/goethe/parsers/cobra/test_data/pkg_search_test/pkg_search_inner_test"
	"github.com/spf13/cobra"
	cmd2 "github.com/spf13/cobra/cobra/cmd"
)

var localRootCmd = &cobra.Command{
	Use: "testing 1 2 3",
}

//func Execute() {
//	pkg_search_inner_test.ExportedTestCmd.Execute()
//}

func main() {
	localRootCmd.Execute()
	pkg_search_inner_test.Execute()
}

func anotherTest() {
	cmd2.Execute()
}
