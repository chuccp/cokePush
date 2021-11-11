package message

import (
	"github.com/chuccp/cokePush/util"
	"math/rand"
)

type IMessage interface {
	GetMessageId() uint32
	GetTimestamp() uint32
	GetMessageLength() uint32
	GetMessageType() byte
	GetClassId() byte
	GetValue(key byte) []byte
	GetKeys() []byte
	SetValue(key byte, value interface{})
}

type Message struct {
}

func (message *Message) GetMessageId() uint32 {
	return 0
}
func (message *Message) GetTimestamp() uint32 {
	return 0
}
func (message *Message) GetMessageLength() uint32 {
	return 0
}
func (message *Message)GetClassId() byte{
	return 0
}
func (message *Message) GetMessageType() byte {
	return 0
}
func (message *Message) GetValue(key byte) []byte {
	return nil
}
func (message *Message) GetString(key byte) string {
	return ""
}
func (message *Message) SetValue(key byte, value interface{}) {

}
func (message *Message) GetKeys() []byte {
	return nil
}
func GetMsgId() uint32 {
	num := rand.Intn(1024)
	return util.Millisecond()<<10 | (uint32(num))
}