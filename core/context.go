package core

import (
	"github.com/chuccp/cokePush/config"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
)

type registerHandle func(value ...interface{}) interface{}

type Context struct {
	handleFuncMap map[string]registerHandle
	dock          *dock
	config             *config.Config
	queryForwardHandle registerHandle
}

func (context *Context) RegisterHandle(handleName string, handle registerHandle) {
	context.handleFuncMap[handleName] = handle
}
func (context *Context) GetHandle(handleName string) registerHandle {
	return context.handleFuncMap[handleName]
}

func (context *Context) GetConfig() *config.Config {
	return context.config
}
func (context *Context) AddUser(iUser user.IUser) {
	context.dock.AddUser(iUser)
}
func (context *Context) GetUser(username string,f func(user.IUser)bool) {
	context.dock.UserStore.GetUser(username,f)
}

func (context *Context) Query(queryName string, value ...interface{}) interface{} {
	handle := context.GetHandle(queryName)
	v := handle(value...)
	cluHandle:=context.queryForwardHandle
	if cluHandle!=nil{
		iv:=make([]interface{},0)
		iv = append(iv, queryName)
		iv = append(iv, v)
		for _,vi:=range value{
			iv = append(iv, vi)
		}
		vs:=cluHandle(iv...)
		return vs
	}else{
		vvs:=make([]interface{},0)
		vvs = append(vvs, v)
		return vvs
	}
}

func (context *Context) UserNum() int32{
	return context.dock.UserNum()
}
func (context *Context) DeleteUser(iUser user.IUser) {
	context.dock.DeleteUser(iUser)
}
func (context *Context) SendMessage(msg message.IMessage, write WriteFunc) {
	context.dock.sendMessage(msg, write)
}
func (context *Context) SendMessageNoForward(msg message.IMessage, write WriteFunc) {
	context.dock.SendMessageNoForward(msg, write)
}

func (context *Context) HandleAddUser(handleAddUser HandleAddUser) {
	context.dock.handleAddUser = handleAddUser
}

func (context *Context) SendNum()int{
	return context.dock.sendNum()
}
func (context *Context) ReplyNum()int{
	return context.dock.replyNum()
}

func (context *Context) HandleDeleteUser(handleDeleteUser HandleDeleteUser) {
	context.dock.handleDeleteUser = handleDeleteUser
}
func (context *Context) HandleSendMessage(handleSendMessage HandleSendMessage) {
	context.dock.handleSendMessage = handleSendMessage
}
func (context *Context) QueryForwardHandle(queryHandle registerHandle) {
	context.queryForwardHandle = queryHandle
}

func newContext(config *config.Config) *Context {
	return &Context{handleFuncMap: make(map[string]registerHandle), dock: newDock(), config: config}
}
func (context *Context) Init() {
	go context.dock.exchangeSendMsg()
	go context.dock.exchangeReplyMsg()
}
