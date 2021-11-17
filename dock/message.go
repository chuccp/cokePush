package dock

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
)
var handleMessage = &MessageDock{}
var handleFunction = &FunctionDock{}
type MessageDock struct {
}

func (messageDock *MessageDock) handleMessage(iMessage message.IMessage) {

	switch iMessage.GetMessageType() {
	case message.BasicMessageType:
		messageDock.basicMessage(iMessage)
	case message.GroupMessageType:
		messageDock.groupMessage(iMessage)
	case message.BroadcastMessageType:
		messageDock.broadcastMessage(iMessage)
	}

}
func (messageDock *MessageDock) basicMessage(iMessage message.IMessage) {
	toUser:=iMessage.GetString(message.ToUser)
	iu:=user.GetUser(toUser)
	iu.WriteMessage(iMessage)
}
func (messageDock *MessageDock) groupMessage(iMessage message.IMessage) {

	//groupId:=iMessage.GetValue(message.Group)



}
func (messageDock *MessageDock) broadcastMessage(iMessage message.IMessage) {

	user.Range(func(iUser user.IUser) bool {
		iUser.WriteMessage(iMessage)
		return true
	})
}

type FunctionDock struct {
}
func (function *FunctionDock) handleFunction(iMessage message.IMessage,iUser user.IUser) {
	switch iMessage.GetMessageType() {
	case message.LoginType:
		function.handleLogin(iMessage,iUser)
	case message.JoinGroupType:
		function.joinGroup(iMessage,iUser)
	}
}

func (function *FunctionDock) handleLogin(iMessage message.IMessage,iUser user.IUser) {
	username:=iMessage.GetString(message.Username)
	iUser.SetUsername(username)
	user.AddUser(iUser)
}
func (function *FunctionDock) joinGroup(iMessage message.IMessage,iUser user.IUser) {
	groupId:=iMessage.GetString(message.GroupId)
	user.JoinGroup(groupId,iUser)
}

func OnMessage(iMessage message.IMessage) {
	handleMessage.handleMessage(iMessage)
}
func OnFunction(iMessage message.IMessage,iUser user.IUser) {
	handleFunction.handleFunction(iMessage,iUser)
}
