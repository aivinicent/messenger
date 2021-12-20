package httpserver

import (
	"encoding/json"
	"messenger/internal/dbclient"
	"net/http"
)

func Start() {
	http.HandleFunc("/messages", messagesHandler)
	http.ListenAndServe(":8000", nil)
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)

		var body messagesPost
		err := decoder.Decode(&body)
		if err != nil {
			panic("httpserver.messagesHandler: decoder.Decode: " + err.Error())
		}

		err = dbclient.AddMessage(body.Body)
		if err != nil {
			panic("httpserver.messagesHandler: dbclient.AddMessage: " + err.Error())
		}
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

type messagesPost struct {
	Body string
}
