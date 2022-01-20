package user

import (
	"github.com/chuccp/cokePush/message"
	"time"
)
type WriteFunc func(err error, hasUser bool)
type IUser interface {
	WriteMessage(iMessage message.IMessage) error
	GetId() string
	GetUsername() string
	GetRemoteAddress() string
	SetUsername(username string)
	LastLiveTime()*time.Time
	CreateTime()*time.Time
	WriteMessageFunc(iMessage message.IMessage,writeFunc WriteFunc)
}