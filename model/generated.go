// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Event interface {
	IsEvent()
}

// Событие с новым сообщением.
type MessageEvent struct {
	Chat    *Chat    `json:"chat"`
	Message *Message `json:"message"`
}

func (MessageEvent) IsEvent() {}

type ChatType string

const (
	// Канал: может писать владелец, подписчики – читают.
	ChatTypeChannel ChatType = "CHANNEL"
	// Групповой чат с произвольным количеством людей.
	ChatTypeGroup ChatType = "GROUP"
	// Личный чат между двумя пользователями.
	ChatTypePrivate ChatType = "PRIVATE"
)

var AllChatType = []ChatType{
	ChatTypeChannel,
	ChatTypeGroup,
	ChatTypePrivate,
}

func (e ChatType) IsValid() bool {
	switch e {
	case ChatTypeChannel, ChatTypeGroup, ChatTypePrivate:
		return true
	}
	return false
}

func (e ChatType) String() string {
	return string(e)
}

func (e *ChatType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ChatType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ChatType", str)
	}
	return nil
}

func (e ChatType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
