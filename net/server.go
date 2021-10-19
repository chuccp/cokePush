package net

import "net"

type TCPServer struct {
	listener     *net.TCPListener
	port int
}

func NewTCPServer(port int) *TCPServer {
	return &TCPServer{port: port}
}

func (server *TCPServer) bind() error {
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: server.port,
	})
	if err==nil{
		server.listener = listener
	}
	return err
}

func (server *TCPServer) accept() (*Stream,error) {
	tcn, err := server.listener.AcceptTCP()
	if err!=nil{
		stream := NewStream(tcn)
		return stream,err
	}
	return nil,nil
}
