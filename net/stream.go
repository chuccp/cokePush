package net

import (
	"bufio"
	"bytes"
	"net"
)

type Stream struct {
	*net.TCPConn
	read_         *bufio.Reader
	write_        *bufio.Writer
	isManualClose bool
}

func NewStream(conn *net.TCPConn) *Stream {
	var sm *Stream = &Stream{TCPConn: conn, isManualClose: false}
	sm.read_ = bufio.NewReader(conn)
	sm.write_ = bufio.NewWriter(conn)
	return sm
}
func (stream *Stream) ReadLine() ([]byte, error) {
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
func (stream *Stream) read(len int) ([]byte, error) {
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
func (stream *Stream) readUint(len uint32) ([]byte, error) {
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
func (stream *Stream) ReadUintBytes(len uint32) ([]byte, error) {
	return stream.readUint(len)
}

func (stream *Stream) ReadBytes(len int) ([]byte, error) {
	return stream.read(len)
}
func (stream *Stream) ReadByte() (byte, error) {
	return stream.read_.ReadByte()
}
func (stream *Stream) Write(data []byte) (int, error) {
	num, err := stream.TCPConn.Write(data)
	if err != nil {
		return num, err
	}
	return num, err
}
func (stream *Stream) GetLocalAddress() *net.TCPAddr {
	if stream.LocalAddr()==nil{
		return nil
	}
	return stream.LocalAddr().(*net.TCPAddr)
}
func (stream *Stream) GetRemoteAddress() *net.TCPAddr {
	return stream.RemoteAddr().(*net.TCPAddr)
}

func (stream *Stream) ManualClose() {
	stream.isManualClose = true
	stream.Close()
}
func (stream *Stream) IsManualClose() bool {
	return stream.isManualClose
}
