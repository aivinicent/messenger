package app

import (
	"messenger/internal/dbclient"
	"messenger/internal/httpserver"
)

func Run() {
	dbclient.New()

	allMessages, err := dbclient.GetMessages()
	if err != nil {
		panic("app.Run: dbclient.GetMessages: " + err.Error())
	}

	lastMessageID := int64(0)
	for _, message := range allMessages {
		if message.Id > lastMessageID {
			lastMessageID = message.Id
		}
	}

	httpserver.LastNewMessageID = lastMessageID

	httpserver.Start()
}
