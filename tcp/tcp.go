package tcp

import (
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/net"
)

type Server struct {
	tcpserver *net.TCPServer
	port      int
	context *core.Context
}

func (server *Server) Init(context *core.Context) {
	server.context = context
	server.port = context.GetConfig().GetIntOrDefault("tcp.server.port", 6464)
	server.tcpserver = net.NewTCPServer(server.port)
}
func (server *Server) Start() error {
	err := server.tcpserver.Bind()
	if err != nil {
		return err
	}
	go server.AcceptConn()
	return nil
}

func (server *Server) AcceptConn() {
	for {
		io, err := server.tcpserver.Accept()
		if err != nil {
			break
		} else {
			client, err := NewClient(io,server.context)
			if err == nil {
				go client.Start()
			}
		}
	}
}
func (server *Server) Name() string {

	return "TCP"
}
func NewServer() *Server {
	return &Server{ port: -1}
}