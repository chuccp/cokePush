package km

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/util"
	"io"
	"sync"
)

var chunk0Pool *sync.Pool
var chunk1Pool *sync.Pool
var chunk2Pool *sync.Pool
var chunkWriteStreamPool *sync.Pool
var chunkReadStreamPool *sync.Pool

func init() {
	chunk0Pool = &sync.Pool{New: func() interface{} {
		return new(chunk0)
	}}
	chunk1Pool = &sync.Pool{New: func() interface{} {
		return new(chunk1)
	}}
	chunk2Pool = &sync.Pool{New: func() interface{} {
		return new(chunk2)
	}}

	chunkWriteStreamPool = &sync.Pool{New: func() interface{} {
		return new(chunkWriteStream)
	}}
	chunkReadStreamPool = &sync.Pool{New: func() interface{} {
		return new(chunkReadStream)
	}}

}
func createChunkWriteStreamPool(msg message.IMessage) *chunkWriteStream {
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
func createChunkReadStream(read_ *net.IOReadStream,chunkMap map[uint16]*chunkRecord) *chunkReadStream {
	var chunkReadStream = chunkReadStreamPool.Get().(*chunkReadStream)
	chunkReadStream.read_ =read_
	chunkReadStream.maxBodySize = 512
	chunkReadStream.recordMap = chunkMap
	return chunkReadStream
}
func freeChunkWriteStream(crm *chunkWriteStream) {
	chunkWriteStreamPool.Put(crm)
}
func freeChunkReadStream(crm *chunkReadStream) {
	chunkReadStreamPool.Put(crm)
}

func createPoolChunk0(chunk *chunk) *chunk0 {
	var chk = (chunk0Pool.Get()).(*chunk0)
	chk.chunk = chunk
	return chk
}
func freePoolChunk0(chunk *chunk0) {
	chunk0Pool.Put(chunk)
}

func createPoolChunk1(chunk *chunk) *chunk1 {
	var chk = (chunk1Pool.Get()).(*chunk1)
	chk.chunk = chunk
	return chk
}
func freePoolChunk1(chunk *chunk1) {
	chunk1Pool.Put(chunk)
}

func createPoolChunk2(chunk *chunk) *chunk2 {
	var chk = (chunk2Pool.Get()).(*chunk2)
	chk.chunk = chunk
	return chk
}
func freePoolChunk2(chunk *chunk2) {
	chunk2Pool.Put(chunk)
}

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
	if chunk0.messageLength>0{
		bytesArray = append(bytesArray, chunk0.key)
		bytesArray = append(bytesArray, lengthToBytes(chunk0.dataLen)...)
		bytesArray = append(bytesArray, chunk0.data...)
	}
	return bytesArray
}

