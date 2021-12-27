package cluster

import (
	"github.com/chuccp/cokePush/km"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
)

type Client struct {
	stream *km.Stream
}

func (client *Client) Start() {

	for {
		msg, err := client.stream.ReadMessage()
		if err == nil {
			client.handleMessage(msg)
		} else {
			break
		}
	}

}

func (client *Client) handleMessage(msg message.IMessage) {
	switch msg.GetClassId() {

	case message.FunctionMessageClass:

	case message.LiveMessageClass:

	}

}
func NewClient(stream *net.IONetStream) (*Client, error) {
	kmStream, err := km.NewStream(stream)
	return &Client{stream: kmStream}, err
}
