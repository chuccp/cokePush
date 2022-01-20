package net

import (
	log "github.com/chuccp/coke-log"
	"net"
	"strconv"
)

type XConn struct {
	port   int
	host   string
	address string
	addr   *net.TCPAddr
	stream *IONetStream
}

func NewXConn(host string, port int) *XConn {
	addr:= host+":"+strconv.Itoa(port)
	return NewXConn2(addr)
}
func NewXConn2(address string) *XConn {
	addr, _ := net.ResolveTCPAddr("tcp", address)
	return &XConn{port: addr.Port, host: addr.Network(), addr: addr}
}
func (x *XConn) Create() (*IONetStream,error) {
	log.InfoF("创建连接 {}",x.addr.String())
	conn, err := net.DialTCP("tcp", nil, x.addr)
	if err != nil {
		return nil,err
	}
	x.stream = NewIOStream(conn)
	return x.stream,nil
}
func (x *XConn) Close() {
	x.stream.Close()
}
func (x *XConn) LocalAddress() *net.TCPAddr {
	return x.stream.GetLocalAddress()
}
func (x *XConn) RemoteAddress() *net.TCPAddr {
	return x.stream.GetRemoteAddress()
}

