package models

type Message struct {
	Id        int64  `json:"id"`
	Body      string `json:"body"`
	Timestamp string `json:"timestamp"`
}
