package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cobra_readme -m ./path/to/module",
	Short: "Generate a README.md from a cobra (github.com/spf13/cobra) command.",
}

var modulePath string

func init() {
	b, err := ioutil.ReadFile("./description.md")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.Long = string(b)
	rootCmd.Flags().StringVarP(&modulePath, "module", "m", "", "Path to the Go module directory. We will make copy this into the temp directory and statically update it.")
	rootCmd.Flags().StringVarP(&modulePath, "module", "m", "", "Path to the Go module directory")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}