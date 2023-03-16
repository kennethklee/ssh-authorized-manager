package headerAuth

import (
	"github.com/fatih/color"
	auth "github.com/kennethklee/pb-auth"
	"github.com/kennethklee/ssh-authorized-manager/ssham"
	"github.com/pocketbase/pocketbase/core"
)

type Plugin struct{}

func (Plugin) Info() ssham.PluginInfo {
	return ssham.PluginInfo{
		Name:        "HeaderAuth",
		Version:     "0.0.1",
		Description: "Authenticate users based on HTTP headers",
	}
}

func (Plugin) OnPreload(app core.App) error {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		bold := color.New(color.Bold).Add(color.FgGreen)
		bold.Println("> Config")
		authConfig := auth.HeaderAuthConfigFromEnv()
		authConfig.AdminLogin = true // Need this to manage users & servers
		auth.InstallHeaderAuth(e.App, e.Router, auth.HeaderAuthConfigFromEnv())
		return nil
	})

	return nil
}

func (Plugin) OnLoad(app core.App) error {
	return nil
}

func (Plugin) OnServe(e *core.ServeEvent) error {
	return nil
}

func init() {
	ssham.RegisterPlugin(&Plugin{})
}
