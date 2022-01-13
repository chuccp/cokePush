package user

import (
	"github.com/chuccp/cokePush/message"
)

type IUser interface {
	WriteMessage(iMessage message.IMessage) error
	GetId() string
	GetUsername() string
	GetRemoteAddress() string
	SetUsername(username string)
}