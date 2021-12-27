package net

import (
	"net"
	"strconv"
)

type XConn struct {
	port   int
	host   string
	addr   *net.TCPAddr
	stream *IONetStream
}

func NewXConn(host string, port int) *XConn {
	addr, _ := net.ResolveTCPAddr("tcp", host+":"+strconv.Itoa(port))
	return &XConn{port: port, host: host, addr: addr}
}
func (x *XConn) Create() (*IONetStream,error) {
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

