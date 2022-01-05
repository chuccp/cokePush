package core

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
)

type dock struct {
	sendMsg       *Queue
	writeMsg      *Queue
	UserStore     *user.Store
}

func newDock() *dock {
	return &dock{sendMsg: NewQueue(), writeMsg: NewQueue(),UserStore: user.NewStore()}
}

func (dock *dock)sendMessage(iMessage message.IMessage)error  {
	dock.sendMsg.offer(iMessage)
	return nil
}

func (dock *dock)writeUserMsg(msg message.IMessage) (bool,error) {
	userId := msg.GetString(message.ToUser)
	log.InfoF("信息发送:{}", userId)
	var ee error = nil
	return dock.UserStore.GetUser(msg.GetString(message.ToUser), func(iUser user.IUser) bool {
		err:=iUser.WriteMessage(msg)
		ee=err
		return err!=nil
	}),ee
}

func (dock *dock)handleMessage(iMessage message.IMessage,writeRead user.IUser)  {


}
func (dock *dock) exchangeSendMsg() {
	log.Info("启动信息发送处理")
	var i = 0
	for {
		msg := dock.sendMsg.poll()
		if msg != nil {
			//dock.writeUserMsg(msg)
		}
		i++
		if i>>10 == 1 {
			i = 0
			log.InfoF("当前信息池剩下 :{} 未处理", dock.sendMsg.Num())
		}
	}
}