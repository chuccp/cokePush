package ukm

import (
	"bytes"
	"github.com/chuccp/queue"
	"io"
	"net"
)

type MODE uint8

type STATUS uint8

const (
	CLIENT MODE = iota
	SERVER
)

type ukm struct {
	mode    MODE
	udpConn *net.UDPConn
	connMap map[string]*Conn
}

func newUkm(udpConn *net.UDPConn, mode MODE) *ukm {
	return &ukm{udpConn: udpConn, mode: mode, connMap: make(map[string]*Conn)}
}

func (u *ukm) accept() (conn *Conn, err1 error) {
	for {
		data := make([]byte, 1024)
		num, rAddr, err := u.udpConn.ReadFromUDP(data)
		if err != nil {
			return nil, err
		} else {
			key := rAddr.String()
			conn = u.connMap[key]
			if conn != nil {
				conn.push(data[0:num])
			} else {
				conn = NewConn(rAddr)
				conn.push(data[0:num])
				u.connMap[key] = conn
				return
			}
		}
	}
}

type Conn struct {
	buffer *bytes.Buffer
	queue  *queue.VQueue
	rAddr  *net.UDPAddr
	write  chan bool
	isWait bool
}

func (conn *Conn) push(p []byte) {
	conn.buffer.Write(p)
	if conn.buffer.Len() >= 1024 {
		data := make([]byte, 1024)
		start := 0
		for {
			n, err := conn.buffer.Read(data[start:])
			if err == nil {
				start = start + n
				if start == 1024 {
					conn.queue.Offer(data)
					break
				}
			}
		}
	}
}
func NewConn(addr *net.UDPAddr) *Conn {
	return &Conn{rAddr: addr, write: make(chan bool), queue: queue.NewVQueue(), isWait: false, buffer: new(bytes.Buffer)}
}
func (conn *Conn) Read(data []byte) (num int, err error) {
	v, _ := conn.queue.Poll()
	if v != nil {
		return
	} else {
		return 0, io.EOF
	}

}
func (conn *Conn) Write(p []byte) (n int, err error) {
	return 0, nil
}
func (conn *Conn) Close() error {
	return nil
}
