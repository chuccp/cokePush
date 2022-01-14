package core

import (
	"github.com/chuccp/cokePush/message"
)

const (
	wait int8 = iota
	run
)

type DockMessage struct {
	InputMessage message.IMessage
	write        WriteFunc
	flag         bool
	err          error
	IsForward    bool
}

func (dm *DockMessage) GetToUsername() string {
	return dm.InputMessage.GetString(message.ToUser)
}
func newDockMessage(inputMessage message.IMessage, write WriteFunc) *DockMessage {
	return &DockMessage{InputMessage: inputMessage, flag: false, write: write}
}