func lengthToBytes(length uint32) []byte {
	if length <= 32_767 {
		b := []byte{0, 0}
		b[1] = byte(length)
		b[0] = byte(length >> 8)
		return b
	} else {
		b := []byte{0, 0, 0, 0}
		b[3] = byte(length)
		b[2] = byte(length >> 4)
		b[1] = byte(length >> 8)
		b[0] = byte(length>>12) | 128
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

	maxBodySize uint32
}

var countChunkId uint16 = 0
var lock = new(sync.Mutex)

func getChunkId() uint16 {
	lock.Lock()
	defer lock.Unlock()
	countChunkId = countChunkId + 4
	return countChunkId >> 2
}

func newChunkWriteStream(message message.IMessage) *chunkWriteStream {
	return &chunkWriteStream{maxBodySize: 512, chunkId: getChunkId(), message: message, process: 0, keyIndex: 0, messageLength: message.GetMessageLength(), rMessageLength: 0}
}
func (stream *chunkWriteStream) hasNext() bool {

	return stream.rMessageLength < stream.messageLength
}
func (stream *chunkWriteStream) readChunk0() IChunk {
	if stream.messageLength>0{
		stream.process = 0
		return stream.readChunk()
	}else{
		chunk := createChunk0(stream.chunkId, stream.message.GetClassId(), stream.message.GetMessageType(), stream.message.GetMessageLength(), stream.message.GetMessageId(), stream.message.GetTimestamp(), stream.keyTemp, stream.dataLenTemp, stream.dataTemp)
		return chunk
	}
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
	end := start + stream.maxBodySize
	if end > stream.dataLenTemp {
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
	msg            *message.Message
}

func newChunkRecord(msg *message.Message) *chunkRecord {
	return &chunkRecord{messageLength: 0, rMessageLength: 0, dataLength: 0, rDataLength: 0, msg: msg, data: make([]byte, 0)}
}

type chunkReadStream struct {
	read_       *net.IOReadStream
	recordMap   map[uint16]*chunkRecord
	maxBodySize uint32
}

func newChunkReadStream(read_ io.Reader) *chunkReadStream {
	return &chunkReadStream{read_: net.NewIOReadStream(read_), recordMap: make(map[uint16]*chunkRecord), maxBodySize: 512}
}
func (stream *chunkReadStream) InitChunkRecord(chunkId uint16, chunk *chunk0) {
	msg := message.NewMessage()
	msg.SetClassId(chunk.classId)
	msg.SetMessageType(chunk.messageType)
	msg.SetTimestamp(chunk.time)
	msg.SetMessageId(chunk.messageId)
	stream.recordMap[chunkId] = newChunkRecord(msg)
	stream.recordMap[chunkId].messageLength = chunk.messageLength
}
func (stream *chunkReadStream) SetDataLength(chunkId uint16, length uint32) {
	stream.recordMap[chunkId].dataLength = length
	stream.recordMap[chunkId].rDataLength = 0
}
func (stream *chunkReadStream) GetRDataLength(chunkId uint16) uint32 {
	chunkRecord := stream.recordMap[chunkId]
	return chunkRecord.dataLength - chunkRecord.rDataLength
}
func (stream *chunkReadStream) putChunk0(chunkId uint16, chunk *chunk0) bool {
	chunkRecord := stream.recordMap[chunkId]
	chunkRecord.rDataLength = chunkRecord.rDataLength + uint32(len(chunk.data))
	if chunkRecord.dataLength == chunkRecord.rDataLength {
		chunkRecord.data = make([]byte, 0)
		chunkRecord.msg.SetValue(chunk.key, chunk.data)
		chunkRecord.rMessageLength = chunkRecord.rMessageLength + chunkRecord.dataLength
		if chunkRecord.rMessageLength == chunkRecord.messageLength {
			return true
		}
	} else {
		chunkRecord.key = chunk.key
		chunkRecord.data = append(chunkRecord.data, chunk.data...)
	}
	return false
}
func (stream *chunkReadStream) putChunk1(chunkId uint16, chunk *chunk1) bool {
	chunkRecord := stream.recordMap[chunkId]
	chunkRecord.rDataLength = chunkRecord.rDataLength + uint32(len(chunk.data))
	if chunkRecord.dataLength == chunkRecord.rDataLength {
		chunkRecord.data = make([]byte, 0)
		chunkRecord.msg.SetValue(chunk.key, chunk.data)
		chunkRecord.rMessageLength = chunkRecord.rMessageLength + chunkRecord.dataLength
		if chunkRecord.rMessageLength == chunkRecord.messageLength {
			return true
		}
	} else {
		chunkRecord.key = chunk.key
		chunkRecord.data = append(chunkRecord.data, chunk.data...)
	}
	return false
}
func (stream *chunkReadStream) putChunk2(chunkId uint16, chunk *chunk2) bool {
	chunkRecord := stream.recordMap[chunkId]
	chunkRecord.rDataLength = chunkRecord.rDataLength + uint32(len(chunk.data))
	chunkRecord.data = append(chunkRecord.data, chunk.data...)
	if chunkRecord.dataLength == chunkRecord.rDataLength {
		chunkRecord.msg.SetValue(chunkRecord.key, chunkRecord.data)
		chunkRecord.rMessageLength = chunkRecord.rMessageLength + chunkRecord.dataLength
		if chunkRecord.rMessageLength == chunkRecord.messageLength {
			return true
		}
	}
	return false
}
func (stream *chunkReadStream) readMessage() (*message.Message, error) {

	for {
		chunkId, fa, err := stream.readChunk()
		if err != nil {
			return nil, err
		} else {
			if fa {
				msg:= stream.recordMap[chunkId].msg
				delete(stream.recordMap,chunkId)
				return msg, err
			}
		}
	}
	return nil, nil
}

func (stream *chunkReadStream) readChunk() (uint16, bool, error) {
	data, err := stream.read_.ReadBytes(2)
	if err == nil {
		ck := createChunkHeader2(data)
		chunkId := ck.chunkId()
		if ck.chunkType() == 0 {
			chunk0 := createPoolChunk0(ck)
			data, err = stream.read_.ReadBytes(10)
			if err == nil {
				chunk0.classId = data[0]
				chunk0.messageType = data[1]
				chunk0.time = util.U32BE(data[2:6])
				chunk0.messageId = util.U32BE(data[6:10])
				chunk0.messageLength, err = stream.readMessageLength()
				stream.InitChunkRecord(chunkId, chunk0)
				if chunk0.messageLength==0{
					freePoolChunk0(chunk0)
					return chunkId, true, nil
				}
				chunk0.key, err = stream.read_.ReadByte()
				if err == nil {
					chunk0.dataLen, err = stream.readMessageLength()
					stream.SetDataLength(chunkId, chunk0.dataLen)
					if err == nil {
						if chunk0.dataLen <= stream.maxBodySize {
							chunk0.data, err = stream.read_.ReadUintBytes(chunk0.dataLen)
						} else {
							chunk0.data, err = stream.read_.ReadUintBytes(stream.maxBodySize)
						}
						fa := stream.putChunk0(chunkId, chunk0)
						freePoolChunk0(chunk0)
						return chunkId, fa, nil
					}
				}
			}
			freePoolChunk0(chunk0)
		} else if ck.chunkType() == 1 {
			chunk1 := createPoolChunk1(ck)
			chunk1.key, err = stream.read_.ReadByte()
			if err == nil {
				chunk1.dataLen, err = stream.readMessageLength()
				stream.SetDataLength(chunkId, chunk1.dataLen)
				if err == nil {
					if chunk1.dataLen <= stream.maxBodySize {
						chunk1.data, err = stream.read_.ReadUintBytes(chunk1.dataLen)
					} else {
						chunk1.data, err = stream.read_.ReadUintBytes(stream.maxBodySize)
					}
					fa := stream.putChunk1(chunkId, chunk1)
					freePoolChunk1(chunk1)
					return chunkId, fa, nil
				}
			}
			freePoolChunk1(chunk1)
		} else if ck.chunkType() == 2 {
			chunk2 := createPoolChunk2(ck)
			dataLen := stream.GetRDataLength(chunk2.chunkId())
			if dataLen > stream.maxBodySize {
				chunk2.data, err = stream.read_.ReadUintBytes(stream.maxBodySize)
				fa := stream.putChunk2(chunkId, chunk2)
				freePoolChunk2(chunk2)
				return chunkId, fa, nil
			} else {
				chunk2.data, err = stream.read_.ReadUintBytes(dataLen)
				fa := stream.putChunk2(chunkId, chunk2)
				freePoolChunk2(chunk2)
				return chunkId, fa, nil
			}
			freePoolChunk2(chunk2)
		}
	}
	return 0, false, err
}
func (stream *chunkReadStream) readMessageLength() (uint32, error) {
	var num uint32
	data, err := stream.read_.ReadUintBytes(2)
	if err == nil {
		pre := data[0] & 128
		if pre == 0 {
			num = uint32(data[0])
			num = (num << 8) | uint32(data[1])
			return num, err
		} else {
			data2, err := stream.read_.ReadUintBytes(2)
			if err == nil {
				num =  uint32(data[0])&127
				num = (num << 8) | uint32(data[1])
				num = (num << 8) | uint32(data2[0])
				num = (num << 8) | uint32(data2[1])
			}
			return num, err
		}
	}

	return num, err
}
