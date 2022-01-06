package core

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
)

type registerHandle func(value ...interface{}) interface{}

type Context struct {
	handleFuncMap map[string]registerHandle
	dock *dock
}
func (context *Context) RegisterHandle(handleName string, handle registerHandle) {
	context.handleFuncMap[handleName] = handle
}
func (context *Context) GetHandle(handleName string)registerHandle {
	return context.handleFuncMap[handleName]
}
func (context *Context) DeleteUser(iUser user.IUser) {
	context.dock.UserStore.DeleteUser(iUser)

}
func (context *Context) SendMessage(msg message.IMessage)error {
	return context.dock.sendMessage(msg)
}
func (context *Context) Handle(msg message.IMessage,writeRead user.IUser)error{
	 context.dock.handleMessage(msg,writeRead)
	 return nil
}
func newContext() *Context {
	return &Context{handleFuncMap: make(map[string]registerHandle),dock:newDock()}
}
func (context *Context) Init() {
	go context.dock.exchangeSendMsg()
}
