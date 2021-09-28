package api

import (
	"github.com/chuccp/cokePush/config"
	"net/http"
)

type Server struct {
	serveMux *http.ServeMux

}

func (server *Server)root(w http.ResponseWriter, re *http.Request)  {

}
func (server *Server) Start(config *config.Config) {

}
func (server *Server) AddRoute(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	server.serveMux.HandleFunc(pattern,handler)
}


