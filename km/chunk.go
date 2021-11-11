package km

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/util"
)

type chunk0 struct {
	*Chunk
	messageHeader []byte
	key           byte
	dataLen       int
	data          []byte
}

func createChunk0(classType byte, messageLength uint32, messageId uint32, time uint32, key byte, dataLen int, data []byte) *chunk0 {

	return &chunk0{Chunk: &Chunk{classType}, messageHeader: messageHeader(messageLength, messageId, time), key: key, dataLen: dataLen, data: data}
}

func messageHeader(messageLength uint32, messageId uint32, time uint32) []byte {

	if messageLength-32_767 < 0 {
		b := []byte{0, 0}
		b[0] = byte(messageLength)
		b[1] = byte(messageLength >> 8)
	} else {
		b := []byte{0, 0, 0, 0}
		var pre = (messageLength) | 32_768
		b[0] = byte(pre)
		b[1] = byte(pre >> 8)
		b[2] = byte(messageLength >> 15)
		b[3] = byte(messageLength >> 23)
	}

	util.U32TOBytes(messageId)

	return nil
}

type chunk1 struct {
	*Chunk
	key     byte
	dataLen int
	data    []byte
}

func createChunk1(classType byte, key byte, dataLen int, data []byte) *chunk1 {

	return &chunk1{Chunk: &Chunk{64 | classType}}
}

type chunk2 struct {
	*Chunk
	data []byte
}

func createChunk2(classType byte, data []byte) *chunk2 {
	return &chunk2{Chunk: &Chunk{128 | classType}}
}

type Chunk struct {
	chunkHeader byte
}

func (chunk *Chunk) setChunkHeader(chunkType byte, classType byte) {
	chunk.chunkHeader = chunkType<<6 | classType
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
