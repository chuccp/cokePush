package api

import (
	"github.com/chuccp/cokePush/config"
	"net/http"
	"strconv"
)

type Server struct {
	serveMux *http.ServeMux
	config   *config.Config
	port int
}

func (server *Server) root(w http.ResponseWriter, re *http.Request) {

	w.Write([]byte("test"))
}
func (server *Server) Start() error {
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(server.port),
		Handler: server.serveMux,
	}
	error := srv.ListenAndServe()
	return error
}
func (server *Server) Init() {
	server.port = server.config.GetIntOrDefault("rest.server.port", 8080)
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