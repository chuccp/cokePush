package message

type IMessage interface {
	GetMessageId() int
	GetTimestamp() int
	GetMessageLength() int
	GetMessageType() byte
	GetClassId() byte
	GetValue(key int) interface{}
	GetKeys() []byte
	SetValue(key int, value interface{})
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
func (message *Message) GetValue(key int) interface{} {
	return nil
}
func (message *Message) GetString(key int) string {
	return ""
}
func (message *Message) SetValue(key int, value interface{}) {

}
func (message *Message) GetKeys() []byte {
	return nil
}
