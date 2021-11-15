package km

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/util"
	"io"
	"sync"
)

type chunk0 struct {
	*chunk
	classId       byte
	messageType   byte
	messageLength uint32
	messageId     uint32
	time          uint32
	key           byte
	dataLen       uint32
	data          []byte
}

func createChunk0(chunkId uint16, classId byte, messageType byte, messageLength uint32, messageId uint32, time uint32, key byte, dataLen uint32, data []byte) *chunk0 {
	return &chunk0{chunk: createChunkHeader(0, chunkId), classId: classId, messageType: messageType, messageLength: messageLength, messageId: messageId, time: time, key: key, dataLen: dataLen, data: data}
}
func newChunk0(chunk *chunk) *chunk0 {
	return &chunk0{chunk: chunk}
}

func (chunk0 chunk0) toByte() []byte {
	bytesArray := make([]byte, 0)
	//写chunk
	bytesArray = append(bytesArray, chunk0.chunk.toByte()...)
	//写header
	bytesArray = append(bytesArray, chunk0.classId)
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
	if length < 32_767 {
		b := []byte{0, 0}
		b[1] = byte(length)
		b[0] = byte(length >> 8)
		return b
	} else {
		b := []byte{0, 0, 0, 0}
		b[3] = byte(length)
		b[2] = byte(length >> 4)
		b[1] = byte(length >> 8)
		b[0] = byte(length>>12) | 64
		return b
	}
}

type chunk1 struct {
	*chunk
	key     byte
	dataLen uint32
	data    []byte
}

func createChunk1(chunkId uint16, key byte, dataLen uint32, data []byte) *chunk1 {

	return &chunk1{chunk: createChunkHeader(1, chunkId), dataLen: dataLen, key: key, data: data}
}
func newChunk1(chunk *chunk) *chunk1 {
	return &chunk1{chunk: chunk}
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

func newChunk2(chunk *chunk) *chunk2 {
	return &chunk2{chunk: chunk}
}
func createChunk2(chunkId uint16, data []byte) *chunk2 {
	return &chunk2{chunk: createChunkHeader(2, chunkId), data: data}
}
func (chunk2 chunk2) toByte() []byte {
	bytesArray := make([]byte, 0)
	bytesArray = append(bytesArray, chunk2.chunk.toByte()...)
	bytesArray = append(bytesArray, chunk2.data...)
	return bytesArray
}

type chunk struct {
	chunkHeader uint16
}

func createChunkHeader(chunkType byte, chunkId uint16) *chunk {

	return &chunk{uint16(chunkType)<<14 | chunkId}
}
func createChunkHeader2(data []byte) *chunk {
	return &chunk{chunkHeader: util.BytesTOU16(data)}
}
func (chunk *chunk) chunkType() byte {
	return byte(chunk.chunkHeader >> 14)
}
func (chunk *chunk) chunkId() uint16 {
	return chunk.chunkHeader << 2 >> 2
}
func (chunk *chunk) toByte() []byte {
	return util.U16TOBytes(chunk.chunkHeader)
}

type IChunk interface {
	chunkType() byte
	chunkId() uint16
	toByte() []byte
}

type chunkWriteStream struct {
	message message.IMessage
	process int

	messageLength  uint32
	rMessageLength uint32

	keyIndex byte
	chunkId  uint16

	rdataLenTemp uint32
	dataLenTemp  uint32
	dataTemp     []byte
	keyTemp      byte
}

var chunkId uint16 = 0
var lock = new(sync.Mutex)

func getChunkId() uint16 {
	lock.Lock()
	defer lock.Unlock()
	chunkId = chunkId + 4
	return chunkId >> 2
}

func newChunkWriteStream(message message.IMessage) *chunkWriteStream {
	return &chunkWriteStream{chunkId: getChunkId(), message: message, process: 0, keyIndex: 0, messageLength: message.GetMessageLength(), rMessageLength: 0}
}
func (stream *chunkWriteStream) hasNext() bool {

	return stream.rMessageLength < stream.messageLength
}
func (stream *chunkWriteStream) readChunk() IChunk {

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

		chunk := createChunk1(stream.chunkId, stream.keyTemp, stream.dataLenTemp, stream.dataTemp[start:end])
		if stream.rdataLenTemp < stream.dataLenTemp {
			stream.process = 2
		}
		return chunk
	} else if stream.process == 2 {

		chunk := createChunk2(stream.chunkId, stream.dataTemp[start:end])
		if stream.rdataLenTemp == stream.dataLenTemp {
			stream.process = 1
		}
		return chunk
	}
	return nil
}

type chunkRecord struct {
	messageLength  uint32
	rMessageLength uint32
	key            byte
	data           []byte
	dataLength     uint32
	rDataLength    uint32
	isOver         bool
	msg            *message.Message
}

func newChunkRecord(msg *message.Message) *chunkRecord {
	return &chunkRecord{messageLength: 0, rMessageLength: 0, dataLength: 0, rDataLength: 0, isOver: false, msg: msg,data:make([]byte,0)}
}
type chunkReadStream struct {
	read_     *net.IOReadStream
	recordMap map[uint16]*chunkRecord
}

func newChunkReadStream(read_ io.Reader) *chunkReadStream {
	return &chunkReadStream{read_: net.NewIOReadStream(read_), recordMap: make(map[uint16]*chunkRecord)}
}
func (stream *chunkReadStream) InitChunkRecord(chunkId uint16, classId byte, messageType byte, time uint32) {
	msg := message.NewMessage()
	msg.SetClassId(classId)
	msg.SetMessageType(messageType)
	msg.SetTimestamp(time)
	stream.recordMap[chunkId] = newChunkRecord(msg)
}
func (stream *chunkReadStream) SetMessageLength(chunkId uint16, length uint32) {
	stream.recordMap[chunkId].messageLength = length
}
func (stream *chunkReadStream) SetDataLength(chunkId uint16, length uint32) {
	stream.recordMap[chunkId].dataLength = length
	stream.recordMap[chunkId].rDataLength = 0
}

