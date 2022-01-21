package core

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
	"sync"
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
	replay    bool
}

func (dm *DockMessage) GetToUsername() string {
	return dm.InputMessage.GetString(message.ToUser)
}


var poolDockMessage = &sync.Pool{
	New: func() interface{} {
		return new(DockMessage)
	},
}

func newDockMessage(inputMessage message.IMessage, write user.WriteFunc) *DockMessage {
	dk:=poolDockMessage.Get().(*DockMessage)
	dk.InputMessage = inputMessage
	dk.flag = false
	dk.write =  write
	dk.replay = true
	return dk

}
func newDockMessageNoReplay(inputMessage message.IMessage) *DockMessage {
	dk:=poolDockMessage.Get().(*DockMessage)
	dk.InputMessage = inputMessage
	dk.flag = false
	dk.replay = false
	return dk
}
func free(dk *DockMessage)  {
	poolDockMessage.Put(dk)
}