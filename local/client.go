package local

import (
	"github.com/chuccp/cokePush/dock"
	"github.com/chuccp/cokePush/message"
)

type Client struct {
	user   *User
}

func newClient(user *User)*Client  {
	return &Client{user}
}

func (client *Client)login()  {

}

func (client *Client)handleMessage(msg message.IMessage)error{
	switch msg.GetClassId() {
	case message.OrdinaryMessageClass:
		msg.SetString(message.FromUser, client.user.GetUsername())
		dock.OnMessage(msg)
	case message.FunctionMessageClass:
		dock.OnFunction(msg, client.user)
	}
	return nil
}