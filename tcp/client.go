package tcp

import (
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/km"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"strconv"
	"unsafe"
)

type client struct {
	stream *km.Stream
	context *core.Context
	id string
	username string
}
func NewClient(stream *net.IONetStream,context *core.Context) (*client, error) {
	kmStream, err := km.NewStream(stream)
	client:=&client{stream: kmStream,context:context }
	return client, err
}
func (client *client) Start() {
	msg, err := client.stream.ReadMessage()
	if err != nil {
		client.context.Handle(msg,client)
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