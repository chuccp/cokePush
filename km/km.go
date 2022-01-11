package km

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"sync"
)

type km interface {
	ReadMessage() (message.IMessage, error)
	WriteMessage(msg message.IMessage) error
}

type km00001 struct {
	readWrite *net.IONetStream
	lock *sync.RWMutex
}

func NewKm00001(readWrite *net.IONetStream) *km00001 {
	return &km00001{readWrite: readWrite,lock:new(sync.RWMutex)}
}
func (km *km00001) ReadMessage() (message.IMessage, error) {
	chunkStream := createChunkReadStream(km.readWrite.IOReadStream)
	msg, err := chunkStream.readMessage()
	freeChunkReadStream(chunkStream)
	return msg, err
}
func (km *km00001) WriteMessage(msg message.IMessage) (err error) {
	km.lock.Lock()
	defer km.lock.Unlock()
	chunkStream := createChunkWriteStreamPool(msg)
	chunk := chunkStream.readChunk0()
	err = km.writeChunk(chunk)
	if err != nil {
		freeChunkWriteStream(chunkStream)
		return
	}
	for chunkStream.hasNext() {
		chunk = chunkStream.readChunk()
		err = km.writeChunk(chunk)
		if err != nil {
			freeChunkWriteStream(chunkStream)
			return
		}
	}
	freeChunkWriteStream(chunkStream)
	if err==nil{
		km.readWrite.IOWriteStream.Flush()
	}
	return
}
func (km *km00001) writeChunk(chunk IChunk) error {
	_, err := km.readWrite.IOWriteStream.Write(chunk.toByte())
	return err
}
