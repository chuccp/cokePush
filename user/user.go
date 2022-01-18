package user

import (
	"github.com/chuccp/cokePush/message"
)
type WriteFunc func(err error, hasUser bool)
type IUser interface {
	WriteMessage(iMessage message.IMessage) error
	GetId() string
	GetUsername() string
	GetRemoteAddress() string
	SetUsername(username string)
	WriteMessageFunc(iMessage message.IMessage,writeFunc WriteFunc)
}