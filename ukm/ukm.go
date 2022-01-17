package ukm

import (
	"bytes"
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

func (u *ukm) accept() (*Conn, error) {
	for {
		data := make([]byte, 1024)
		num, rAddr, err := u.udpConn.ReadFromUDP(data)
		if err != nil {
			return nil, err
		} else {
			key := rAddr.String()
			conn := u.connMap[key]
			if conn != nil {
				conn.push(data[0:num])
			} else {
				cn := NewConn(rAddr)
				cn.push(data[0:num])
				u.connMap[key] = cn
			}
		}
	}
}

type Conn struct {
	data   []byte
	buffer *bytes.Buffer
	rAddr  *net.UDPAddr
}

func (conn *Conn) push(p []byte) (n int, err error) {
	return 0, nil
}
func NewConn(addr *net.UDPAddr) *Conn {
	return &Conn{rAddr: addr}
}
func (conn *Conn) Read(p []byte) (n int, err error) {

	return 0, nil
}
func (conn *Conn) Write(p []byte) (n int, err error) {
	return 0, nil
}
func (conn *Conn) Close() error {
	return nil
}
