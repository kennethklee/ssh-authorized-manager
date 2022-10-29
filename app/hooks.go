package main

import (
	"os"
	"sshauthman/worker"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type HooksConfig struct {
	AutoVerifyUser bool
}

func HooksConfigFromEnv() HooksConfig {
	return HooksConfig{
		AutoVerifyUser: os.Getenv("AUTO_VERIFY") != "",
	}
}

func RegisterHooks(app core.App, config HooksConfig) {
	// auto verify feature flag
	if config.AutoVerifyUser {
		app.OnUserBeforeCreateRequest().Add(func(e *core.UserCreateEvent) error {
			e.User.Verified = true
			return nil
		})
	}

	app.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		return syncServerHookHandler("create", app.Dao(), e.Record)
	})
	app.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		return syncServerHookHandler("update", app.Dao(), e.Record)
	})
	app.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		return syncServerHookHandler("delete", app.Dao(), e.Record)
	})
}

// hooks for server record events
// hooks for userServer record events
// hooks for publickey record events
// conditionally sync servers when a record is created or updated or deleted
func syncServerHookHandler(action string, dao *daos.Dao, r *models.Record) error {
	if r.Collection().Name == "servers" {
		worker.SubmitAndForget(&worker.SyncServerWork{Server: r})
	}
	if r.Collection().Name == "publicKeys" {
		// find all servers that need this public key updated
		servers, err := collectServersByUserID(dao, r.GetStringDataValue("userId"))
		if err != nil {
			return err
		}

		for _, server := range servers {
			worker.SubmitAndForget(&worker.SyncServerWork{Server: server})
		}
	}
	if r.Collection().Name == "userServers" {
		serverCollection, _ := dao.FindCollectionByNameOrId("servers")

		// fetch server
		server, _ := dao.FindRecordById(serverCollection, r.GetStringDataValue("serverId"), nil)
		worker.SubmitAndForget(&worker.SyncServerWork{Server: server})
	}

	return nil
}

func collectServersByUserID(dao *daos.Dao, userId string) ([]*models.Record, error) {
	userServerCollection, err := dao.FindCollectionByNameOrId("userServers")
	if err != nil {
		return nil, err
	}
	serverCollection, err := dao.FindCollectionByNameOrId("servers")
	if err != nil {
		return nil, err
	}

	userServers, err := dao.FindRecordsByExpr(userServerCollection, &dbx.HashExp{"userId": userId})
	if err != nil {
		return nil, err
	}

	var serverIds = []interface{}{}
	for _, userServer := range userServers {
		serverIds = append(serverIds, userServer.GetStringDataValue("serverId"))
	}

	servers, err := dao.FindRecordsByExpr(serverCollection, dbx.In("id", serverIds...))
	if err != nil {
		return nil, err
	}
	return servers, nil
}
