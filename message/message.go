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
	SetString(key byte, value string)
	SetValue(key byte, value []byte)
}

type Message struct {
	messageType   byte
	classId       byte
	keys          []byte
	data          map[byte][]byte
	messageLength uint32
	time          uint32
	messageId     uint32
}

func CreateMessage() *Message {

	return &Message{messageId: msgId()}
}
func (message *Message) GetMessageId() uint32 {
	return 0
}
func (message *Message) GetTimestamp() uint32 {
	return 0
}
func (message *Message) GetMessageLength() uint32 {
	return message.messageLength
}
func (message *Message) GetClassId() byte {
	return message.classId
}
func (message *Message) GetMessageType() byte {
	return message.messageType
}
func (message *Message) GetValue(key byte) []byte {
	return message.data[key]
}
func (message *Message) GetString(key byte) string {
	return ""
}
func (message *Message) SetString(key byte, value string) {
	message.SetValue(key, []byte(value))
}
func (message *Message) SetValue(key byte, value []byte) {
	message.keys = append(message.keys, key)
	message.data[key] = value
	message.messageLength = message.messageLength + uint32(len(value))
}
func (message *Message) GetKeys() []byte {
	return nil
}
func msgId() uint32 {
	num := rand.Intn(1024)
	return util.Millisecond()<<10 | (uint32(num))
}
