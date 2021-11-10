package km

import "github.com/chuccp/cokePush/message"

type km interface {
	ReadMessage() (message.IMessage, error)
	WriteMessage(msg message.IMessage ) error
}

type km00001 struct {
}



func NewKm00001() *km00001 {
	return &km00001{}
}
func (km *km00001) ReadMessage() (message.IMessage, error) {

	return nil, nil
}
func (km *km00001) WriteMessage(msg message.IMessage ) error{
		newChunkStream(msg)
	return nil
}
