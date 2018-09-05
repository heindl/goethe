package cobrareadme

import (
	"bytes"
	"github.com/go-errors/errors"
	"github.com/spf13/cobra"
)

type Command struct{
	Name string `json:"name"`
	Long string `json:"long"`
	Short string `json:"short"`
	Example string `json:"example"`
	FlagString string `json:"flagString"`
	Children []*Command `json:"command"`
}

func staticFunctionToParseCobraCommand(parentCmd *cobra.Command) (*Command, error) {

	if parentCmd == nil {
		return nil, errors.New("Invalid parent command")
	}

	parent := &Command{
		Name:    parentCmd.Use, // TODO: If use is null use package name.
		Short:   parentCmd.Short,
		Long:    parentCmd.Long,
		Example: parentCmd.Example,
	}
	parentFlagBuffer := bytes.NewBuffer([]byte{})
	parentCmd.PersistentFlags().SetOutput(parentFlagBuffer)
	parent.FlagString = parentFlagBuffer.String()

	for _, childCmd := range parentCmd.Commands() {
		child := &Command{
			Name:    childCmd.Use,
			Short:   childCmd.Short,
			Long:    childCmd.Long,
			Example: childCmd.Example,
		}
		childFlagBuffer := bytes.NewBuffer([]byte{})
		childCmd.LocalFlags().SetOutput(childFlagBuffer)
		child.FlagString = childFlagBuffer.String()
		parent.Children = append(parent.Children, child)
	}

	return parent, nil
}
