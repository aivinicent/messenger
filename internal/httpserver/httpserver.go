package httpserver

import (
	"encoding/json"
	"messenger/internal/dbclient"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func Start() {
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

	lastNewMessageID = lastMessageID

	http.HandleFunc("/messages", messagesHandler)
	http.HandleFunc("/live-messages", liveMessagedHandler)
	http.ListenAndServe(":8080", nil)
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)

		var body messagesPost
		err := decoder.Decode(&body)
		if err != nil {
			panic("httpserver.messagesHandler: decoder.Decode: " + err.Error())
		}

		newMessageID, err := dbclient.AddMessage(body.Body)
		if err != nil {
			panic("httpserver.messagesHandler: dbclient.AddMessage: " + err.Error())
		}

		if lastNewMessageID < newMessageID {
			lastNewMessageID = newMessageID
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	} else if r.Method == "GET" {
		messages, err := dbclient.GetMessages()
		if err != nil {
			panic("httpserver.messagesHandler: dbclient.GetMessages: " + err.Error())
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		err = json.NewEncoder(w).Encode(messages)
		if err != nil {
			panic("httpserver.messagesHandler: json.NewEncoder.Encode: " + err.Error())
		}
	}
}

func liveMessagedHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic("httpserver.liveMessagedHandler: upgrader.Upgrade: " + err.Error())
	}
	defer c.Close()

	lastSendedMessageID := lastNewMessageID

	for {
		if lastNewMessageID != lastSendedMessageID {
			lastSendedMessageID++

			message, err := dbclient.GetMessage(lastSendedMessageID)
			if err != nil {
				panic("httpserver.liveMessagedHandler: dbclient.GetMessage: " + err.Error())
			}

			err = c.WriteMessage(1, []byte(message.Body))
			if err != nil {
				panic("httpserver.liveMessagedHandler: c.WriteMessage: " + err.Error())
			}
		}

		time.Sleep(1000)
	}
}

var lastNewMessageID int64 = 0

type messagesPost struct {
	Body string
}
