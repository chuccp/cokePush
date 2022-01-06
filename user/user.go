package user

import (
	"github.com/chuccp/cokePush/message"
)

type IUser interface {
	WriteMessage(iMessage message.IMessage) error
	ReadMessage() (message.IMessage,error)
	GetId() string
	GetUsername() string
	SetUsername(username string)
}