package local

import (
	"github.com/chuccp/cokePush/core"
)

type Server struct {
	context *core.Context
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) Start() error {
	return nil
}
func (server *Server) Init(context *core.Context) {
	server.context = context
}

func (server *Server) CreateClient(user *User) *Client {
	return newClient(user, server.context)
}

func (server *Server) Name() string {
	return "local"
}
