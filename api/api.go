package api

import (
	"github.com/chuccp/cokePush/config"
	"net/http"
	"strconv"
)

type Server struct {
	serveMux *http.ServeMux
	config *config.Config

}



func (server *Server)root(w http.ResponseWriter, re *http.Request)  {

}
func (server *Server) Start(config *config.Config)error {
	port:=config.GetIntOrDefault("rest.server.port",8080)
	srv := &http.Server{
		Addr: ":" + strconv.Itoa(port),
		Handler:server.serveMux,
	}
	error:=srv.ListenAndServe()
	return error
}
func (server *Server) Init() {


}
func (server *Server) AddRoute(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	server.serveMux.HandleFunc(pattern,handler)
}


