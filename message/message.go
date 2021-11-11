package message

type IMessage interface {
	GetMessageId() int
	GetTimestamp() int
	GetMessageLength() int
	GetMessageType() byte
	GetClassId() byte
	GetValue(key byte) []byte
	GetKeys() []byte
	SetValue(key byte, value interface{})
}

type Message struct {
}

func (message *Message) GetMessageId() int {
	return 0
}
func (message *Message) GetTimestamp() int {
	return 0
}
func (message *Message) GetMessageLength() int {
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
