package main

import (
	"fmt"
	"time"
)

type Message struct {
	UserID string    `json:"userId"`
	Text   string    `json:"text"`
	Time   time.Time `json:"time"`
}

func (message *Message) String() string {
	return fmt.Sprintf("%s - %s: %s",
		message.UserID,
		message.Time.Format(time.Kitchen),
		message.Text)
}
