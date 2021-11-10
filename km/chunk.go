package km

import "github.com/chuccp/cokePush/message"

type chunk struct {



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
func (stream *chunkStream) readChunk() *chunk {

	if stream.process == 0 {

	}
	return nil
}
