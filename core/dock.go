package core

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
	"github.com/chuccp/utils/log"
	"github.com/chuccp/utils/queue"
)

type HandleAddUser func(iUser user.IUser)
type HandleDeleteUser func(username string)
type HandleSendMessage func(iMessage *DockMessage, writeFunc user.WriteFunc)

type HandleSendMultiMessage func(fromUser string, usernames *[]string, text string, f func(username string, status int))

type dock struct {
	sendMsg                *queue.Queue
	replyMsg               *queue.Queue
	UserStore              *user.Store
	handleAddUser          HandleAddUser
	handleDeleteUser       HandleDeleteUser
	handleSendMessage      HandleSendMessage
	handleSendMultiMessage HandleSendMultiMessage
	//sendIndexNum      uint32
	//replyIndexNum     uint32
}

func newDock() *dock {
	return &dock{sendMsg: queue.NewQueue(), UserStore: user.NewStore(), replyMsg: queue.NewQueue()}
}
func (dock *dock) sendMessage(iMessage message.IMessage, write user.WriteFunc) {
	msg := newDockMessage(iMessage, write)
	msg.IsForward = true
	dock.sendMsg.Offer(msg)
}
func (dock *dock) SendMessageNoForward(iMessage message.IMessage, write user.WriteFunc) {
	msg := newDockMessage(iMessage, write)
	msg.IsForward = false
	dock.sendMsg.Offer(msg)
}
func (dock *dock) SendMessageNoReplay(iMessage message.IMessage) {
	msg := newDockMessageNoReplay(iMessage)
	msg.IsForward = false
	dock.sendMsg.Offer(msg)
}

func (dock *dock) eachUsers(f func(key string, value *user.StoreUser) bool) {
	dock.UserStore.EachUsers(f)
}

func (dock *dock) writeUserMsg(msg *DockMessage) {
	var flag bool
	var ee error
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
		} else {
			msg.flag = false
			msg.err = NoFoundUser
			dock.replyMessage(msg)
		}
	} else {
		msg.flag = flag
		msg.err = ee
		dock.replyMessage(msg)
	}
}
func (dock *dock) AddUser(iUser user.IUser) {
	fa := dock.UserStore.AddUser(iUser)
	if fa {
		if dock.handleAddUser != nil {
			dock.handleAddUser(iUser)
		}
	}
}
func (dock *dock) DeleteUser(iUser user.IUser) {
	fa := dock.UserStore.DeleteUser(iUser)
	if fa {
		if dock.handleDeleteUser != nil {
			dock.handleDeleteUser(iUser.GetUsername())
		}
	}
}
func (dock *dock) UserNum() int32 {
	return dock.UserStore.GetUserNum()
}

func (dock *dock) replyMessage(msg *DockMessage) {
	if msg.replay {
		log.DebugF("????????????????????????:{}", msg.InputMessage.GetMessageId())
		dock.replyMsg.Offer(msg)
	} else {
		free(msg)
	}
}
func (dock *dock) exchangeReplyMsg() {
	log.DebugF("????????????????????????")
	for {
		msg, _ := dock.replyMsg.Poll()
		dockMessage := msg.(*DockMessage)
		if msg != nil {
			dockMessage.write(dockMessage.err, dockMessage.flag)
			free(dockMessage)
		}
	}
}
func (dock *dock) sendNum() int32 {

	return dock.sendMsg.Num()
}
func (dock *dock) replyNum() int32 {

	return dock.replyMsg.Num()
}
func (dock *dock) exchangeSendMsg() {
	log.DebugF("????????????????????????")

	for {
		msg, _ := dock.sendMsg.Poll()
		if msg != nil {
			dm := msg.(*DockMessage)
			dock.writeUserMsg(dm)
		}
	}
}
