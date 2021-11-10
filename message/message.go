package message

type IMessage interface {
	GetMessageId() int
	GetTimestamp() int
	GetMessageLength() int
	GetMessageType() int8
	GetClassId() int8
	GetValue(key int) interface{}
	GetKeys() []int8
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
func (message *Message)GetClassId() int8{
	return 0
}
func (message *Message) GetMessageType() int8 {
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
func (message *Message) GetKeys() []int8 {
	return nil
}
