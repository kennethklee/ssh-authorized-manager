package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

// Auto generated migration with the most recent collections configuration.
func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `[
			{
				"id": "fqthnfcd9k791ur",
				"created": "2022-09-24 01:52:59.401Z",
				"updated": "2022-12-20 21:21:03.522Z",
				"name": "publicKeys",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "bga2vftx",
						"name": "userId",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"collectionId": "systemprofiles0",
							"cascadeDelete": true
						}
					},
					{
						"system": false,
						"id": "ahmf4ljk",
						"name": "type",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": "^\\S*$"
						}
					},
					{
						"system": false,
						"id": "cyzn4pll",
						"name": "publicKey",
						"type": "text",
						"required": true,
						"unique": true,
						"options": {
							"min": null,
							"max": null,
							"pattern": "^\\S*$"
						}
					},
					{
						"system": false,
						"id": "rayhtawe",
						"name": "comment",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"listRule": "@request.auth.id=userId",
				"viewRule": "@request.auth.id=userId",
				"createRule": "@request.auth.id=userId",
				"updateRule": "@request.auth.id=userId",
				"deleteRule": "@request.auth.id=userId",
				"options": {}
			},
			{
				"id": "0wy8xlddg0uekql",
				"created": "2022-09-24 01:52:59.402Z",
				"updated": "2022-12-20 21:21:03.514Z",
				"name": "servers",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "e7hily1x",
						"name": "lastState",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "9vrkwzfh",
						"name": "name",
						"type": "text",
						"required": true,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "hwmv2syu",
						"name": "host",
						"type": "text",
						"required": true,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "xfz7uwqi",
						"name": "port",
						"type": "number",
						"required": true,
						"unique": false,
						"options": {
							"min": null,
							"max": null
						}
					},
					{
						"system": false,
						"id": "jmui01gw",
						"name": "username",
						"type": "text",
						"required": true,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "wvcfagdd",
						"name": "usePassword",
						"type": "bool",
						"required": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "9huin8v9",
						"name": "#password",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "noyiex3k",
						"name": "#privateKey",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "vl4cotoc",
						"name": "#privateKeyPassphrase",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "w8rbwtjz",
						"name": "hostname",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "ej3fdqgn",
						"name": "hostKey",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"listRule": "@collection.userServers.userId=@request.auth.id \u0026\u0026 @collection.userServers.serverId=id",
				"viewRule": "@collection.userServers.userId=@request.auth.id \u0026\u0026 @collection.userServers.serverId=id",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "xrd5dhd1yoe4x06",
				"created": "2022-09-24 01:52:59.404Z",
				"updated": "2022-12-20 21:21:03.516Z",
				"name": "serverLogs",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "xj7il86q",
						"name": "serverId",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"collectionId": "0wy8xlddg0uekql",
							"cascadeDelete": true
						}
					},
					{
						"system": false,
						"id": "37xvbs9x",
						"name": "type",
						"type": "select",
						"required": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"info",
								"warn",
								"error",
								"hostKey",
								"hostName"
							]
						}
					},
					{
						"system": false,
						"id": "1ybucrs3",
						"name": "message",
						"type": "text",
						"required": true,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "gjl7yonn",
						"name": "payload",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"listRule": "@collection.userServers.userId=@request.auth.id \u0026\u0026 @collection.userServers.serverId=serverId",
				"viewRule": "@collection.userServers.userId=@request.auth.id \u0026\u0026 @collection.userServers.serverId=serverId",
				"createRule": null,
				"updateRule": null,
				"deleteRule": "@collection.userServers.userId=@request.auth.id \u0026\u0026 @collection.userServers.serverId=serverId",
				"options": {}
			},
			{
				"id": "12iz4a90ov5tprg",
				"created": "2022-09-24 01:52:59.405Z",
				"updated": "2022-12-20 21:21:03.522Z",
				"name": "userServers",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "nomwosjc",
						"name": "userId",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"collectionId": "systemprofiles0",
							"cascadeDelete": true
						}
					},
					{
						"system": false,
						"id": "mqupzxls",
						"name": "serverId",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"collectionId": "0wy8xlddg0uekql",
							"cascadeDelete": true
						}
					},
					{
						"system": false,
						"id": "qhlkt19a",
						"name": "options",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"listRule": "@request.auth.id=userId",
				"viewRule": "@request.auth.id=userId",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "systemprofiles0",
				"created": "2022-12-20 21:21:03.518Z",
				"updated": "2022-12-20 21:21:03.518Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "pbfieldname",
						"name": "name",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "pbfieldavatar",
						"name": "avatar",
						"type": "file",
						"required": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"maxSize": 5242880,
							"mimeTypes": [
								"image/jpg",
								"image/jpeg",
								"image/png",
								"image/svg+xml",
								"image/gif"
							],
							"thumbs": null
						}
					}
				],
				"listRule": "id = @request.auth.id",
				"viewRule": "id = @request.auth.id",
				"createRule": "",
				"updateRule": "id = @request.auth.id",
				"deleteRule": null,
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": false,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"requireEmail": true
				}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		// no revert since the configuration on the environment, on which
		// the migration was executed, could have changed via the UI/API
		return nil
	})
}
