package message

import (
	"github.com/chuccp/cokePush/util"
	"math/rand"
	"time"
)

type IMessage interface {
	GetMessageId() uint32
	GetTimestamp() uint32
	GetMessageLength() uint32
	GetMessageType() byte
	GetClassId() byte
	GetValue(key byte) []byte
	GetString(key byte) string

	GetKeys() []byte
	SetString(key byte, value string)
	SetValue(key byte, value []byte)

	SetMessageId(messageId uint32)
	SetTimestamp(time uint32)
	SetMessageLength(msgLen uint32)
	SetMessageType(messageType byte)
	SetClassId(classId byte)
}

type Message struct {
	classId       byte
	messageType   byte
	keys          []byte
	data          map[byte][]byte
	messageLength uint32
	time          uint32
	messageId     uint32
}

func NewMessage() *Message {
	return &Message{data: make(map[byte][]byte)}
}
func CreateMessage(classId byte, messageType byte) *Message {
	return CreateBackMessage(classId,messageType,MsgId())
}
func CreateBackMessage(classId byte, messageType byte,msgId uint32) *Message {
	return &Message{messageId: msgId, time: millisecond(), classId: classId, messageType: messageType, keys: make([]byte, 0), data: make(map[byte][]byte)}
}
func (message *Message) GetMessageId() uint32 {
	return message.messageId
}
func (message *Message) GetTimestamp() uint32 {
	return message.time
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
	return string(message.data[key])
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
	return message.keys
}

func (message *Message) SetMessageId(messageId uint32) {
	message.messageId = messageId
}
func (message *Message) SetTimestamp(time uint32) {
	message.time = time
}
func (message *Message) SetMessageLength(msgLen uint32) {
	message.messageLength = msgLen
}
func (message *Message) SetMessageType(messageType byte) {
	message.messageType = messageType
}
func (message *Message) SetClassId(classId byte) {
	message.classId = classId
}

func MsgId() uint32 {
	num := rand.Intn(1024)
	return util.Millisecond()<<10 | (uint32(num))
}
func millisecond() uint32 {
	ms := time.Now().UnixNano() / 1e6
	return uint32(ms)
}
