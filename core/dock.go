package core

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
)

type dock struct {
	sendMsg   *Queue
	UserStore *user.Store
}

func newDock() *dock {
	return &dock{sendMsg: NewQueue(), UserStore: user.NewStore()}
}

func (dock *dock) sendMessage(iMessage message.IMessage) error {
	mc := dock.sendMsg.offer(iMessage)
	fa, err := mc.wait()
	if err != nil {
		return err
	}
	if !fa {
		return NoFoundUser
	}
	return nil
}

func (dock *dock) writeUserMsg(msg *messageChan) (bool, error) {
	userId := msg.inputMessage.GetString(message.ToUser)
	log.InfoF("信息发送:{}", userId)
	var ee error = nil
	return dock.UserStore.GetUser(msg.inputMessage.GetString(message.ToUser), func(iUser user.IUser) bool {
		err := iUser.WriteMessage(msg.inputMessage)
		ee = err
		return err != nil
	}), ee
}
func (dock *dock)login(iMessage message.IMessage, writeRead user.IUser){
	log.DebugF("添加新用户 :{}",writeRead.GetUsername())
	dock.UserStore.AddUser(writeRead)
}
func (dock *dock)DeleteUser(iUser user.IUser){

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
func (dock *dock) exchangeSendMsg() {
	log.Info("启动信息发送处理")
	var i = 0
	for {
		msg := dock.sendMsg.poll()
		if msg != nil {
			log.InfoF("收到message msgId:", msg.inputMessage.GetMessageId())
			fa, err := dock.writeUserMsg(msg)
			log.InfoF("fa:{} err:{}", fa, err)
			msg.receive(fa, err)
		}
		i++
		if i>>10 == 1 {
			i = 0
			log.InfoF("当前信息池剩下 :{} 未处理", dock.sendMsg.Num())
		}
	}
}
