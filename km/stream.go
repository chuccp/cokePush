package km

import (
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/util"
)

var code []byte = []byte{'c', 'o', 'k', 'e'}

var version []byte=[]byte{0, 0, 0, 1}

type Stream struct {
	io        *net.IONetStream
	chunkSize int
	km        km
	isClient bool
}

func NewStream(io *net.IONetStream) (*Stream, error) {
	st := &Stream{io: io, chunkSize: 512,isClient:false}
	return st, st.init()
}
func NewClientStream(io *net.IONetStream) (*Stream, error) {
	st := &Stream{io: io, chunkSize: 512,isClient:true}
	return st, st.init()
}
func (stream *Stream) init()error {

	if stream.isClient{
		return stream.shakeHandsClient()
	}else{
		return stream.shakeHandsServer()
	}
}

func (stream *Stream) shakeHandsClient() error{
	//写标识位
	_,err:=stream.io.Write(code)
	if err==nil{
		_,err:=stream.io.Write(version)
		if err==nil{
			return nil
		}
	}
	return core.ProtocolError
}
/**
作为服务端
 */
func (stream *Stream) shakeHandsServer() error {

	data, err := stream.io.ReadBytes(4)
	if err == nil {
		if util.Equal(data, code) {
			ver, err := stream.io.ReadBytes(4)
			if err == nil {
				if util.Equal(ver, version) {
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
func (stream *Stream) ReadMessage() (message.IMessage, error) {
	return stream.km.ReadMessage()
}
func (stream *Stream) WriteMessage(msg message.IMessage) error {
	return stream.km.WriteMessage(msg)
}

func (stream *Stream) close(code int) {
	stream.io.Close()
}
