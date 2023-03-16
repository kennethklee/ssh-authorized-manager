package cmd

import (
	"fmt"
	"os/exec"

	"github.com/kennethklee/ssh-authorized-manager/ssham/plugin"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

type Plugin struct {
	Module      string
	Version     string
	Replacement string
}

func NewBuilderCommand(app core.App) *cobra.Command {
	var isRun bool
	var withPlugins []string
	var runCommand = &cobra.Command{
		Use:   "builder",
		Short: "Build the app server with the configured plugins",
		Long: `Build the app server with the configured plugins.

This option requires the go toolchain to be installed.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			builder, err := plugin.NewBuilder(withPlugins...)
			if err != nil {
				return fmt.Errorf("failed to create builder: %w", err)
			}

			_, err = exec.LookPath("go")
			if err != nil {
				return fmt.Errorf("go toolchain not found: %w", err)
			}

			if isRun {
				return builder.Run(args...)
			} else {
				return builder.Compile(args...)
			}
		},
	}
	runCommand.Flags().StringArrayVar(&withPlugins, "with", []string{}, "Plugins to load (format: module@version[=replacement])")
	runCommand.Flags().BoolVar(&isRun, "run", false, "Run the app instead of building")

	return runCommand
}
