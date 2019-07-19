package kaa

import (
	"github.com/spf13/cobra"
	"os"
)

// A wrapper around the cobra.Command. This interface is meant for base commands that should not be called directly. If
// your command is meant to contain a Run method, use SubCmd.
type Cmd interface {
	Command() *cobra.Command
}

// A wrapper around the cobra.Command. This interface is meant for subcommands of a Cmd. Every SubCmd must implement a
// Runner method.
type SubCmd interface {
	Cmd
	Run(ctx Context) error
}

// The alias for the SubCmd.Run(...) method.
type Runner func(ctx Context) error

// The alias for the cobra.Command Run field.
type CobraRunner func(cmd *cobra.Command, args []string)

// Used to wrap the base cobra.Command Run field. This Handler can also be passed in multiple Runners which all share
// the same Context
func Handle(runners ...Runner) CobraRunner {
	return func(cmd *cobra.Command, args []string) {
		ctx := NewContext(cmd, args)
		for _, r := range runners {
			err := r(ctx)
			if err != nil {
				_ = cmd.Usage()
				os.Exit(1)
			}
		}
	}
}
