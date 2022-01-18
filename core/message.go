package core

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
)

const (
	wait int8 = iota
	run
)

type DockMessage struct {
	InputMessage message.IMessage
	write        user.WriteFunc
	flag         bool
	err          error
	IsForward    bool
}

func (dm *DockMessage) GetToUsername() string {
	return dm.InputMessage.GetString(message.ToUser)
}
func newDockMessage(inputMessage message.IMessage, write user.WriteFunc) *DockMessage {
	return &DockMessage{InputMessage: inputMessage, flag: false, write: write}
}

