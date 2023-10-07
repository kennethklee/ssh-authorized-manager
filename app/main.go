package main

import (
	"fmt"
	"os"

	// _ "ssham/migrations"
	"ssham/cmd"
	"ssham/worker"

	"github.com/fatih/color"
	auth "github.com/kennethklee/pb-auth"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

var ProgramName = "SSH Authorized Manager"
var Version = "dev"
var dev = os.Getenv("APP_ENV") == "development"

func main() {
	fmt.Println(ProgramName, Version)

	if err := os.MkdirAll(os.TempDir(), 0755); err == nil {
		fmt.Println("Temp directory missing...created")
	}

	var app = pocketbase.New()
	app.RootCmd.Use = os.Args[0]
	app.RootCmd.Short = ProgramName
	app.RootCmd.Version = Version
	app.RootCmd.AddCommand(
		cmd.NewServersCommand(app),
		cmd.NewHealthCheckCmd(app),
	)

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{Automigrate: dev})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		RegisterRoutes(e.Router)
		RegisterHooks(e.App, HooksConfigFromEnv())

		bold := color.New(color.Bold).Add(color.FgGreen)
		bold.Println("> Config")
		authConfig := auth.HeaderAuthConfigFromEnv()
		authConfig.AdminLogin = true // Need this to manage users & servers
		auth.InstallHeaderAuth(e.App, e.Router, auth.HeaderAuthConfigFromEnv())

		worker.StartWorker(app)

		return nil
	})

	fmt.Println(os.Args)

	if err := app.Start(); err != nil {
		panic(err)
	}
}
