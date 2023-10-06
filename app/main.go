package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	auth "github.com/kennethklee/pb-auth"
	"github.com/kennethklee/ssh-authorized-manager/ssham/cmd"
	_ "github.com/kennethklee/ssh-authorized-manager/ssham/migrations"
	"github.com/kennethklee/ssh-authorized-manager/ssham/worker"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

var ProgramName = "SSH Authorized Manager"
var Version = "dev"
var dev = os.Getenv("APP_ENV") == "development"

func main() {
	fmt.Println(ProgramName, Version)

	var app = pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		RegisterRoutes(e.Router)
		RegisterHooks(e.App, HooksConfigFromEnv())

		bold := color.New(color.Bold).Add(color.FgGreen)
		bold.Println("> Config")
		authConfig := auth.HeaderAuthConfigFromEnv()
		authConfig.AdminLogin = true // Need this to manage users & servers
		auth.InstallHeaderAuth(e.App, e.Router, auth.HeaderAuthConfigFromEnv())

		worker.StartWorker(app)

		// Debug
		// e.App.DB().LogFunc = log.Printf
		// e.App.DB().QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		// 	log.Printf("[%.2fms] Query SQL: %v", float64(t.Milliseconds()), sql)
		// }
		// e.App.DB().ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		// 	log.Printf("[%.2fms] Execute SQL: %v", float64(t.Milliseconds()), sql)
		// }

		return nil
	})

	app.RootCmd.Use = os.Args[0]
	app.RootCmd.Short = ProgramName
	app.RootCmd.Version = Version
	app.RootCmd.AddCommand(
		cmd.NewServersCommand(app),
		cmd.NewHealthCheckCmd(app),
	)

	if err := app.Start(); err != nil {
		panic(err)
	}
}
