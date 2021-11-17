package user

import "github.com/chuccp/cokePush/message"

type IUser interface {
	SetUsername(username string)
	GetUsername() string
	WriteMessage(iMessage message.IMessage) error
}
type User struct {
	username string
	write    message.Write
}

func (u *User) SetUsername(username string) {
	u.username = username
}
func (u *User) GetUsername() string {
	return u.username
}
func (u *User) WriteMessage(iMessage message.IMessage) error {
	return u.write.WriteMessage(iMessage)
}
func AddUser(iUser IUser) {

}
func Range(f func(iUser IUser) bool) {

}

func GetUser(username string)IUser  {
	return nil
}
