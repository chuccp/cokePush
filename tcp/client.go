package tcp

import (
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/km"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
)

type Client struct {
	stream *km.Stream
	context *core.Context
}
func NewClient(stream *net.IONetStream,context *core.Context) (*Client, error) {
	kmStream, err := km.NewStream(stream)
	return &Client{stream: kmStream,context:context }, err
}
func (client *Client) Start() {
	msg, err := client.stream.ReadMessage()
	if err != nil {
		client.context.Handle(msg,client)
	}
}
func (client *Client)WriteMessage(iMessage message.IMessage) error{
	return client.stream.WriteMessage(iMessage)
}
func (client *Client)ReadMessage() (message.IMessage,error){
	return client.stream.ReadMessage()
}
func (client *Client)GetUserId() string{

	return ""
}
func (client *Client)GetUsername() string{
	return ""
}
func (client *Client) Close() {

}
