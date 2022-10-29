package migrations

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		userOnlyRule := "@request.user.id=userId"
		publicKeys := models.Collection{
			Name: "publicKeys",
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     "userId",
					Required: true,
					Type:     schema.FieldTypeUser,
					Options: &schema.UserOptions{
						MaxSelect:     1,
						CascadeDelete: true,
					},
				},
				&schema.SchemaField{
					Name: "type",
					Type: schema.FieldTypeText,
					Options: &schema.TextOptions{
						Pattern: "^\\S*$",
					},
				},
				&schema.SchemaField{
					Name:     "publicKey",
					Required: true,
					Unique:   true,
					Type:     schema.FieldTypeText,
					Options: &schema.TextOptions{
						Pattern: "^\\S*$",
					},
				},
				&schema.SchemaField{
					Name: "comment",
					Type: schema.FieldTypeText,
				},
			),
			ListRule:   &userOnlyRule,
			ViewRule:   &userOnlyRule,
			CreateRule: &userOnlyRule,
			UpdateRule: &userOnlyRule,
			DeleteRule: &userOnlyRule,
		}
		dao.SaveCollection(&publicKeys)

		userOnlyServerRule := "@collection.userServers.userId=@request.user.id && @collection.userServers.serverId=id"
		servers := models.Collection{
			Name: "servers",
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name: "lastState",
					Type: schema.FieldTypeText,
				},
				&schema.SchemaField{
					Name:     "name",
					Required: true,
					Type:     schema.FieldTypeText,
				},
				&schema.SchemaField{
					Name:     "host",
					Required: true,
					Type:     schema.FieldTypeText,
				},
				&schema.SchemaField{
					Name:     "port",
					Required: true,
					Type:     schema.FieldTypeNumber,
				},
				&schema.SchemaField{
					Name:     "username",
					Required: true,
					Type:     schema.FieldTypeText,
				},
				&schema.SchemaField{
					Name: "usePassword",
					Type: schema.FieldTypeBool,
				},
				&schema.SchemaField{
					Name: "#password",
					Type: schema.FieldTypeText,
				},
				&schema.SchemaField{
					Name: "#privateKey",
					Type: schema.FieldTypeText,
				},
				&schema.SchemaField{
					Name: "#privateKeyPassphrase",
					Type: schema.FieldTypeText,
				},
				&schema.SchemaField{
					Name: "hostname",
					Type: schema.FieldTypeText,
				},
				&schema.SchemaField{
					Name: "hostKey",
					Type: schema.FieldTypeText,
				},
			),
			ListRule: &userOnlyServerRule,
			ViewRule: &userOnlyServerRule,
		}
		dao.SaveCollection(&servers)

		userOnlyServerLogRule := "@collection.userServers.userId=@request.user.id && @collection.userServers.serverId=serverId"
		serverLogs := models.Collection{
			Name: "serverLogs",
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     "serverId",
					Required: true,
					Type:     schema.FieldTypeRelation,
					Options: &schema.RelationOptions{
						MaxSelect:     1,
						CascadeDelete: true,
						CollectionId:  servers.Id,
					},
				},
				&schema.SchemaField{
					Name: "type",
					Type: schema.FieldTypeSelect,
					Options: &schema.SelectOptions{
						MaxSelect: 1,
						Values:    []string{"info", "warn", "error", "hostKey", "hostName"},
					},
				},
				&schema.SchemaField{
					Name:     "message",
					Required: true,
					Type:     schema.FieldTypeText,
				},
				&schema.SchemaField{
					Name: "payload",
					Type: schema.FieldTypeText,
				},
			),
			ListRule:   &userOnlyServerLogRule,
			ViewRule:   &userOnlyServerLogRule,
			DeleteRule: &userOnlyServerLogRule,
		}
		dao.SaveCollection(&serverLogs)

		userServers := models.Collection{
			Name: "userServers",
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     "userId",
					Required: true,
					Type:     schema.FieldTypeUser,
					Options: &schema.UserOptions{
						MaxSelect:     1,
						CascadeDelete: true,
					},
				},
				&schema.SchemaField{
					Name:     "serverId",
					Required: true,
					Type:     schema.FieldTypeRelation,
					Options: &schema.RelationOptions{
						MaxSelect:     1,
						CascadeDelete: true,
						CollectionId:  servers.Id,
					},
				},
				&schema.SchemaField{
					Name: "options",
					Type: schema.FieldTypeText,
				},
			),
			ListRule: &userOnlyRule,
			ViewRule: &userOnlyRule,
		}
		dao.SaveCollection(&userServers)

		return nil
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collectionNames := []string{"publicKeys", "userServers", "serverLogs", "servers"}
		for _, name := range collectionNames {
			collection, _ := dao.FindCollectionByNameOrId(name)
			if collection != nil {
				err := dao.DeleteCollection(collection)
				if err != nil {
					return fmt.Errorf("failed to delete collection %s: %s", name, err)
				}
			}
		}

		return nil
	})
}
