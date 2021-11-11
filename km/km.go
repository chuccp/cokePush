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

	return nil, nil
}
func (km *km00001) WriteMessage(msg message.IMessage) error {
	chunkStream := newChunkStream(msg)
	for chunkStream.hasNext() {
		chunk := chunkStream.readChunk()
		err := km.writeChunk(chunk)
		if err != nil {
			return err
		}
	}
	return nil
}
func (km *km00001) writeChunk(chunk IChunk) error {



	return nil
}
