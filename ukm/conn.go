package ukm

import (
	"io"
	"net"
	"strconv"
)

type KmListener struct {
	conn      *net.UDPConn
	kmConnMap map[string]*KMConn
}

func createKmListener(conn *net.UDPConn) *KmListener {
	return &KmListener{conn: conn, kmConnMap: make(map[string]*KMConn)}
}
func (kmListener *KmListener) getConn(rAddr net.Addr, lAddr net.Addr) (*KMConn, bool) {
	var key = rAddr.String() + lAddr.String()
	conn := kmListener.kmConnMap[key]
	if conn != nil {
		return conn, true
	}
	return conn, false
}
func (kmListener *KmListener) putConn(rAddr net.Addr, lAddr net.Addr, data []byte) *KMConn {
	var key = rAddr.String() + lAddr.String()
	kmConn := createKMConn(rAddr, kmListener.conn.LocalAddr(), data)
	kmListener.kmConnMap[key] = kmConn
	return kmConn
}
func (kmListener *KmListener) AcceptConn() (*KMConn, error) {
	data := make([]byte, 1460)
	for {
		num, addr, err := kmListener.conn.ReadFromUDP(data)
		if err == nil {
			kn, fa := kmListener.getConn(addr, kmListener.conn.LocalAddr())
			if fa {
				kn.buffer.Write(data[:num])
			} else {
				var kmConn = kmListener.putConn(addr, kmListener.conn.LocalAddr(), data[:num])
				return kmConn, nil
			}
		}
	}
	return nil, io.EOF
}

type KMConn struct {
	lAddr  net.Addr
	rAddr  net.Addr
	buffer *buffer
}

func createKMConn(rAddr net.Addr, lAddr net.Addr, data []byte) *KMConn {
	return &KMConn{buffer: newBuffer(data), rAddr: rAddr, lAddr: lAddr}
}
func (kmConn *KMConn) Read(b []byte) (int, error) {
	if b == nil || len(b) == 0 {
		return 0, ENI
	}
	num, err := kmConn.buffer.Read(b)
	return num, err
}
func (kmConn *KMConn) LocalAddr() net.Addr {
	return kmConn.lAddr
}
func (kmConn *KMConn) RemoteAddr() net.Addr {

	return kmConn.rAddr
}
func ListenKm(localPort int) (*KmListener, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(localPort))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err == nil {
		return createKmListener(conn), err
	}
	return nil, err
}
