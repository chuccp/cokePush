package local

import (
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
)

type Client struct {
	user    *User
	context *core.Context
}

func newClient(user *User, context *core.Context) *Client {
	return &Client{user: user, context: context}
}

func (client *Client) login() error {
	client.context.AddUser(client.user)
	return nil
}

func (client *Client) handleMessage(msg message.IMessage) error {
	switch msg.GetClassId() {
	case message.OrdinaryMessageClass:
		msg.SetString(message.FromUser, client.user.GetUsername())
		return client.context.SendMessage(msg)
	case message.FunctionMessageClass:

	}
	return core.UnKnownClassId
}
