package ukm

import (
	"bytes"
	log "github.com/chuccp/coke-log"
	"io"
	"net"
	"sync"
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
				_, err1 = conn.push(data[0:num])
				u.connMap[key] = conn
				return
			}
		}
	}
}

type Conn struct {
	data   *bytes.Buffer
	rAddr  *net.UDPAddr
	write  chan bool
	isWait bool
	wLock  *sync.Mutex
}

func (conn *Conn) push(p []byte) (n int, err error) {
	conn.wLock.Lock()
	n, err = conn.data.Write(p)
	if conn.isWait {
		conn.write <- true
		conn.wLock.Unlock()
		conn.isWait = false
	} else {
		conn.wLock.Unlock()
	}
	return
}
func NewConn(addr *net.UDPAddr) *Conn {
	return &Conn{rAddr: addr, write: make(chan bool), data: new(bytes.Buffer), isWait: false, wLock: new(sync.Mutex)}
}
func (conn *Conn) Read(p []byte) (n int, err error) {
	for {
		conn.wLock.Lock()
		n, err = conn.data.Read(p)
		log.InfoF("=====num:{} err:{} key:{}", n, err)
		if err == io.EOF && n == 0 {
			conn.isWait = true
			conn.wLock.Unlock()
			<-conn.write
		} else {
			conn.wLock.Unlock()
			return
		}

	}
}
func (conn *Conn) Write(p []byte) (n int, err error) {
	return 0, nil
}
func (conn *Conn) Close() error {
	return nil
}
