package api

import (
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/util"
	"github.com/pquerna/ffjson/ffjson"
	"net/http"
	"strconv"
)

type Server struct {
	serveMux *http.ServeMux
	port     int
	context  *core.Context
}

func (server *Server) root(w http.ResponseWriter, re *http.Request) {
	var dm map[string]interface{} = make(map[string]interface{})
	dm[VERSION] = core.VERSION
	data, _ := ffjson.Marshal(dm)
	w.Write(data)
}

func (server *Server) WriteMessage(iMessage message.IMessage) error {

	return nil
}
func (server *Server) sendMessage(w http.ResponseWriter, re *http.Request) {
	username := util.GetUsername(re)
	msg := util.GetMessage(re)
	flag := util.GetChanBool()
	server.context.SendMessage(message.CreateBasicMessage("system", username, msg), func(iMessage message.IMessage, err error, u bool) {
		flag <- u
	})
	fa := <-flag
	util.FreeChanBool(flag)
	if fa {
		w.Write([]byte("success"))
	} else {
		w.Write([]byte("NO user"))
	}

}
func (server *Server) clusterInfo(w http.ResponseWriter, re *http.Request) {
	handle := server.context.GetHandle("machineInfo")
	if handle != nil {
		value := handle()
		data, _ := ffjson.Marshal(value)
		w.Write(data)
	} else {
		w.Write([]byte("machineInfo not found"))
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
	if server.port < 0 {
		server.port = context.GetConfig().GetIntOrDefault("rest.server.port", 8080)
	}
	server.AddRoute("/", server.root)
	server.AddRoute("/sendMessage", server.sendMessage)
	server.AddRoute("/clusterInfo", server.clusterInfo)
	context.RegisterHandle("AddRoute", server.addRoute)
}

func (server *Server) Name() string {
	return "api"
}
func (server *Server) addRoute(value ...interface{}) interface{} {

	pattern, ok := value[0].(string)
	if ok {
		handler, ok := value[1].(func(http.ResponseWriter, *http.Request))
		if ok {
			server.serveMux.HandleFunc(pattern, handler)
			return nil
		}
	}
	return http.ErrAbortHandler
}
func (server *Server) AddRoute(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	server.serveMux.HandleFunc(pattern, handler)
}
func NewServer() *Server {
	return &Server{serveMux: http.NewServeMux(), port: -1}
}
