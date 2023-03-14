package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kennethklee/ssh-authorized-manager/app/worker"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/spf13/cobra"
)

func NewServersCommand(app core.App) *cobra.Command {
	var serversCommand = &cobra.Command{
		Use:   "servers",
		Short: "Server commands",
	}

	serversCommand.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List servers",
		Run: func(cmd *cobra.Command, args []string) {
			serverCollection, _ := app.Dao().FindCollectionByNameOrId("servers")
			rows := []dbx.NullStringMap{}
			app.Dao().RecordQuery(serverCollection).All(&rows)
			servers := models.NewRecordsFromNullStringMaps(serverCollection, rows)

			w := tabwriter.NewWriter(os.Stdout, 20, 8, 1, ' ', 0)
			w.Write([]byte("Server ID\tName\tRemote Address\t\n"))
			for _, server := range servers {
				w.Write([]byte(fmt.Sprintf("%s\t%s\t%s\t\n", server.Id, server.GetString("name"), server.GetString("username")+"@"+server.GetString("host"))))
			}
			w.Flush()
		},
	})

	serversCommand.AddCommand(&cobra.Command{
		Use:   "sync [flags] serverId",
		Short: "Sync authorized_keys to server",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			server, err := app.Dao().FindRecordById("servers", args[0], nil)
			if err != nil {
				cmd.PrintErrln("Error:", err)
				cmd.Help()
				return
			}

			// Setup worker to sync authorized_keys on a server
			worker.SetApplication(app)
			work := worker.SyncServerWork{Server: server}
			if err := work.Execute(); err != nil {
				fmt.Println(err)
			}
		},
	})

	return serversCommand
}
