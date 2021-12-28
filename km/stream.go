package km

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/util"
	net2 "net"
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
	log.Debug("客户端发送接收握手信息")
	//写标识位
	num,err:=stream.io.Write(code)
	log.TraceF("shakeHandsClient 写入num:{}",num)
	if err==nil{
		_,err=stream.io.Write(version)
		if err==nil{
			err=stream.io.Flush()
			if err==nil{
				stream.km = NewKm00001(stream.io)
				return nil
			}
			return nil
		}
	}
	return err
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
			ver, err := stream.io.ReadBytes(4)
			if err == nil {
				if util.Equal(ver, version) {
					log.Debug("服务端握手成功")
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
	err:= stream.km.WriteMessage(msg)
	if err==nil{
		err=stream.io.Flush()
	}
	return err
}

func (stream *Stream) close(code int) {
	stream.io.Close()
}

func (stream *Stream) RemoteAddr()*net2.TCPAddr {
	return stream.io.RemoteAddr().(*net2.TCPAddr)
}

func (stream *Stream) LocalAddr() *net2.TCPAddr {
	return  stream.io.LocalAddr().(*net2.TCPAddr)
}