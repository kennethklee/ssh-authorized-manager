package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ox1bxwliw4rxoy6")
		if err != nil {
			return err
		}

		// update
		edit_server := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "olfi3v8h",
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
		}`), edit_server)
		collection.Schema.AddField(edit_server)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ox1bxwliw4rxoy6")
		if err != nil {
			return err
		}

		// update
		edit_server := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
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
		}`), edit_server)
		collection.Schema.AddField(edit_server)

		return dao.SaveCollection(collection)
	})
}
