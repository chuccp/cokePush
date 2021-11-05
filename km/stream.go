package km

import (
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/util"
)

var code []byte = []byte{'c', 'o', 'k', 'e'}

type Stream struct {
	io        *net.IOStream
	chunkSize int
	km km
}

func NewStream(io *net.IOStream) (*Stream, error) {
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
					stream.km = NewKm00001()
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
func (stream *Stream) ReadMessage() (Message,error) {
	 return stream.km.ReadMessage()

}
func (stream *Stream) close(code int) {

}
