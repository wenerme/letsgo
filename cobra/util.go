package wcobra

import "github.com/spf13/cobra"

func AppendRunE(run func(cmd *cobra.Command, args []string) error, f func(cmd *Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	if run == nil {
		run = func(cmd *cobra.Command, args []string) error {
			return f(Wrap(cmd), args)
		}
	} else {
		last := run
		run = func(cmd *cobra.Command, args []string) error {
			if err := last(cmd, args); err != nil {
				return err
			}
			return f(Wrap(cmd), args)
		}
	}
	return run
}
