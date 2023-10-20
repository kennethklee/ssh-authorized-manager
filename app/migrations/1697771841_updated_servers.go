package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("dnu8pjil2uttnkn")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("@collection.userServers.user?=@request.auth.id && @collection.userServers.server?=id")

		collection.ViewRule = types.Pointer("@collection.userServers.user?=@request.auth.id && @collection.userServers.server?=id")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("dnu8pjil2uttnkn")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("@collection.userServers.user=@request.auth.id && @collection.userServers.server=id")

		collection.ViewRule = types.Pointer("@collection.userServers.user=@request.auth.id && @collection.userServers.server=id")

		return dao.SaveCollection(collection)
	})
}
