package user

import (
	"github.com/chuccp/cokePush/message"
)

type IUser interface {
	SetUsername(username string)
	GetUsername() string
	WriteMessage(iMessage message.IMessage) error
}
type User struct {
	Username string
	Write    message.Write
}

func (u *User) SetUsername(username string) {
	u.Username = username
}
func (u *User) GetUsername() string {
	return u.Username
}
func (u *User) WriteMessage(iMessage message.IMessage) error {
	return u.Write.WriteMessage(iMessage)
}

