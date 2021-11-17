package tcp

import (
	"github.com/chuccp/cokePush/dock"
	"github.com/chuccp/cokePush/km"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
)

type Client struct {
	stream *km.Stream
	user   *User
}

func NewClient(stream *net.IONetStream) (*Client, error) {
	kmStream, err := km.NewStream(stream)
	return &Client{stream: kmStream, user: NewUser()}, err
}

func (client *Client) Start() {
	msg, err := client.stream.ReadMessage()
	if err == nil {
		switch msg.GetClassId() {
		case message.OrdinaryMessageClass:
			msg.SetString(message.FromUser, client.user.GetUsername())
			dock.OnMessage(msg)
		case message.FunctionMessageClass:
			dock.OnFunction(msg, client.user)
		}
	}
}
