package km

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/util"
)

var code []byte = []byte{'c', 'o', 'k', 'e'}

type Stream struct {
	io        *net.IONetStream
	chunkSize int
	km km
}

func NewStream(io *net.IONetStream) (*Stream, error) {
	st := &Stream{io: io, chunkSize: 512}
	st.init()
	return st, nil
}
func (stream *Stream) init() {

	stream.verify()

}
func (stream *Stream) verify() error {

	data, err := stream.io.ReadBytes(4)
	if err == nil {
		if util.Equal(data, code) {

			ver, err := stream.io.ReadBytes(4)
			if err == nil {
				if util.Equal(ver, []byte{0, 0, 0, 1}) {
					stream.km = NewKm00001(stream.io)
				} else {
					stream.close(0)
				}
			} else {
				return err
			}

		} else {
			stream.close(0)
		}
	} else {
		return err
	}

	return err
}
func (stream *Stream) ReadMessage() (message.IMessage,error) {
	 return stream.km.ReadMessage()
}
func (stream *Stream) WriteMessage(msg message.IMessage)error {
	return stream.km.WriteMessage(msg)
}

func (stream *Stream) close(code int) {

}
