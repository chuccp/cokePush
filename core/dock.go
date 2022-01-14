package core

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
	"github.com/chuccp/cokePush/util"
)

type WriteFunc func(err error, hasUser bool)

type HandleAddUser func(iUser user.IUser)
type HandleDeleteUser func(username string)
type HandleSendMessage func(iMessage *DockMessage, writeFunc WriteFunc)

type dock struct {
	sendMsg           *util.Queue
	replyMsg          *util.Queue
	UserStore         *user.Store
	handleAddUser     HandleAddUser
	handleDeleteUser  HandleDeleteUser
	handleSendMessage HandleSendMessage
	//sendIndexNum      uint32
	//replyIndexNum     uint32
}

func newDock() *dock {
	return &dock{sendMsg:util.NewQueue(), UserStore: user.NewStore(), replyMsg: util.NewQueue()}
}
func (dock *dock) sendMessage(iMessage message.IMessage, write WriteFunc) {
	msg := newDockMessage(iMessage, write)
	msg.IsForward = true
	dock.sendMsg.Offer(msg)
}
func (dock *dock) SendMessageNoForward(iMessage message.IMessage, write WriteFunc) {
	msg := newDockMessage(iMessage, write)
	msg.IsForward = false
	dock.sendMsg.Offer(msg)
}

func (dock *dock) writeUserMsg(msg *DockMessage) (flag bool, ee error) {
	userId := msg.InputMessage.GetString(message.ToUser)
	log.DebugF("信息发送:{}", userId)
	flag = dock.UserStore.GetUser(msg.InputMessage.GetString(message.ToUser), func(iUser user.IUser) bool {
		err := iUser.WriteMessage(msg.InputMessage)
		ee = err
		return err != nil
	})
	log.DebugF("信息发送本机  IsForward:{}  flag:{} msgId:{}", msg.IsForward,flag,msg.InputMessage.GetMessageId())
	if msg.IsForward && !flag {
		if dock.handleSendMessage != nil {
			 dock.handleSendMessage(msg, func(err error, hasUser bool) {
				log.DebugF("信息发送本机 有反馈 hasUser:{} msgId:{}", hasUser,msg.InputMessage.GetMessageId())
				msg.flag = hasUser
				msg.err = err
				dock.replyMessage(msg)
			})
		}else{
			msg.flag = false
			msg.err = NoFoundUser
			dock.replyMessage(msg)
		}
	}
	return
}
func (dock *dock) AddUser(iUser user.IUser) {
	fa:=dock.UserStore.AddUser(iUser)
	log.InfoF("添加用户：{}  flag:{}",iUser.GetUsername(),fa)
	if fa{
		if dock.handleAddUser!=nil{
			dock.handleAddUser(iUser)
		}
	}
}
func (dock *dock) DeleteUser(iUser user.IUser) {
	fa:=dock.UserStore.DeleteUser(iUser)
	log.InfoF("删除用户：{}  flag:{}",iUser.GetUsername(),fa)
	if fa{
		if dock.handleDeleteUser!=nil{
			dock.handleDeleteUser(iUser.GetUsername())
		}
	}
}
func (dock *dock) UserNum()int32 {
	return dock.UserStore.GetUserNum()
}

func (dock *dock) replyMessage(msg *DockMessage) {
	log.DebugF("加入消息反馈队列:{}",msg.InputMessage.GetMessageId())
	dock.replyMsg.Offer(msg)
}
func (dock *dock) exchangeReplyMsg() {
	log.DebugF("启动信息反馈处理")
	for {
		msg,_ := dock.replyMsg.Poll()
		dockMessage:=msg.(*DockMessage)
		log.DebugF("处理反馈信息：{}",dockMessage.InputMessage.GetMessageId())
		if msg != nil {
			dockMessage.write(dockMessage.err, dockMessage.flag)
		}
		log.DebugF("处理反馈信息：{} 完成",dockMessage.InputMessage.GetMessageId())
	}
}
func (dock *dock) sendNum()int32{

	return dock.sendMsg.Num()
}
func (dock *dock) replyNum()int32{

	return dock.replyMsg.Num()
}
func (dock *dock) exchangeSendMsg() {
	log.DebugF("启动信息发送处理")

	for {
		msg,_ := dock.sendMsg.Poll()
		if msg != nil {
			dm:=msg.(*DockMessage)
			fa, err := dock.writeUserMsg(dm)
			if !dm.IsForward || fa{
				log.DebugF("fa:{} err:{}", fa, err)
				dm.flag = fa
				dm.err = err
				dock.replyMessage(dm)
			}
		}
	}
}
