package cmd

import (
	"github.com/spf13/cobra"
)

// Command is cli command interface
type Command interface {
	// Init command
	Init()

	// CobraCmd
	CobraCmd() *cobra.Command

	// Name
	Name() string
}

// baseCommand
type baseCommand struct {
	command *cobra.Command
}

func (bc *baseCommand) Init() {
}

func (bc *baseCommand) CobraCmd() *cobra.Command {
	return bc.command
}

func (bc *baseCommand) Name() string {
	return bc.command.Name()
}

//AddCommand is add child command to the parent command
func (bc *baseCommand) AddCommand(child Command) {
	child.Init()
	childCmd := child.CobraCmd()
	childCmd.SilenceUsage = true
	childCmd.DisableFlagsInUseLine = true
	childCmd.SilenceErrors = true
	bc.CobraCmd().AddCommand(childCmd)
}
