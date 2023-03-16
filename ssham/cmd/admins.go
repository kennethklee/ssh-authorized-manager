package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/spf13/cobra"
)

func NewAdminCommand(app core.App) *cobra.Command {
	var adminCommand = &cobra.Command{
		Use:   "admins",
		Short: "Admin commands",
	}

	adminCommand.AddCommand(&cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List admins",
		Run: func(cmd *cobra.Command, args []string) {
			var admins []models.Admin
			err := app.Dao().AdminQuery().All(&admins)
			if err != nil {
				panic(err)
			}

			w := tabwriter.NewWriter(os.Stdout, 20, 8, 1, ' ', 0)
			w.Write([]byte("Admin ID\tEmail\tUpdated\t\n"))
			for _, admin := range admins {
				w.Write([]byte(fmt.Sprintf("%s\t%s\t%s\t\n", admin.Id, admin.Email, admin.Updated)))
			}
			w.Flush()
		},
	})

	adminCommand.AddCommand(&cobra.Command{
		Use:   "add [flags] email",
		Short: "Add admin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			admin := models.Admin{}
			admin.Email = args[0]
			admin.RefreshTokenKey()

			if err := app.Dao().SaveAdmin(&admin); err != nil {
				return err
			}
			fmt.Println("Admin added")
			return nil
		},
	})

	return adminCommand
}
