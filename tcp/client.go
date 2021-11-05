package tcp

import (
	"github.com/chuccp/cokePush/km"
	"github.com/chuccp/cokePush/net"
)

type Client struct {
	stream *km.Stream
}

func NewClient(stream *net.IOStream) (*Client, error) {
	kmStream, err := km.NewStream(stream)
	return &Client{stream: kmStream}, err
}

func (client *Client) Start() {

	msg,err := client.stream.ReadMessage()
	if err==nil{
		type_, err := msg.GetMessageType()
		if err == nil {
			if type_ == km.BasicMessage {



			}
		}
	}
}
