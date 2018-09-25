## Motivation

Standard formatting and static analysis tools have helped make Go, not only a good language for cloud programming, but also for open source development.

Go code on Github is remarkably reliable to use. Yet developer speed and productivity - both the writer and reader - could be improved through better tools.

The README.md is generally the first interaction a developer has with new code, and it should be as consistent, reliable and clear as the underlying language.

#### Reliability
The worst way to keep documentation up to date is to ask a writer to update it in multiple places. So it should run automatically upon each release. And the text should be as tightly bound to the code as possible, whcih go already accomplishes with GoDocs.

#### Consistency
Requirements for good open source distributions are numerous and growing - licenses, binary distributions, usage examples, contributors - and are hard for developers to remember. But they are easy for computer programs to remember.

#### Clarity
Simple and to the point. Standard formatting achieves this well.

## Other Programs and data;
    - http://docopt.org/
    - https://developers.google.com/style/code-syntax

## Other Libraries
    https://github.com/segmentio/terraform-docs

## Future
- Include lines of code, number of subpackages, etc?
- Generate a libraries used section from analyzing imports.
    - Generate short summaries of each from readme generator, if github summary doesn't exist
    - Include versioning.
    - Try to ignore those that augment the standard library.
- Goal is to be a no configuration program, run right after you tag your release.
 - It takes best practices - through historical knowledge and new research - and applies them to your go code.
 - The reame is the first interface between you and new code, so it's very important to be consistent.
 - And so Goethe is also a linter. It makes recommendations to help you share your work better.
- Render the first command, or child command, under the description text.
    - Need to determine which child command to invoke by i guess counting the complexity.
- Add latest release, which is only useful to not have to switch over to another tab, but nevertheless
- Learn more about standards for man page documentation:
http://pubs.opengroup.org/onlinepubs/9699919799/
- Divide into sections, where each section has a template. This can be the be the first step toward plugins.
- Actually want to do some testing to see what helps developers understand new code the fastest. Which complex apis, this may not be that useful. But do programmers like one library with lots, of features, or many small ones? Do they prefer just examples with as little writing as possible?
- Note that no matter how complicated this description becomes, it has to work on itself. What good would it be if it only generated it's own documentation?
- Tests instead of examples.
    - Often examples do not work for code that developers forget about. This is painful for end users.
    - Can even generate tests from the underlying data.
- Generate related library list from imports (and possibly exports scraped from gGithub)
- Include Godoc exported values into the ReadMe. This requires the surface area of exported functions to be small.

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

