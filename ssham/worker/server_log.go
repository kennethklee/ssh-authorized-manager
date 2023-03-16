package worker

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type ServerLog struct {
	serverId string
	msgType  string
	message  string
	payload  string
}

func NewServerLog(serverId string, msgType string, message string, payload string) ServerLog {
	return ServerLog{
		serverId: serverId,
		msgType:  msgType,
		message:  message,
		payload:  payload,
	}
}

func SaveServerLog(message ServerLog) error {
	serverLogsCollection, _ := app.Dao().FindCollectionByNameOrId("serverLogs")

	serverLog := models.NewRecord(serverLogsCollection)

	serverLog.Set("serverId", message.serverId)
	serverLog.Set("type", message.msgType)
	serverLog.Set("message", message.message)
	serverLog.Set("payload", message.payload)

	return app.Dao().SaveRecord(serverLog)
}

func CreateServerLog(serverRecord *models.Record, msgType string, message string, payload string) error {
	serverLogsCollection, _ := app.Dao().FindCollectionByNameOrId("serverLogs")

	serverLog := models.NewRecord(serverLogsCollection)
	serverLog.Set("serverId", serverRecord.Id)
	serverLog.Set("type", msgType)
	serverLog.Set("message", message)
	serverLog.Set("payload", payload)

	if err := app.Dao().SaveRecord(serverLog); err != nil {
		return err
	}

	serverRecord.Set("lastState", msgType)
	app.Dao().SaveRecord(serverRecord)

	event := &core.RecordCreateEvent{Record: serverLog}
	return app.OnRecordAfterCreateRequest().Trigger(event)
}
