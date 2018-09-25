package cobra

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

// CommandData is the aggregated documentation data from a Cobra command.
// Must remain in the parse.go file to be transferred into the cobra directory.
// This file should only import from cobra and the standard library.
type CommandData struct {
	// TODO: Eventually need to adjust the installation command based on multiple execution paths.
	// ExecutionRemotePath string `json:""`
	Name    string `json:""`
	Use     string `json:""`
	UseLine string `json:""`
	Long    string `json:""`
	Short   string `json:""`
	Example string `json:""`
	// HasExample determines if the command has example.
	HasExample bool `json:""`
	// Runnable determines if the command is itself runnable.
	Runnable bool `json:""`
	// HasSubCommands determines if the command has children commands.
	HasSubCommands bool           `json:""`
	SubCommands    []*CommandData `json:""`

	// IsAvailableCommand determines if a command is available as a non-help command
	// (this includes all non deprecated/hidden commands).
	IsAvailableCommand bool     `json:""`
	PersistentFlags    []string `json:""`
	LocalFlags         []string `json:""`
}

func (Ω *CommandData) formatFlagUsages(flags string) (res []string) {
	for _, line := range strings.Split(flags, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			res = append(res, line)
		}
	}
	return
}

func (Ω *CommandData) update(cmd *cobra.Command) {
	Ω.Name = cmd.Name()
	Ω.Use = cmd.Use
	Ω.UseLine = cmd.UseLine()
	Ω.Short = cmd.Short
	Ω.Long = cmd.Long
	Ω.Example = cmd.Example
	Ω.HasExample = cmd.HasExample()
	Ω.Runnable = cmd.Runnable()
	Ω.IsAvailableCommand = cmd.IsAvailableCommand()

	Ω.PersistentFlags = Ω.formatFlagUsages(cmd.PersistentFlags().FlagUsages())

	for _, localFlag := range Ω.formatFlagUsages(cmd.LocalFlags().FlagUsages()) {
		ex := false
		for _, persistentFlag := range Ω.PersistentFlags {
			ex = ex || localFlag == persistentFlag
		}
		if !ex {
			Ω.LocalFlags = append(Ω.LocalFlags, localFlag)
		}
	}

}

// Establish as variables to they can be mocked in tests.
var cobraReadmeErrorWriter io.Writer = os.Stderr
var cobraReadmeOutputWriter io.Writer = os.Stdout

// getCobraReadMeRootCommand is a hack to account for lack of pointers in root command variable.
func getCobraReadMeRootCommand(i interface{}) *cobra.Command {
	if reflect.ValueOf(i).Kind() != reflect.Ptr {
		cmd := i.(cobra.Command)
		return &cmd
	} else {
		return i.(*cobra.Command)
	}
}

func main() {

	parentCmd := getCobraReadMeRootCommand(cobraReadmeRootCmd)
	if parentCmd == nil {
		fmt.Fprintf(cobraReadmeErrorWriter, "root command [cobraReadmeRootCmd] is nil")
		return
	}
	parent := &CommandData{}
	parent.update(parentCmd)

	for i, childCmd := range parentCmd.Commands() {
		parent.SubCommands = append(parent.SubCommands, &CommandData{})
		parent.SubCommands[i].update(childCmd)
	}

	b, err := json.Marshal(parent)
	if err != nil {
		fmt.Fprintf(cobraReadmeErrorWriter, err.Error())
		return
	}
	fmt.Fprint(cobraReadmeOutputWriter, string(b))
}
