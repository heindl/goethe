// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package cobra

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/heindl/goethe/parsers/cobra/test_data/pkg_compile_test"
	"github.com/stretchr/testify/assert"
)

func expectedCommandOutput() *CommandData {
	return &CommandData{
		Name:               "go-maths",
		Use:                "go-maths",
		UseLine:            "go-maths",
		Long:               "A longer description of an example math function.",
		Short:              "run some maths",
		Example:            "",
		HasExample:         false,
		Runnable:           false,
		HasSubCommands:     false,
		IsAvailableCommand: true,
		PersistentFlags: []string{
			"--config string   config file (default is $HOME/.goethe-test.yaml)",
			"-r, --round           Round the result",
		},
		// TODO: LocalFlags should probably not be the same as persistent flags.
		LocalFlags: nil,
		SubCommands: []*CommandData{
			{
				Name:               "add",
				Use:                "add float64 float64 [...float64]",
				UseLine:            "go-maths add float64 float64 [...float64]",
				Long:               "A long description of adding numbers together.",
				Short:              "add numbers",
				Example:            "",
				HasExample:         false,
				Runnable:           true,
				HasSubCommands:     false,
				SubCommands:        nil,
				IsAvailableCommand: true,
				PersistentFlags:    nil,
				LocalFlags:         nil,
			},
			{
				Name:               "subtract",
				Use:                "subtract float64 float64 [...float64]",
				UseLine:            "go-maths subtract float64 float64 [...float64] [flags]",
				Long:               "A long description of subtracting numbers.",
				Short:              "subtract numbers",
				Example:            "",
				HasExample:         false,
				Runnable:           true,
				HasSubCommands:     false,
				SubCommands:        nil,
				IsAvailableCommand: true,
				PersistentFlags:    nil,
				LocalFlags:         []string{"-a, --absolute   Round the result to the absolute value"},
			},
		},
	}
}

func checkCommandData(t *testing.T, expected, actual *CommandData) {

	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Use, actual.Use)
	assert.Equal(t, expected.UseLine, actual.UseLine)
	assert.Equal(t, expected.Long, actual.Long)
	assert.Equal(t, expected.Short, actual.Short)
	assert.Equal(t, expected.Example, actual.Example)
	assert.Equal(t, expected.HasExample, actual.HasExample)
	assert.Equal(t, expected.Runnable, actual.Runnable)
	assert.Equal(t, expected.IsAvailableCommand, actual.IsAvailableCommand)
	assert.Equal(t, expected.PersistentFlags, actual.PersistentFlags)
	assert.Equal(t, expected.LocalFlags, actual.LocalFlags)
	assert.Equal(t, expected.HasSubCommands, actual.HasSubCommands)
	assert.Equal(t, len(expected.SubCommands), len(actual.SubCommands))
	if len(expected.SubCommands) > 0 {
		for i := range expected.SubCommands {
			checkCommandData(t, expected.SubCommands[i], actual.SubCommands[i])
		}
	}
}

func TestParser(t *testing.T) {

	cobraReadmeRootCmd = pkg_compile_test.RootCmd

	testStderr := bytes.NewBufferString("")
	testStdout := bytes.NewBufferString("")

	cobraReadmeErrorWriter = testStderr
	cobraReadmeOutputWriter = testStdout

	main()

	assert.Equal(t, "", testStderr.String())

	var cmdData CommandData
	assert.NoError(t, json.Unmarshal(testStdout.Bytes(), &cmdData))

	checkCommandData(t, expectedCommandOutput(), &cmdData)

}

// Hint: If parser test is passing, but compilation test is failing, try go:generate.
func TestIntegratedCompilation(t *testing.T) {

	cmdPath, err := filepath.Abs("./test_data/pkg_compile_test")
	assert.NoError(t, err)

	cmdData, err := compile(cmdPath, "RootCmd")
	assert.NoError(t, err)

	checkCommandData(t, expectedCommandOutput(), cmdData)

}
