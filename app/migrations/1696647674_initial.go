package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `[
			{
				"id": "_pb_users_auth_",
				"created": "2023-10-06 17:55:56.577Z",
				"updated": "2023-10-06 17:55:56.578Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "users_name",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "users_avatar",
						"name": "avatar",
						"type": "file",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"maxSize": 5242880,
							"mimeTypes": [
								"image/jpeg",
								"image/png",
								"image/svg+xml",
								"image/gif",
								"image/webp"
							],
							"thumbs": null,
							"protected": false
						}
					}
				],
				"indexes": [],
				"listRule": "id = @request.auth.id",
				"viewRule": "id = @request.auth.id",
				"createRule": "",
				"updateRule": "id = @request.auth.id",
				"deleteRule": "id = @request.auth.id",
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": true,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"requireEmail": false
				}
			},
			{
				"id": "dnu8pjil2uttnkn",
				"created": "2023-10-07 01:23:15.411Z",
				"updated": "2023-10-07 02:29:50.490Z",
				"name": "servers",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "6smizrql",
						"name": "state",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "0o8qjw3o",
						"name": "name",
						"type": "text",
						"required": true,
						"presentable": true,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "vmcpds7w",
						"name": "host",
						"type": "text",
						"required": true,
						"presentable": true,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "pdhm3skj",
						"name": "port",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "wm9ucxdt",
						"name": "username",
						"type": "text",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "bfk8oh2l",
						"name": "usePassword",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "otfh4hum",
						"name": "password",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "lueuer9r",
						"name": "privateKey",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "hjwuys2r",
						"name": "privateKeyPassPhrase",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "1qkoco1j",
						"name": "hostName",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "oxt8xljt",
						"name": "hostKey",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": "@collection.userServers.user=@request.auth.id && @collection.userServers.server=id",
				"viewRule": "@collection.userServers.user=@request.auth.id && @collection.userServers.server=id",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "ox1bxwliw4rxoy6",
				"created": "2023-10-07 01:24:13.233Z",
				"updated": "2023-10-07 01:29:04.286Z",
				"name": "userServers",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "s0nizb1i",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": true,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "olfi3v8h",
						"name": "server",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "dnu8pjil2uttnkn",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "mffy3wtu",
						"name": "options",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "mrc76h3nygjqnbw",
				"created": "2023-10-07 01:28:46.605Z",
				"updated": "2023-10-07 02:34:17.340Z",
				"name": "publicKeys",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "tso0xqny",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": true,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "jun5yuce",
						"name": "type",
						"type": "text",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "daqyijv2",
						"name": "key",
						"type": "text",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "iyrbl8gx",
						"name": "comment",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id=user",
				"viewRule": "@request.auth.id=user",
				"createRule": "@request.auth.id=user",
				"updateRule": "@request.auth.id=user",
				"deleteRule": "@request.auth.id=user",
				"options": {}
			},
			{
				"id": "nv0a97w7tcczhyz",
				"created": "2023-10-07 01:30:40.211Z",
				"updated": "2023-10-07 01:30:40.211Z",
				"name": "serverLogs",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "abrzk3kb",
						"name": "server",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "dnu8pjil2uttnkn",
							"cascadeDelete": true,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "vyhmkfu2",
						"name": "type",
						"type": "select",
						"required": true,
						"presentable": false,
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
						"id": "d9u7gnxz",
						"name": "message",
						"type": "text",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "3a9igfzb",
						"name": "payload",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}
