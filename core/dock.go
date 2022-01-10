package core

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
)

type WriteFunc func(iMessage message.IMessage,err error,hasUser bool)

type dock struct {
	sendMsg   *Queue
	replyMsg *Queue
	UserStore *user.Store
}

func newDock() *dock {
	return &dock{sendMsg: NewQueue(), UserStore: user.NewStore(),replyMsg:NewQueue()}
}

func (dock *dock) sendMessage(iMessage message.IMessage,write WriteFunc)  {
	 dock.sendMsg.offer(newDockMessage(iMessage,write))
}

func (dock *dock) writeUserMsg(msg *dockMessage) (bool, error) {
	userId := msg.inputMessage.GetString(message.ToUser)
	log.DebugF("信息发送:{}", userId)
	var ee error = nil
	return dock.UserStore.GetUser(msg.inputMessage.GetString(message.ToUser), func(iUser user.IUser) bool {
		err := iUser.WriteMessage(msg.inputMessage)
		ee = err
		return err != nil
	}), ee
}
func (dock *dock)login(iMessage message.IMessage, writeRead user.IUser){
	writeRead.SetUsername(iMessage.GetString(message.Username))
	log.DebugF("添加新用户 :{}",writeRead.GetUsername())
	if writeRead.GetUsername()==""{
		log.ErrorF("用户名不能为空")
		return
	}
	dock.UserStore.AddUser(writeRead)
}
func (dock *dock)DeleteUser(iUser user.IUser){
	dock.UserStore.DeleteUser(iUser)
}
func (dock *dock) handleMessage(iMessage message.IMessage, writeRead user.IUser) {

	switch iMessage.GetClassId() {

	case message.FunctionMessageClass:
		switch iMessage.GetMessageType() {
				case message.LoginType:
					dock.login(iMessage,writeRead)

		}

	}

}
func (dock *dock)  replyMessage(msg *dockMessage){
	dock.replyMsg.offer(msg)
}
func (dock *dock) exchangeReplyMsg(){
	log.Info("启动信息反馈处理")
	var i = 0
	for {
		msg := dock.replyMsg.poll()
		if msg != nil {
			msg.write(msg.inputMessage,msg.err,msg.flag)
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
			log.InfoF("收到message msgId:", msg.inputMessage.GetMessageId())
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
