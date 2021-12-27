package cluster

import (
	"github.com/chuccp/cokePush/config"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/user"
)

type Server struct {
	config    *config.Config
	port      int
	machineId string
	tcpserver *net.TCPServer
	userStore *user.Store

	remotePort int
	remoteHost string
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

func (server *Server) Init(context *core.Context) {
	server.port = server.config.GetIntOrDefault("cluster.local.port", 6361)
	server.machineId = server.config.GetString("cluster.local.machineId")
	if server.machineId == "" {
		server.machineId = MachineId()
	}
	server.tcpserver = net.NewTCPServer(server.port)
	server.userStore = context.UserStore
}
func (server *Server) Name() string {
	return "cluster"
}

func NewServer(config *config.Config) *Server {
	return &Server{config: config}
}

type machine struct {


}