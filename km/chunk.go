package km

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/util"
)

type chunk0 struct {
	*chunk
	classType     byte
	messageType   byte
	messageLength uint32
	messageId     uint32
	time          uint32
	key           byte
	dataLen       uint32
	data          []byte
}

func createChunk0(chunkId byte, classType byte, messageType byte, messageLength uint32, messageId uint32, time uint32, key byte, dataLen uint32, data []byte) *chunk0 {
	return &chunk0{chunk: &chunk{chunkId}, messageType: messageType, messageLength: messageLength, messageId: messageId, time: time, key: key, dataLen: dataLen, data: data}
}

func (chunk0 chunk0) toByte() []byte {
	bytesArray := make([]byte, 0)
	//写chunk
	bytesArray = append(bytesArray, chunk0.chunk.toByte()...)
	//写header
	bytesArray = append(bytesArray, chunk0.classType)
	bytesArray = append(bytesArray, chunk0.messageType)
	bytesArray = append(bytesArray, util.U32TOBytes(chunk0.time)...)
	bytesArray = append(bytesArray, util.U32TOBytes(chunk0.messageId)...)
	bytesArray = append(bytesArray, lengthToBytes(chunk0.messageLength)...)
	//写body
	bytesArray = append(bytesArray, chunk0.key)
	bytesArray = append(bytesArray, lengthToBytes(chunk0.dataLen)...)
	bytesArray = append(bytesArray, chunk0.data...)
	return bytesArray
}

func lengthToBytes(length uint32) []byte {
	if length-32_767 < 0 {
		b := []byte{0, 0}
		b[0] = byte(length)
		b[1] = byte(length >> 8)
		return b
	} else {
		b := []byte{0, 0, 0, 0}
		var pre = (length) | 32_768
		b[0] = byte(pre)
		b[1] = byte(pre >> 8)
		b[2] = byte(length >> 15)
		b[3] = byte(length >> 23)
		return b
	}
}

type chunk1 struct {
	*chunk
	key     byte
	dataLen uint32
	data    []byte
}

func createChunk1(classType byte, key byte, dataLen uint32, data []byte) *chunk1 {

	return &chunk1{chunk: &chunk{64 | classType}}
}

func (chunk1 chunk1) toByte() []byte {
	bytesArray := make([]byte, 0)
	bytesArray = append(bytesArray, chunk1.chunk.toByte()...)
	//写body
	bytesArray = append(bytesArray, chunk1.key)
	bytesArray = append(bytesArray, lengthToBytes(chunk1.dataLen)...)
	bytesArray = append(bytesArray, chunk1.data...)

	return bytesArray
}

type chunk2 struct {
	*chunk
	data []byte
}

func createChunk2(chunkId byte, data []byte) *chunk2 {
	return &chunk2{chunk: &chunk{128 | chunkId}}
}
func (chunk2 chunk2) toByte() []byte {
	bytesArray := make([]byte, 0)
	bytesArray = append(bytesArray, chunk2.chunk.toByte()...)
	bytesArray = append(bytesArray, chunk2.data...)
	return bytesArray
}

type chunk struct {
	chunkHeader byte
}

func (chunk *chunk) setChunkHeader(chunkType byte, chunkId byte) {
	chunk.chunkHeader = chunkType<<6 | chunkId
}
func (chunk *chunk) chunkType() byte {
	return chunk.chunkHeader >> 6
}
func (chunk *chunk) chunkId() byte {
	return chunk.chunkHeader << 2 >> 2
}
func (chunk *chunk) toByte() []byte {
	return nil
}

type IChunk interface {
	chunkType() byte
	chunkId() byte
	toByte() []byte
}

type chunkStream struct {
	message message.IMessage
	process int

	messageLength  uint32
	rMessageLength uint32

	keyIndex byte
	chunkId  byte

	rdataLenTemp uint32
	dataLenTemp  uint32
	dataTemp     []byte
	keyTemp      byte
}

func newChunkStream(message message.IMessage) *chunkStream {
	return &chunkStream{message: message, process: 0, keyIndex: 0, messageLength: message.GetMessageLength(), rMessageLength: 0}
}
func (stream *chunkStream) hasNext() bool {

	return stream.rMessageLength < stream.messageLength
}
func (stream *chunkStream) readChunk() IChunk {

	if stream.process == 0 || stream.process == 1 {
		stream.keyTemp = stream.message.GetKeys()[stream.keyIndex]
		stream.dataTemp = stream.message.GetValue(stream.keyTemp)
		stream.dataLenTemp = uint32(len(stream.dataTemp))
		stream.rdataLenTemp = 0
		stream.keyIndex++
	}

	start := stream.rdataLenTemp
	end := start + 512
	if end > stream.dataLenTemp-start {
		end = stream.dataLenTemp
	}
	stream.rdataLenTemp = end
	stream.rMessageLength = stream.rMessageLength + end - start
	if stream.process == 0 {
		chunk := createChunk0(stream.chunkId, stream.message.GetClassId(), stream.message.GetMessageType(), stream.message.GetMessageLength(), stream.message.GetMessageId(), stream.message.GetTimestamp(), stream.keyTemp, stream.dataLenTemp, stream.dataTemp[start:end])
		if stream.rdataLenTemp < stream.dataLenTemp {
			stream.process = 2
		} else {
			stream.process = 1
		}
		return chunk
	} else if stream.process == 1 {
		chunk := createChunk1(stream.message.GetMessageType(), stream.keyTemp, stream.dataLenTemp, stream.dataTemp[start:end])
		if stream.rdataLenTemp < stream.dataLenTemp {
			stream.process = 2
		}
		stream.keyIndex++
		return chunk
	} else if stream.process == 2 {
		chunk := createChunk2(stream.message.GetMessageType(), stream.dataTemp[start:end])
		if stream.rdataLenTemp == stream.dataLenTemp {
			stream.process = 1
		}
		return chunk
	}
	return nil
}
