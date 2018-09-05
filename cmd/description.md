A small script that generates a github.com flavored README.md file from a Go package and Cobra (github.com/spf13/cobra) command.

The intention is to keep the markdown file up to date on every release, write documentation in one place, and save time while maintaining standards.

It statically generates a Go file that wraps the Cobra Command, executes the command, and parses the command tree for all sub-commands.

Consider GoReleaser and GoDownloader for compiling cross-platform binary distributions. If a `goreleaser.yaml` file is found or a `godownloader.sh` script, the README.md will reflect simplified installation instructions.

The script has very specific expectations in the way the Go module is structured, which match the Cobra Generator (https://github.com/spf13/cobra#using-the-cobra-generator) pattern:

1. A Go Module (go version >= 1.11) has been initiated and the Cobra Command exists in a "./cmd" directory relative to module root.

    /a_new_module
        abc.go
        def.go
        go.mod
        /cmd
            uvw.go
            xyz.go

2. The Cobra root command variable is named "rootCmd", and not enclosed in any function.

```go
package main

var rootCmd = &cobra.Command{
	Use:   "subtract [number] [number]",
	Short: "Subtract two numbers",
}
```

3. Sub-commands are not added to the root command in `func main() {}`. This script statically toggles the main command, so it will not be called.

The Cobra Generator (https://github.com/spf13/cobra#using-the-cobra-generator) handles all of this for you so consider starting your next project with it.

