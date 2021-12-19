package httpserver

import (
	"encoding/json"
	"messenger/internal/dbclient"
	"net/http"
)

func Start() {
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/messages", messagesHandler)
	http.ListenAndServe(":8000", nil)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("ok")
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
	}
}

type messagesPost struct {
	Body string
}
