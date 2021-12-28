package km

import (
	log "github.com/chuccp/coke-log"
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
	chunkStream := createChunkReadStream(km.io)
	msg, err := chunkStream.readMessage()
	freeChunkReadStream(chunkStream)
	return msg, err
}
func (km *km00001) WriteMessage(msg message.IMessage) error {
	log.DebugF("WriteMessage 写入信息class:{} type:{} msgId:{}",msg.GetClassId(),msg.GetMessageType(),msg.GetMessageId())
	chunkStream := createChunkWriteStreamPool(msg)
	for chunkStream.hasNext() {
		chunk := chunkStream.readChunk()
		err := km.writeChunk(chunk)
		if err != nil {
			freeChunkWriteStream(chunkStream)
			return err
		}
	}
	freeChunkWriteStream(chunkStream)
	return nil
}
func (km *km00001) writeChunk(chunk IChunk) error {
	_, err := km.io.Write(chunk.toByte())
	return err
}
