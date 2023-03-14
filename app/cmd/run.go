package cmd

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

type Plugin struct {
	Module      string
	Version     string
	Replacement string
}

// Parse plugin string into module, version, and replacement

func parsePlugin(plugin string) (module, version, replacement string, err error) {
	// module@version[=replacement]
	pluginPattern, err := regexp.Compile(`(?P<module>.+)@(?P<version>.+)(=(?P<replacement>.+))?`)
	if err != nil {
		return
	}

	match := pluginPattern.FindStringSubmatch(plugin)
	if len(match) == 0 {
		panic(fmt.Errorf("Invalid plugin: %s", plugin))
	}

	for i, name := range pluginPattern.SubexpNames() {
		switch name {
		case "module":
			module = match[i]
		case "version":
			version = match[i]
		case "replacement":
			replacement = match[i]
		}
	}
	return
}

func runGoMod(projectPath string, args ...string) error {
	args = append([]string{"mod"}, args...)
	goModCmd := exec.Command("go", args...)
	goModCmd.Dir = projectPath
	return goModCmd.Run()
}

func NewRunCommand(app core.App) *cobra.Command {
	var withPlugins []string
	var runCommand = &cobra.Command{
		Use:   "run",
		Short: "Run the server with the configured plugins",
		Long: `Run the server with the configured plugins.

This option requires the go toolchain to be installed.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse plugins module@version[=replacement]
			plugins := []Plugin{}
			for _, plugin := range withPlugins {
				module, version, replacement, err := parsePlugin(plugin)
				if err != nil {
					return err
				}

				plugins = append(plugins, Plugin{
					Module:      module,
					Version:     version,
					Replacement: replacement,
				})
			}

			// Generate temp go project in /tmp
			projectPath, err := ioutil.TempDir("", "ssham")
			if err != nil {
				return fmt.Errorf("failed to create temp project path: %w", err)
			}
			fmt.Println("[INFO]", "Temp project path:", projectPath)
			defer func() {
				err := os.RemoveAll(projectPath)
				if err != nil {
					log.Fatalln("[ERROR]", "Failed to remove temp project path:", projectPath)
				}
			}()

			// Create project `go mod init`
			err = runGoMod(projectPath, "init", "ssham")
			if err != nil {
				return fmt.Errorf("failed to create temp project: %w", err)
			}

			// Go get plugins `go get module@version`
			for _, plugin := range plugins {
				err := runGoMod(projectPath, "get", plugin.Module+"@"+plugin.Version)
				if err != nil {
					return fmt.Errorf("failed to get plugin %s@%s: %w", plugin.Module, plugin.Module, err)
				}
			}

			// Go replace plugins `go mod edit -replace module@version=replacement`
			for _, plugin := range plugins {
				if plugin.Replacement != "" {
					err := runGoMod(projectPath, "edit", "-replace", plugin.Module+"@"+plugin.Version+"="+plugin.Replacement)
					if err != nil {
						return fmt.Errorf("failed to replace plugin %s: %w", plugin.Module, err)
					}
				}
			}

			// Generate main.go (import main, plugins, and call main)
			tpl, err := template.New("main.go").Parse(mainGoTemplate)
			if err != nil {
				return fmt.Errorf("failed to parse main.go template: %w", err)
			}

			mainGoFile, err := os.Create(projectPath + "/main.go")
			if err != nil {
				return fmt.Errorf("failed to create main.go: %w", err)
			}
			defer mainGoFile.Close()

			err = tpl.Execute(mainGoFile, struct{ Plugins []Plugin }{Plugins: plugins})
			if err != nil {
				return fmt.Errorf("failed to execute main.go template: %w", err)
			}

			// Go mod tidy `go mod tidy`
			err = runGoMod(projectPath, "tidy")
			if err != nil {
				// Wait for input to see error
				fmt.Println("[ERROR]", "Failed to tidy project. Press enter to continue")
				fmt.Scanln()

				return fmt.Errorf("failed to tidy project: %w", err)
			}

			// Run `go run main.go args`
			goRunCmd := exec.Command("go", "run", "main.go")
			goRunCmd.Dir = projectPath
			goRunCmd.Args = append(goRunCmd.Args, args...)
			goRunCmd.Stdout = os.Stdout
			goRunCmd.Stderr = os.Stderr
			goRunCmd.Stdin = os.Stdin
			return goRunCmd.Run()
		},
	}
	runCommand.Flags().StringArrayVar(&withPlugins, "with", []string{}, "Plugins to load (format: module@version[=replacement])")

	return runCommand
}

const mainGoTemplate = `package main

import (
	ssham "github.com/kennethklee/ssh-authorized-manager/app"
	{{ range .Plugins }}
	_ {{.Module}} {{.Version}}"
	{{ end }}
)

func main() {
	ssham.Main()
}
`
