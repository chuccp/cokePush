package tcp

import (
	"github.com/chuccp/cokePush/config"
	"github.com/chuccp/cokePush/net"
)

type Server struct {
	config   *config.Config
	tcpserver *net.TCPServer
	port int
}
func (server *Server) Init() {
	server.port = server.config.GetIntOrDefault("tcp.server.port", 6464)
	server.tcpserver =net.NewTCPServer(server.port)
}
func (server *Server) Start() error{
	err:=server.tcpserver.Bind()
	if err!=nil{
		return err
	}else{
		go server.AcceptConn()
	}
	return nil
}
func (server *Server) AcceptConn(){
	for{
		io,err:=server.tcpserver.Accept()
		if err!=nil{
			break
		}else{
			client, err :=NewClient(io)
			if err ==nil{
				go client.Start()
			}
		}
	}
}