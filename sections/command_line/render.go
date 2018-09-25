package command_line

import (
	"github.com/heindl/goethe/parsers/cobra"
	"github.com/heindl/goethe/utilities"
	"github.com/pkg/errors"
)

//go:generate parcello -m command.md
func Render(mod *utilities.ModuleInfo) ([]byte, error) {

	commands, err := cobra.Parse(mod)
	if err != nil {
		return nil, err
	}

	if len(commands) == 0 {
		return nil, nil
	}

	if len(commands) > 1 {
		return nil, errors.Errorf("found %d root commands, though the system can only handle one right now. post a github issue and i'll work on this", len(commands))
	}

	return utilities.ExecuteTemplate("command.md", commands[0])
}
