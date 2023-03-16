package plugin

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"text/template"
)

type PluginModule struct {
	Module      string
	Version     string
	Replacement string
}

func (p PluginModule) String() string {
	return p.Module + "@" + p.Version
}

type Builder struct {
	modules []PluginModule
}

func NewBuilder(pluginArgs ...string) (*Builder, error) {
	plugins := []PluginModule{}
	for _, pluginStr := range pluginArgs {
		plugin, err := parsePlugin(pluginStr)
		if err != nil {
			return nil, err
		}

		plugins = append(plugins, plugin)
	}

	return &Builder{modules: plugins}, nil
}

func parsePlugin(plugin string) (result PluginModule, err error) {
	// module@version[=replacement]
	pluginPattern, err := regexp.Compile(`^(?P<module>.+)@(?P<version>[^=]+)(=(?P<replacement>.+))?$`)
	if err != nil {
		return
	}

	match := pluginPattern.FindStringSubmatch(plugin)
	if len(match) == 0 {
		err = fmt.Errorf("invalid plugin: %s", plugin)
		return
	}

	for i, name := range pluginPattern.SubexpNames() {
		switch name {
		case "module":
			result.Module = match[i]
		case "version":
			result.Version = match[i]
		case "replacement":
			result.Replacement = match[i]
		}
	}
	fmt.Println("[INFO]", "Plugin:", result.Module, result.Version, result.Replacement)
	return
}

func runGoMod(projectPath string, args ...string) error {
	args = append([]string{"mod"}, args...)
	goCmd := exec.Command("go", args...)
	goCmd.Dir = projectPath
	goCmd.Stdout = os.Stdout
	goCmd.Stderr = os.Stderr
	goCmd.Stdin = os.Stdin
	return goCmd.Run()
}

func runGoGet(projectPath string, args ...string) error {
	args = append([]string{"get"}, args...)
	goCmd := exec.Command("go", args...)
	goCmd.Dir = projectPath
	goCmd.Stdout = os.Stdout
	goCmd.Stderr = os.Stderr
	goCmd.Stdin = os.Stdin
	return goCmd.Run()
}

func (b *Builder) build() (string, error) {
	// Generate temp go project in /tmp
	projectPath, err := os.MkdirTemp("", "ssham")
	if err != nil {
		return "", fmt.Errorf("failed to create temp project path: %w", err)
	}
	fmt.Println("[INFO]", "Temp project path:", projectPath)

	// Create project `go mod init`
	err = runGoMod(projectPath, "init", "ssham")
	if err != nil {
		return projectPath, fmt.Errorf("failed to create temp project: %w", err)
	}

	// Go get plugins `go get module@version`
	for _, module := range b.modules {
		err := runGoGet(projectPath, module.String())
		if err != nil {
			return projectPath, fmt.Errorf("failed to get plugin %s: %w", module.String(), err)
		}
	}

	// Go replace plugins `go mod edit -replace module@version=replacement`
	for _, module := range b.modules {
		if module.Replacement != "" {
			err := runGoMod(projectPath, "edit", "-replace", module.Module+"@"+module.Version+"="+module.Replacement)
			if err != nil {
				return projectPath, fmt.Errorf("failed to replace plugin %s: %w", module.Module, err)
			}
		}
	}

	// Generate main.go (import main, plugins, and call main)
	tpl, err := template.New("main.go").Parse(mainGoTemplate)
	if err != nil {
		return projectPath, fmt.Errorf("failed to parse main.go template: %w", err)
	}

	mainGoFile, err := os.Create(projectPath + "/main.go")
	if err != nil {
		return projectPath, fmt.Errorf("failed to create main.go: %w", err)
	}
	defer mainGoFile.Close()

	err = tpl.Execute(mainGoFile, struct{ Plugins []PluginModule }{Plugins: b.modules})
	if err != nil {
		return projectPath, fmt.Errorf("failed to execute main.go template: %w", err)
	}

	// Go mod tidy `go mod tidy`
	err = runGoMod(projectPath, "tidy")
	if err != nil {
		return projectPath, fmt.Errorf("failed to tidy project: %w", err)
	}

	return projectPath, nil
}

func (b *Builder) Compile(buildArgs ...string) error {
	projectPath, err := b.build()
	if projectPath != "" {
		defer func() {
			err := os.RemoveAll(projectPath)
			if err != nil {
				log.Fatalln("[ERROR]", "Failed to remove temp project path:", projectPath)
			}
		}()
	}
	if err != nil {
		return err
	}

	// Build `go build -o ssham main.go`
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	buildArgs = append(buildArgs, ".")
	buildArgs = append([]string{"build", "-o", path.Join(pwd, "ssham")}, buildArgs...)
	fmt.Println("[INFO] RUN:", "go", buildArgs)
	goBuildCmd := exec.Command("go", buildArgs...)
	goBuildCmd.Env = os.Environ()
	goBuildCmd.Dir = projectPath
	goBuildCmd.Stdout = os.Stdout
	goBuildCmd.Stderr = os.Stderr
	return goBuildCmd.Run()
}

func (b *Builder) Run(args ...string) error {
	projectPath, err := b.build()
	if projectPath != "" {
		defer func() {
			err := os.RemoveAll(projectPath)
			if err != nil {
				log.Fatalln("[ERROR]", "Failed to remove temp project path:", projectPath)
			}
		}()
	}
	if err != nil {
		// Wait for input to see error
		fmt.Println("[ERROR]", "Failed to tidy project. Press enter to continue")
		fmt.Scanln()
		return err
	}

	// Run `go run main.go args`
	goRunCmd := exec.Command("go", "run", "main.go")
	goRunCmd.Env = os.Environ()
	goRunCmd.Dir = projectPath
	goRunCmd.Args = append(goRunCmd.Args, args...)
	goRunCmd.Stdout = os.Stdout
	goRunCmd.Stderr = os.Stderr
	goRunCmd.Stdin = os.Stdin
	err = goRunCmd.Run()

	// Wait for input to see error
	if err != nil {
		fmt.Println("[ERROR]", "Failed to run project. Press enter to continue")
		fmt.Scanln()
	}
	return err
}

const mainGoTemplate = `package main

import (
	"github.com/kennethklee/ssh-authorized-manager/ssham"
	{{ range .Plugins }}
	_ "{{.Module}}"
	{{ end }}
)

func main() {
	ssham.Main()
}
`
