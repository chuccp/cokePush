package km

import "github.com/chuccp/cokePush/message"

type chunk0 struct {
	*Chunk
	messageHeader []byte
	key           byte
	dataLen       int
	data          []byte
}

func createChunk0(messageType byte, messageLength int, messageId int, time int, key byte, dataLen int, data []byte) *chunk0 {

	return &chunk0{}
}

type chunk1 struct {
	*Chunk
	key     byte
	dataLen int
	data    []byte
}

func createChunk1(messageType byte, key byte, dataLen int, data []byte) *chunk1 {

	return &chunk1{}
}

type chunk2 struct {
	*Chunk
	data []byte
}

func createChunk2(messageType byte, data []byte) *chunk2 {

	return &chunk2{}
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
	message  message.IMessage
	process  int
	dataLen  int
	keyIndex byte
}

func newChunkStream(message message.IMessage) *chunkStream {
	return &chunkStream{message: message, process: 0, keyIndex: 0}
}
func (stream *chunkStream) hasNext() bool {

	return true
}
func (stream *chunkStream) readChunk() IChunk {
	key := stream.message.GetKeys()[stream.keyIndex]
	data := stream.message.GetValue(key)
	start := stream.dataLen
	end := start + 512
	if end > stream.dataLen-start {
		end = stream.dataLen - start
		stream.dataLen = end
	}
	if stream.process == 0 {
		stream.process = 1
		return createChunk0(stream.message.GetMessageType(), stream.message.GetMessageLength(), stream.message.GetMessageId(), stream.message.GetTimestamp(), key, len(data), data[start:end])
	} else if stream.process == 1 {

		chunk := createChunk1(stream.message.GetMessageType(), key, len(data), data[start:end])
		if stream.dataLen > 0 {
			stream.process = 2
		}
		stream.keyIndex++
		return chunk
	} else if stream.process == 2 {
		chunk := createChunk2(stream.message.GetMessageType(), data[start:end])
		if stream.dataLen == 0 {
			stream.process = 1
		}
		return chunk
	}
	return nil
}
