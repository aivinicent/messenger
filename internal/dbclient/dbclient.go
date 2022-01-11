package dbclient

import (
	"fmt"
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

func AddMessage(body string) (id int64, err error) {
	message := messages{
		Body:      body,
		Timestamp: time.Now().Unix(),
	}

	err = db.Table("messenger.messages").Create(&message).Error
	if err != nil {
		return
	}

	id = message.ID

	return
}

func GetMessages() (allMessages []models.Message, err error) {
	var messagesInternal []messages

	tx := db.Raw("SELECT * FROM messenger.messages").Scan(&messagesInternal)
	if tx.Error != nil {
		err = tx.Error
	}

	for _, internalMessage := range messagesInternal {
		allMessages = append(allMessages, convertInternalMessage(internalMessage))
	}

	return
}

func GetMessage(id int64) (message models.Message, err error) {
	var messagesInternal []messages

	tx := db.Raw("SELECT * FROM messenger.messages").Where("id=?", id).Scan(&messagesInternal)
	if tx.Error != nil {
		err = tx.Error
	}

	if len(messagesInternal) == 0 {
		err = fmt.Errorf("Not found")
		return
	}

	message = convertInternalMessage(messagesInternal[0])

	return
}

func convertInternalMessage(internalMessage messages) (message models.Message) {
	message = models.Message{
		Id:        internalMessage.ID,
		Body:      internalMessage.Body,
		Timestamp: time.Unix(internalMessage.Timestamp, 0).Format(time.RFC822),
	}

	return
}

type messages struct {
	ID        int64
	Body      string
	Timestamp int64
}
