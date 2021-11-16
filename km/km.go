package km

import (
	"github.com/chuccp/cokePush/message"
	"io"
)


type km interface {
	ReadMessage() (message.IMessage, error)
	WriteMessage(msg message.IMessage) error
}

type km00001 struct {
	io io.ReadWriter
}

func NewKm00001(io io.ReadWriter) *km00001 {
	return &km00001{io: io}
}
func (km *km00001) ReadMessage() (message.IMessage, error) {
	chunkStream := getChunkReadStream(km.io)
	msg,err:= chunkStream.readMessage()
	putChunkReadStream(chunkStream)
	return msg, err
}
func (km *km00001) WriteMessage(msg message.IMessage) error {
	chunkStream := getChunkWriteStreamPool(msg)
	for chunkStream.hasNext() {
		chunk := chunkStream.readChunk()
		err := km.writeChunk(chunk)
		if err != nil {
			putChunkWriteStream(chunkStream)
			return err
		}
	}
	putChunkWriteStream(chunkStream)
	return nil
}
func (km *km00001) writeChunk(chunk IChunk) error {
	_,err:=km.io.Write(chunk.toByte())
	return err
}


