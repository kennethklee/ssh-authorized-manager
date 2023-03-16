package ssham

import (
	"fmt"
	"os"

	"github.com/kennethklee/ssh-authorized-manager/ssham/cmd"
	_ "github.com/kennethklee/ssh-authorized-manager/ssham/migrations"
	"github.com/kennethklee/ssh-authorized-manager/ssham/routes"
	"github.com/kennethklee/ssh-authorized-manager/ssham/worker"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

var ProgramName = "SSH Authorized Manager"
var Version = "dev"

func Main() {
	var app = pocketbase.New()

	err := TriggerPluginsPreload(app)
	if err != nil {
		panic(err)
	}

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		fmt.Println(ProgramName, Version)

		// bold := color.New(color.Bold).Add(color.FgGreen)
		// bold.Println("> Config")
		// authConfig := auth.HeaderAuthConfigFromEnv()
		// authConfig.AdminLogin = true // Need this to manage users & servers
		// auth.InstallHeaderAuth(e.App, e.Router, auth.HeaderAuthConfigFromEnv())
		routes.Register(e.Router)
		RegisterHooks(e.App, HooksConfigFromEnv())

		err = TriggerPluginsServe(e)
		if err != nil {
			return err
		}

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
	app.RootCmd.AddCommand(cmd.NewBuilderCommand(app))

	err = TriggerPluginsLoad(app)
	if err != nil {
		panic(err)
	}

	if err := app.Start(); err != nil {
		panic(err)
	}
}
