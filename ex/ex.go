package ex

import (
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/util"
	"github.com/chuccp/utils/config"
	"github.com/chuccp/utils/log"
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
		err=handle("/ex/sendMsg", server.sendMsg)
		if err==nil{
			go server.store.timeoutCheck()
			go server.store.writeBlank()
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
	util.HttpCross(w)
	server.store.jack(w,re)

}
func (server *Server) sendMsg(w http.ResponseWriter, re *http.Request) {
	server.store.sendMsg(w,re)
}
func (server *Server) Name() string {
	return "ex"
}
func NewServer() *Server {
	return &Server{}
}