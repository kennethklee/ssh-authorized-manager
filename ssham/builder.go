package ssham

import (
	"fmt"
	"os/exec"
	"regexp"
)

type PluginModule struct {
	Module      string
	Version     string
	Replacement string
}

type Builder struct {
	plugins []PluginModule
}

func (b *Builder) New(pluginArgs ...string) *Builder {

	return &Builder{}
}

func (b *Builder) parsePluginArg(plugin string) (module, version, replacement string, err error) {
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

func (b *Builder) runGoMod(projectPath string, args ...string) error {
	args = append([]string{"mod"}, args...)
	goModCmd := exec.Command("go", args...)
	goModCmd.Dir = projectPath
	return goModCmd.Run()
}

func (b *Builder) build() error {
	return nil
}

func (b *Builder) Compile() error {
	return b.build()
}

func (b *Builder) Run() error {
	return b.build()
}
