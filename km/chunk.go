package km

import "github.com/chuccp/cokePush/message"

type chunk struct {

}

type chunkStream struct {
	message message.IMessage
}
func newChunkStream(message message.IMessage)*chunkStream{
	return &chunkStream{message:message}
}
func(chunk *chunkStream)readChunk()*chunk{

	


	return nil
}