package ex

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/config"
	"github.com/chuccp/cokePush/core"
	"net/http"
)

type Server struct {
	config *config.Config
	context *core.Context
	store *store
}

func (server *Server) Start() error {
	handle := server.context.GetHandle("AddRoute")
	if handle != nil {
		err:=handle("/ex", server.ex)
		err=handle("/sendMsg", server.sendMsg)
		if err==nil{
			log.InfoF("add ex route success")
		}else{
			log.ErrorF("add ex route fail err:{}",err)
		}
	}else{
		log.ErrorF("please start after api")
	}
	return nil
}
func (server *Server) Init(context *core.Context) {
	server.context = context
	server.store = newStore(context)

}
func (server *Server) ex(w http.ResponseWriter, re *http.Request) {
	server.store.jack(w,re)

}
func (server *Server) sendMsg(w http.ResponseWriter, re *http.Request) {
	server.store.sendMsg(w,re)
}
func (server *Server) Name() string {
	return "ex"
}
func NewServer(config *config.Config) *Server {
	return &Server{config:config}
}