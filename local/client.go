package local

import (
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
)

type Client struct {
	user   *User
	context *core.Context
}

func newClient(user *User,context *core.Context)*Client  {
	return &Client{user:user,context:context}
}

func (client *Client)login()  {
	client.context.AddUser(client.user)
}

func (client *Client)handleMessage(msg message.IMessage)error{
	switch msg.GetClassId() {
	case message.OrdinaryMessageClass:
		msg.SetString(message.FromUser, client.user.GetUsername())

	case message.FunctionMessageClass:

	}
	return nil
}