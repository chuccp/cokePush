package core

import (
	"github.com/chuccp/cokePush/config"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
)

type registerHandle func(value ...interface{}) interface{}

type Context struct {
	handleFuncMap      map[string]registerHandle
	dock               *dock
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
func (context *Context) GetUser(username string, f func(user.IUser) bool) {
	context.dock.UserStore.GetUser(username, f)
}

func (context *Context) Query(queryName string, value ...interface{}) interface{} {
	handle := context.GetHandle(queryName)
	v := handle(value...)
	cluHandle := context.queryForwardHandle
	if cluHandle != nil {
		iv := make([]interface{}, 0)
		iv = append(iv, queryName)
		iv = append(iv, v)
		for _, vi := range value {
			iv = append(iv, vi)
		}
		vs := cluHandle(iv...)
		return vs
	} else {
		vvs := make([]interface{}, 0)
		vvs = append(vvs, v)
		return vvs
	}
}

func (context *Context) UserNum() int32 {
	return context.dock.UserNum()
}
func (context *Context) DeleteUser(iUser user.IUser) {
	context.dock.DeleteUser(iUser)
}
func (context *Context) SendMessage(msg message.IMessage, write user.WriteFunc) {
	context.dock.sendMessage(msg, write)
}

func (context *Context) SendMultiMessage(fromUser string, usernames []string, text string, f func(username string, status int)) {

	localUser := make([]string, 0)
	remoteLocalUser := make([]string, 0)
	for _, v := range usernames {
		if context.dock.UserStore.Has(v) {
			localUser = append(localUser, v)
			f(v, 1)
		} else {
			remoteLocalUser = append(remoteLocalUser, v)
		}
	}
	if len(localUser) > 0 {
		go context.SendMultiMessageNoReplay(fromUser, &localUser, text)
	}
	if len(remoteLocalUser) > 0 {
		if context.dock.handleSendMultiMessage != nil {
			context.dock.handleSendMultiMessage(fromUser, &remoteLocalUser, text, f)
		}
	}
}
func (context *Context) SendMultiMessageNoReplay(fromUser string, usernames *[]string, text string) {
	for _, v := range *usernames {
		msg := message.CreateBasicMessage(fromUser, v, text)
		context.SendMessageNoReplay(msg)
	}
}

func (context *Context) SendMessageNoForward(msg message.IMessage, write user.WriteFunc) {
	context.dock.SendMessageNoForward(msg, write)
}

func (context *Context) SendMessageNoReplay(msg message.IMessage) {
	context.dock.SendMessageNoReplay(msg)
}
func (context *Context) HandleAddUser(handleAddUser HandleAddUser) {
	context.dock.handleAddUser = handleAddUser
}

func (context *Context) SendNum() int32 {
	return context.dock.sendNum()
}
func (context *Context) ReplyNum() int32 {
	return context.dock.replyNum()
}

func (context *Context) HandleDeleteUser(handleDeleteUser HandleDeleteUser) {
	context.dock.handleDeleteUser = handleDeleteUser
}
func (context *Context) HandleSendMessage(handleSendMessage HandleSendMessage) {
	context.dock.handleSendMessage = handleSendMessage
}
func (context *Context) HandleSendMultiMessage(handleSendMultiMessage HandleSendMultiMessage) {
	context.dock.handleSendMultiMessage = handleSendMultiMessage
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

func (context *Context) EachUsers(f func(key string, value *user.StoreUser) bool) {
	context.dock.eachUsers(f)
}
