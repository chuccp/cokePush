package km

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
)

type km interface {
	ReadMessage() (message.IMessage, error)
	WriteMessage(msg message.IMessage) error
}

type km00001 struct {
	readWrite *net.IONetStream
}

func NewKm00001(readWrite *net.IONetStream) *km00001 {
	return &km00001{readWrite: readWrite}
}
func (km *km00001) ReadMessage() (message.IMessage, error) {
	chunkStream := createChunkReadStream(km.readWrite.IOReadStream)
	msg, err := chunkStream.readMessage()
	freeChunkReadStream(chunkStream)
	return msg, err
}
func (km *km00001) WriteMessage(msg message.IMessage) error {
	chunkStream := createChunkWriteStreamPool(msg)
	chunk := chunkStream.readChunk0()
	err := km.writeChunk(chunk)
	if err != nil {
		freeChunkWriteStream(chunkStream)
		return err
	}
	for chunkStream.hasNext() {
		chunk = chunkStream.readChunk()
		err := km.writeChunk(chunk)
		if err != nil {
			freeChunkWriteStream(chunkStream)
			return err
		}
	}
	freeChunkWriteStream(chunkStream)
	km.readWrite.IOWriteStream.Flush()
	return nil
}
func (km *km00001) writeChunk(chunk IChunk) error {
	_, err := km.readWrite.IOWriteStream.Write(chunk.toByte())
	return err
}
