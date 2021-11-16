package km

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"io"
	"sync"
)
var chunkWriteStreamPool *sync.Pool
var chunkReadStreamPool *sync.Pool
func init() {
	chunkWriteStreamPool = &sync.Pool{New: func() interface{} {
		return new(chunkWriteStream)
	}}
	chunkReadStreamPool = &sync.Pool{New: func() interface{} {
		return new(chunkReadStream)
	}}
}

func getChunkWriteStreamPool(msg message.IMessage)*chunkWriteStream  {
	var chunkWriteStream = chunkWriteStreamPool.Get().(*chunkWriteStream)
	chunkWriteStream.maxBodySize = 512
	chunkWriteStream.chunkId = getChunkId()
	chunkWriteStream.message = msg
	chunkWriteStream.process = 0
	chunkWriteStream.keyIndex = 0
	chunkWriteStream.messageLength = msg.GetMessageLength()
	chunkWriteStream.rMessageLength = 0
	return chunkWriteStream
}
func getChunkReadStream(read_ io.Reader)* chunkReadStream {
	var chunkReadStream = chunkReadStreamPool.Get().(*chunkReadStream)
	chunkReadStream.read_ = net.NewIOReadStream(read_)
	chunkReadStream.maxBodySize = 512
	chunkReadStream.recordMap = make(map[uint16]*chunkRecord)
	return chunkReadStream
}
func putChunkWriteStream(crm * chunkWriteStream )  {
	chunkWriteStreamPool.Put(crm)
}
func putChunkReadStream(crm * chunkReadStream )  {
	chunkReadStreamPool.Put(crm)
}

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


