package tcp

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/km"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/user"
	"strconv"
	"unsafe"
)

type client struct {
	stream *km.Stream
	context *core.Context
	id string
	username string
}

func (client *client) GetRemoteAddress() string {
	return client.stream.RemoteAddr().String()
}

func NewClient(stream *net.IONetStream,context *core.Context) (*client, error) {
	kmStream, err := km.NewStream(stream)
	client:=&client{stream: kmStream,context:context }
	return client, err
}
func (client *client) Start() {
	msg, err := client.stream.ReadMessage()
	if err != nil {
		client.handle(msg,client)
	}
}

func (client *client) handle(msg message.IMessage,writeRead user.IUser)error{
	client.handleMessage(msg,writeRead)
	return nil
}

func (client *client) handleMessage(iMessage message.IMessage, iUser user.IUser) {
	switch iMessage.GetClassId() {
	case message.FunctionMessageClass:
		switch iMessage.GetMessageType() {
		case message.LoginType:
			client.login(iMessage, iUser)

		}
	}
}
func (client *client)login(iMessage message.IMessage, iUser user.IUser){
	iUser.SetUsername(iMessage.GetString(message.Username))
	log.DebugF("添加新用户 :{}", iUser.GetUsername())
	if iUser.GetUsername()==""{
		log.ErrorF("用户名不能为空")
		return
	}else{
		client.context.AddUser(iUser)
	}

}
func (client *client)WriteMessage(iMessage message.IMessage) error{
	return client.stream.WriteMessage(iMessage)
}
func (client *client)GetId() string{
	return client.id
}
func (client *client)GetUsername() string{
	return client.username
}
func (client *client)SetUsername(username string){
	client.username = username
	client.id = username+strconv.FormatUint(uint64(uintptr(unsafe.Pointer(client))), 36)
}