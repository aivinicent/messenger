package dbclient

import (
	"messenger/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func New() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("dbclient.New: " + err.Error())
	}
}

func AddMessage(body string) (err error) {
	timestamp := time.Now().Unix()

	tx := db.Exec("INSERT INTO messenger.messages VALUES (DEFAULT, ?, ?)", body, timestamp)
	if tx.Error != nil {
		err = tx.Error
	}

	return
}

func GetMessages() (messages []models.Message, err error) {
	var messagesInternal []messageInternal

	tx := db.Raw("SELECT * FROM messenger.messages").Scan(&messagesInternal)
	if tx.Error != nil {
		err = tx.Error
	}

	for _, messageInternal := range messagesInternal {
		message := models.Message{
			Id:        messageInternal.Id,
			Body:      messageInternal.Body,
			Timestamp: time.Unix(messageInternal.Timestamp, 0).Format(time.RFC822),
		}

		messages = append(messages, message)
	}

	return
}

type messageInternal struct {
	Id        int64
	Body      string
	Timestamp int64
}
