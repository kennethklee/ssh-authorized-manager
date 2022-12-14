package main

import (
	"fmt"
	"os"

	// "sshauthman/api/auth"
	"sshauthman/cmd"
	_ "sshauthman/migrations"
	"sshauthman/routes"
	"sshauthman/worker"

	"github.com/fatih/color"
	auth "github.com/kennethklee/pb-auth"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

var ProgramName = "SSH Authorized Manager"
var Version = "dev"

func main() {
	var app = pocketbase.New()
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		fmt.Println(ProgramName, Version)
		bold := color.New(color.Bold).Add(color.FgGreen)
		bold.Println("> Config")

		auth.InstallHeaderAuth(e.App, e.Router, auth.HeaderAuthConfigFromEnv())
		routes.Register(e.Router)
		RegisterHooks(e.App, HooksConfigFromEnv())

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
	app.RootCmd.AddCommand(cmd.NewServersCommand(app))
	app.RootCmd.AddCommand(cmd.NewAdminCommand(app))

	if err := app.Start(); err != nil {
		panic(err)
	}
}
