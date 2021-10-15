package api

import (
	"github.com/chuccp/cokePush/config"
	"net/http"
	"strconv"
)

type Server struct {
	serveMux *http.ServeMux
	config   *config.Config
}

func (server *Server) root(w http.ResponseWriter, re *http.Request) {

	w.Write([]byte("test"))
}
func (server *Server) Start() error {
	port := server.config.GetIntOrDefault("rest.server.port", 8080)
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: server.serveMux,
	}
	error := srv.ListenAndServe()
	return error
}
func (server *Server) Init() {
	server.AddRoute("/",server.root)
}

func (server *Server)Name()string {
	return "api"
}

func (server *Server) AddRoute(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	server.serveMux.HandleFunc(pattern, handler)
}
func NewServer(config   *config.Config) *Server {
	return &Server{serveMux:http.NewServeMux(),config:config}
}