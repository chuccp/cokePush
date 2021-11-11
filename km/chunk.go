package km

import "github.com/chuccp/cokePush/message"

type chunk0 struct {
	*Chunk
	messageHeader []byte
}

func newChunk0(messageType byte, messageLength int, messageId int, time int) *chunk0 {
	return &chunk0{}
}

type chunk1 struct {
	*Chunk
	key  byte
	data []byte
}
type chunk2 struct {
	*Chunk
	data []byte
}
type Chunk struct {
	chunkHeader byte
}

func (chunk *Chunk) chunkType() byte {
	return chunk.chunkHeader >> 6
}
func (chunk *Chunk) classType() byte {
	return chunk.chunkHeader << 2 >> 2
}

type IChunk interface {
	chunkType() byte
	classType() byte
}

type chunkStream struct {
	message message.IMessage
	process int
}

func newChunkStream(message message.IMessage) *chunkStream {
	return &chunkStream{message: message, process: 0}
}
func (stream *chunkStream) hasNext() bool {

	return true
}
func (stream *chunkStream) readChunk() IChunk {
	if stream.process == 0 {
		stream.process = 1
		return newChunk0(stream.message.GetMessageType(),0,0,0)
	}else if stream.process==1{

	}
	return nil
}
