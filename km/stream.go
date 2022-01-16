package km

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/util"
	net2 "net"
)

var code []byte = []byte{'c', 'o', 'k', 'e'}

var version []byte = []byte{0, 0, 0, 1}

type Stream struct {
	io        *net.IONetStream
	chunkSize int
	km        km
	isClient  bool
}

func NewStream(io *net.IONetStream) (*Stream, error) {
	st := &Stream{io: io, chunkSize: 512, isClient: false}
	return st, st.init()
}
func NewClientStream(io *net.IONetStream) (*Stream, error) {
	st := &Stream{io: io, chunkSize: 512, isClient: true}
	return st, st.init()
}
func (stream *Stream) init() error {

	if stream.isClient {
		return stream.shakeHandsClient()
	} else {
		return stream.shakeHandsServer()
	}
}

func (stream *Stream) shakeHandsClient() (err error) {

	//写标识位
	num, err := stream.io.Write(code)
	log.InfoF("shakeHandsClient 写入num:{}  {}", num, err)
	if err == nil {
		_, err = stream.io.Write(version)
		if err == nil {
			err = stream.io.Flush()
			if err == nil {
				ver, err1 := stream.io.ReadBytes(4)
				if err1 != nil {
					return err1
				}
				if util.Equal(ver, version) {
					log.Debug("客户端握手成功")
					stream.km = NewKm00001(stream.io)
					return nil
				}else{
					stream.close(0)
					return core.UnKnownVersion
				}
			}else{
				return err
			}
		}
	}
	return core.UnKnownConn
}

/**
作为服务端
*/
func (stream *Stream) shakeHandsServer() error {
	log.Debug("服务端接收握手信息")
	data, err := stream.io.ReadBytes(4)
	log.Debug("服务端接收握手信息")
	if err == nil {
		if util.Equal(data, code) {
			ver, err1 := stream.io.ReadBytes(4)
			if err1 != nil {
				return err1
			}
			if util.Equal(ver, version) {
				_, err2 := stream.io.Write(version)
				err2 = stream.io.Flush()
				if err2 == nil {
					log.Debug("服务端握手成功")
					stream.km = NewKm00001(stream.io)
					return nil
				} else {
					return err2
				}
			} else {
				stream.close(0)
				return core.UnKnownVersion
			}
		}
	} else {
		stream.close(0)
		return err
	}
	return core.UnKnownConn
}
func (stream *Stream) ReadMessage() (message.IMessage, error) {
	return stream.km.ReadMessage()
}
func (stream *Stream) WriteMessage(msg message.IMessage) error {
	err := stream.km.WriteMessage(msg)
	return err
}

func (stream *Stream) close(code int) {
	stream.io.Close()
}

func (stream *Stream) RemoteAddr() *net2.TCPAddr {
	return stream.io.RemoteAddr().(*net2.TCPAddr)
}

func (stream *Stream) LocalAddr() *net2.TCPAddr {
	return stream.io.LocalAddr().(*net2.TCPAddr)
}
