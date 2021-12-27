package core

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
)

type Context struct {
	UserStore *user.Store
	sendMsg   *queue
	writeMsg  *queue
}

func (context *Context) exchangeSendMsg() {
	log.Info("启动信息发送处理")
	var i = 0
	for {
		msg := context.sendMsg.poll()
		if msg != nil {
			userId:=msg.GetString(message.ToUser)
			context.GetUser(msg.GetString(message.ToUser))
			log.InfoF("信息发送:{}",userId)
		}
		i++
		if i>>10==1{
			i = 0
			log.InfoF("当前信息池剩下 :{} 未处理",context.sendMsg.Num())
		}
	}
}

func (context *Context) AddUser(iUser user.IUser) {
	context.UserStore.AddUser(iUser)
}
func (context *Context) GetUser(username string) user.IUser {
	return context.UserStore.GetUser(username)
}
func newContext() *Context {
	return &Context{UserStore: user.NewStore(), sendMsg: newQueue(), writeMsg: newQueue()}
}

func (context *Context) SendMessage(iMessage message.IMessage) error {
	context.sendMsg.offer(iMessage)

	return nil
}

func (context *Context) Init() {
	go context.exchangeSendMsg()
}
