package user

import (
	"github.com/chuccp/cokePush/message"
)

type IUser interface {
	WriteMessage(iMessage message.IMessage) error
	ReadMessage() (message.IMessage,error)
	Close()
	GetId() string
	GetUsername() string
}