package net

import (
	"bufio"
	"bytes"
	"io"
	"net"
)

type IOReadStream struct {
	read_ *bufio.Reader
}

func NewIOReadStream(read io.Reader) *IOReadStream {
	return &IOReadStream{read_: bufio.NewReader(read)}
}

func (stream *IOReadStream) ReadLine() ([]byte, error) {
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
func (stream *IOReadStream) read(len int) ([]byte, error) {
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
func (stream *IOReadStream) readUint(len uint32) ([]byte, error) {
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
func (stream *IOReadStream) ReadUintBytes(len uint32) ([]byte, error) {
	return stream.readUint(len)
}

func (stream *IOReadStream) ReadBytes(len int) ([]byte, error) {
	return stream.read(len)
}
func (stream *IOReadStream) ReadByte() (byte, error) {
	return stream.read_.ReadByte()
}

type IOWriteStream struct {
	write_ *bufio.Writer
}

func NewIOWriteStream(write io.Writer) *IOWriteStream {
	return &IOWriteStream{write_: bufio.NewWriter(write)}
}

func (stream *IOWriteStream) Write(data []byte) (int, error) {
	return stream.write_.Write(data)
}
func (stream *IOWriteStream) Flush() error {
	return stream.write_.Flush()
}

type IONetStream struct {
	*net.TCPConn
	*IOReadStream
	*IOWriteStream
	isManualClose bool
}

func NewIOStream(cnn *net.TCPConn) *IONetStream {
	var sm *IONetStream= &IONetStream{TCPConn: cnn, isManualClose: false}
	sm.IOWriteStream = NewIOWriteStream(cnn)
	sm.IOReadStream = NewIOReadStream(cnn)
	return sm
}
func (stream *IONetStream) GetLocalAddress() *net.TCPAddr {
	if stream.LocalAddr() == nil {
		return nil
	}
	return stream.LocalAddr().(*net.TCPAddr)
}
func (stream *IONetStream) GetRemoteAddress() *net.TCPAddr {
	return stream.RemoteAddr().(*net.TCPAddr)
}

func (stream *IONetStream) ManualClose() {
	stream.isManualClose = true
	stream.Close()
}
func (stream *IONetStream) IsManualClose() bool {
	return stream.isManualClose
}
