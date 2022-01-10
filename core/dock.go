package core

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
)

type WriteFunc func(err error, hasUser bool)

type HandleAddUser func(iUser user.IUser)
type HandleDeleteUser func(username string)
type HandleSendMessage func(iMessage *DockMessage, writeFunc WriteFunc)

type dock struct {
	sendMsg           *Queue
	replyMsg          *Queue
	UserStore         *user.Store
	handleAddUser     HandleAddUser
	handleDeleteUser  HandleDeleteUser
	handleSendMessage HandleSendMessage
}

func newDock() *dock {
	return &dock{sendMsg: NewQueue(), UserStore: user.NewStore(), replyMsg: NewQueue()}
}

func (dock *dock) sendMessage(iMessage message.IMessage, write WriteFunc) {
	msg := newDockMessage(iMessage, write)
	msg.IsForward = true
	dock.sendMsg.offer(newDockMessage(iMessage, write))
}
func (dock *dock) SendMessageNoForward(iMessage message.IMessage, write WriteFunc) {
	msg := newDockMessage(iMessage, write)
	msg.IsForward = false
	dock.sendMsg.offer(msg)
}

func (dock *dock) writeUserMsg(msg *DockMessage) (flag bool, ee error) {
	userId := msg.InputMessage.GetString(message.ToUser)
	log.DebugF("信息发送:{}", userId)
	flag = dock.UserStore.GetUser(msg.InputMessage.GetString(message.ToUser), func(iUser user.IUser) bool {
		err := iUser.WriteMessage(msg.InputMessage)
		ee = err
		return err != nil
	})
	if msg.IsForward && !flag {
		if dock.handleSendMessage != nil {
			dock.handleSendMessage(msg, func(err error, hasUser bool) {
				msg.flag = hasUser
				msg.err = err
				dock.replyMessage(msg)
			})
		}
	}
	return
}
func (dock *dock) AddUser(iUser user.IUser) {
	dock.UserStore.AddUser(iUser)
}
func (dock *dock) DeleteUser(iUser user.IUser) {
	dock.UserStore.DeleteUser(iUser)
}

func (dock *dock) replyMessage(msg *DockMessage) {
	dock.replyMsg.offer(msg)
}
func (dock *dock) exchangeReplyMsg() {
	log.Info("启动信息反馈处理")
	var i = 0
	for {
		msg := dock.replyMsg.poll()
		if msg != nil {
			msg.write(msg.err, msg.flag)
		}
		i++
		if i>>10 == 1 {
			i = 0
			log.InfoF("当前反馈池剩下 :{} 未处理", dock.replyMsg.Num())
		}
	}
}

func (dock *dock) exchangeSendMsg() {
	log.Info("启动信息发送处理")
	var i = 0
	for {
		msg := dock.sendMsg.poll()
		if msg != nil {
			fa, err := dock.writeUserMsg(msg)
			log.DebugF("fa:{} err:{}", fa, err)
			msg.flag = fa
			msg.err = err
			dock.replyMessage(msg)
		}
		i++
		if i>>10 == 1 {
			i = 0
			log.InfoF("当前信息池剩下 :{} 未处理", dock.sendMsg.Num())
		}
	}
}
