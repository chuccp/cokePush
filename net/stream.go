package net

import (
	"bufio"
	"bytes"
	"net"
)
type IOStream struct {
	*net.TCPConn
	read_         *bufio.Reader
	write_        *bufio.Writer
	isManualClose bool
}

func NewIOStream(conn *net.TCPConn) *IOStream {
	var sm = &IOStream{TCPConn: conn, isManualClose: false}
	sm.read_ = bufio.NewReader(conn)
	sm.write_ = bufio.NewWriter(conn)
	return sm
}
func (stream *IOStream) ReadLine() ([]byte, error) {
	buffer := bytes.Buffer{}
	for {
		data, is, err := stream.read_.ReadLine()
		if err != nil {
			return data, err
		}
		if is {
			if len(data) > 0 {
				buffer.Write(data)
			}
		} else {
			buffer.Write(data)
			return buffer.Bytes(), nil
		}
	}
	return nil, nil
}
func (stream *IOStream) read(len int) ([]byte, error) {
	data := make([]byte, len)
	var l = 0
	for l < len {
		n, err := stream.read_.Read(data[l:])
		if err != nil {
			return nil, err
		}
		l += n
	}
	return data, nil
}
func (stream *IOStream) readUint(len uint32) ([]byte, error) {
	data := make([]byte, len)
	var l uint32 = 0
	for l < len {
		n, err := stream.read_.Read(data[l:])
		if err != nil {
			return nil, err
		}
		l += (uint32)(n)
	}
	return data, nil
}
func (stream *IOStream) ReadUintBytes(len uint32) ([]byte, error) {
	return stream.readUint(len)
}

func (stream *IOStream) ReadBytes(len int) ([]byte, error) {
	return stream.read(len)
}
func (stream *IOStream) ReadByte() (byte, error) {
	return stream.read_.ReadByte()
}
func (stream *IOStream) Write(data []byte) (int, error) {
	return stream.TCPConn.Write(data)
}
func (stream *IOStream) GetLocalAddress() *net.TCPAddr {
	if stream.LocalAddr()==nil{
		return nil
	}
	return stream.LocalAddr().(*net.TCPAddr)
}
func (stream *IOStream) GetRemoteAddress() *net.TCPAddr {
	return stream.RemoteAddr().(*net.TCPAddr)
}

func (stream *IOStream) ManualClose() {
	stream.isManualClose = true
	stream.Close()
}
func (stream *IOStream) IsManualClose() bool {
	return stream.isManualClose
}
