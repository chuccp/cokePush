package api

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/util"
	"github.com/pquerna/ffjson/ffjson"
	"net/http"
	"strconv"
	"sync"
)

type Server struct {
	serveMux *http.ServeMux
	port     int
	context  *core.Context
	query *Query
}

func (server *Server) root(w http.ResponseWriter, re *http.Request) {
	var dm map[string]interface{} = make(map[string]interface{})
	dm[VERSION] = core.VERSION
	data, _ := ffjson.Marshal(dm)
	w.Write(data)
}
func (server *Server) sendMessage(w http.ResponseWriter, re *http.Request) {
	username := util.GetUsername(re)
	msg := util.GetMessage(re)
	flag := util.GetChanBool()
	var once sync.Once
	log.InfoF("发送消息：{}",username)
	server.context.SendMessage(message.CreateBasicMessage("system", username, msg), func(err error, u bool) {
		if u {
			w.Write([]byte("success"))
		} else {
			if err!=nil{
				w.Write([]byte(err.Error()))
			}else{
				w.Write([]byte("NO user"))
			}
		}
		once.Do(func() {
			flag <- u
		})
	})
	 <-flag
	util.FreeChanBool(flag)
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
func (server *Server) queryUser(w http.ResponseWriter, re *http.Request){

	username:=util.GetUsername(re)
	value:=server.context.Query("queryUser",username)
	if value != nil {
		data, _ := ffjson.Marshal(value)
		w.Write(data)
	} else {
		w.Write([]byte("queryUser error"))
	}

}
func (server *Server) systemInfo(w http.ResponseWriter, re *http.Request){
	value:=server.context.Query("systemInfo")
	if value != nil {
		data, _ := ffjson.Marshal(value)
		w.Write(data)
	} else {
		w.Write([]byte("queryUser error"))
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
	server.query = newQuery(context)
	server.port = context.GetConfig().GetIntOrDefault("rest.server.port", 8080)
	server.AddRoute("/", server.root)
	server.AddRoute("/sendMessage", server.sendMessage)
	server.AddRoute("/clusterInfo", server.clusterInfo)
	server.AddRoute("/queryUser", server.queryUser)
	server.AddRoute("/systemInfo", server.systemInfo)
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
