package core

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
)

type Context struct {
	UserStore *user.Store
}

func (context *Context) AddUser(iUser user.IUser) {
	context.UserStore.AddUser(iUser)
}
func (context *Context) GetUser(username string) user.IUser {
	return context.UserStore.GetUser(username)
}
func newContext() *Context {
	return &Context{UserStore: user.NewStore()}
}

func (context *Context) SendMessage(iMessage message.IMessage) error {
	iUser:=context.GetUser(iMessage.GetString(message.ToUser))
	if iUser==nil{
		return NoFoundUser
	}
	return iUser.WriteMessage(iMessage)
}
