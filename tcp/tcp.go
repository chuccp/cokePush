package tcp

import (
	"github.com/chuccp/cokePush/config"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/user"
)

type Server struct {
	config    *config.Config
	tcpserver *net.TCPServer
	port      int
	userStore *user.Store
}

func (server *Server) Init(context *core.Context) {
	server.port = server.config.GetIntOrDefault("tcp.server.port", 6464)
	server.tcpserver = net.NewTCPServer(server.port)
	server.userStore = context.UserStore
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
			client, err := NewClient(io)
			if err == nil {
				go client.Start()
			}
		}
	}
}
func (server *Server) Name() string {

	return "TCP"
}