func (stream *chunkReadStream) CompareDataLength(chunkId uint16, length uint32) {
	chunkRecord := stream.recordMap[chunkId]
	chunkRecord.rDataLength = chunkRecord.rDataLength + length
	chunkRecord.rMessageLength = chunkRecord.rMessageLength + length
}
func (stream *chunkReadStream) GetRDataLength(chunkId uint16) uint32 {
	chunkRecord := stream.recordMap[chunkId]
	return chunkRecord.dataLength - chunkRecord.rDataLength
}
func (stream *chunkReadStream) putChunk0(chunkId uint16, chunk *chunk0) {
	chunkRecord := stream.recordMap[chunkId]
	if chunkRecord.dataLength == chunkRecord.rDataLength {
		chunkRecord.msg.SetValue(chunk.key, chunk.data)
	}else{
		chunkRecord.key = chunk.key
		chunkRecord.data = append(chunkRecord.data,chunk.data...)
	}
}
func (stream *chunkReadStream) putChunk1(chunkId uint16, chunk *chunk1) {
	chunkRecord := stream.recordMap[chunkId]
	if chunkRecord.dataLength == chunkRecord.rDataLength {
		chunkRecord.msg.SetValue(chunk.key, chunk.data)
	}else{
		chunkRecord.key = chunk.key
		chunkRecord.data = append(chunkRecord.data,chunk.data...)
	}
}
func (stream *chunkReadStream) putChunk2(chunkId uint16, chunk *chunk2) {
	chunkRecord := stream.recordMap[chunkId]
	chunkRecord.data = append(chunkRecord.data,chunk.data... )
	if chunkRecord.dataLength == chunkRecord.rDataLength {
		chunkRecord.msg.SetValue(chunkRecord.key, chunkRecord.data)
	}
}
func (stream *chunkReadStream) readChunk() (IChunk, error) {
	data, err := stream.read_.ReadBytes(2)
	if err == nil {
		ck := createChunkHeader2(data)
		if ck.chunkType() == 0 {
			chunk0 := newChunk0(ck)
			data, err = stream.read_.ReadBytes(14)
			chunk0.classId = data[0]
			chunk0.messageType = data[1]
			chunk0.time = util.U32BE(data[2:6])
			chunk0.messageId = util.U32BE(data[6:10])
			chunk0.messageLength, err = stream.readMessageLength()
			stream.InitChunkRecord(chunk0.chunkId(), chunk0.classId, chunk0.messageType, chunk0.time)
			stream.SetMessageLength(chunk0.chunkId(), chunk0.messageLength)
			if err == nil {
				chunk0.key, err = stream.read_.ReadByte()
				if err == nil {
					chunk0.dataLen, err = stream.readMessageLength()
					stream.SetDataLength(chunk0.chunkId(), chunk0.dataLen)
					if err == nil {
						if chunk0.dataLen < 512 {
							chunk0.data, err = stream.read_.ReadUintBytes(chunk0.dataLen)
						} else {
							chunk0.data, err = stream.read_.ReadUintBytes(512)
						}
						stream.CompareDataLength(chunk0.chunkId(), uint32(len(chunk0.data)))
						stream.putChunk0(chunk0.chunkId(), chunk0)
						return chunk0, nil
					}
				}
			}
		} else if ck.chunkType() == 1 {
			chunk1 := newChunk1(ck)
			chunk1.key, err = stream.read_.ReadByte()
			if err == nil {
				chunk1.dataLen, err = stream.readMessageLength()
				stream.SetDataLength(chunk1.chunkId(), chunk1.dataLen)
				if err == nil {
					if chunk1.dataLen < 512 {
						chunk1.data, err = stream.read_.ReadUintBytes(chunk1.dataLen)
					} else {
						chunk1.data, err = stream.read_.ReadUintBytes(512)
					}
					stream.CompareDataLength(chunk1.chunkId(), uint32(len(chunk1.data)))
					stream.putChunk1(chunk1.chunkId(), chunk1)
					return chunk1, nil
				}
			}

		} else if ck.chunkType() == 2 {
			chunk2 := newChunk2(ck)
			dataLen := stream.GetRDataLength(chunk2.chunkId())
			if dataLen>512{
				chunk2.data, err = stream.read_.ReadUintBytes(512)
				stream.putChunk2(chunk2.chunkId(), chunk2)
				return chunk2, nil
			}else{
				chunk2.data, err = stream.read_.ReadUintBytes(dataLen)
				stream.putChunk2(chunk2.chunkId(), chunk2)
				return chunk2, nil
			}


		}
	}
	return nil, err
}
func (stream *chunkReadStream) readMessageLength() (uint32, error) {
	var num uint32
	data, err := stream.read_.ReadUintBytes(2)
	if err == nil {
		pre := data[0] & 64
		if pre == 0 {
			num = num | uint32(data[0])
			num = (num << 8) | uint32(data[1])
			return num, err
		} else {
			data2, err := stream.read_.ReadUintBytes(2)
			if err == nil {
				num = (num | 64) | uint32(data[0])
				num = (num << 8) | uint32(data[1])
				num = (num << 8) | uint32(data2[2])
				num = (num << 8) | uint32(data2[3])
			}
			return num, err
		}
	}

	return num, err
}
