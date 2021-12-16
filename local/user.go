package local

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
)

type User struct {
	*user.User
}

func NewUser(username string, f func(iMessage message.IMessage) error) *User {
	return &User{User: &user.User{Username: username, Write: newWrite(f)}}
}

type Write struct {
	f func(iMessage message.IMessage) error
}

func newWrite(f func(iMessage message.IMessage) error) *Write {
	return &Write{f: f}
}
func (write *Write) WriteMessage(iMessage message.IMessage) error {
	return write.f(iMessage)
}
