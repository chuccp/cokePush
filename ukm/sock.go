package ukm

import "net"

func ListenerUkm(network string, address *net.UDPAddr) (*Listener, error) {
	li := Listener{address: address}
	return li.init()
}

type Listener struct {
	address *net.UDPAddr
	conn *net.UDPConn
	ukm *ukm
}

func (l *Listener) init() ( li *Listener,err error) {
	l.conn, err = net.ListenUDP("udp", l.address)
	li = l
	if err==nil{
		l.ukm = newUkm(l.conn,SERVER)
	}
	return
}
func (l *Listener) Accept() (*Conn, error) {
	return l.ukm.accept()
}


