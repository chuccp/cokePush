package ukm

import "net"

func ListenerUkm(network string, address *net.UDPAddr) (*Listener, error) {
	li := Listener{address: address}
	return li.init(network)
}

type Listener struct {
	address *net.UDPAddr
	conn *net.UDPConn
	ukm *ukm
}

func (l *Listener) init(network string) (  li *Listener,err error) {
	l.conn, err = net.ListenUDP(network, l.address)
	li = l
	if err==nil{
		l.ukm = newUkm(l.conn,SERVER)
	}
	return
}
func (l *Listener) Accept() (*Conn, error) {
	return l.ukm.accept()
}


