package api

import (
	"github.com/chuccp/cokePush/config"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/util"
	"github.com/pquerna/ffjson/ffjson"
	"net/http"
	"strconv"
)

type Server struct {
	serveMux *http.ServeMux
	config   *config.Config
	port     int
	context *core.Context
}

func (server *Server) root(w http.ResponseWriter, re *http.Request) {
	var dm map[string]interface{} = make(map[string]interface{})
	dm[VERSION] = core.VERSION
	data, _ := ffjson.Marshal(dm)
	w.Write(data)
}
func (server *Server) sendMessage(w http.ResponseWriter, re *http.Request){
	username:=util.GetUsername(re)
	msg:=util.GetMessage(re)
	err:=server.context.SendMessage(message.CreateBasicMessage("system",username,msg))
	if err==nil{
		w.Write([]byte("success"))
	}else{
		w.Write([]byte(err.Error()))
	}
}


func (server *Server) Start() error {
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(server.port),
		Handler: server.serveMux,
	}
	error := srv.ListenAndServe()
	return error
}
func (server *Server) Init(context *core.Context) {
	server.context = context
	server.port = server.config.GetIntOrDefault("rest.server.port", 8080)
	server.AddRoute("/", server.root)
	server.AddRoute("/sendMessage", server.sendMessage)
}

func (server *Server) Name() string {
	return "api"
}

func (server *Server) AddRoute(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	server.serveMux.HandleFunc(pattern, handler)
}
func NewServer(config *config.Config) *Server {
	return &Server{serveMux: http.NewServeMux(), config: config}
}
