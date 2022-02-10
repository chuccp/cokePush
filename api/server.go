package api

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
	"github.com/chuccp/cokePush/util"
	"github.com/pquerna/ffjson/ffjson"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Server struct {
	serveMux *http.ServeMux
	port     int
	context  *core.Context
	query    *Query
}

func (server *Server) root(w http.ResponseWriter, re *http.Request) {
	var dm map[string]interface{} = make(map[string]interface{})
	dm[VERSION] = core.VERSION
	data, _ := ffjson.Marshal(dm)
	w.Write(data)
}
func (server *Server) sendMsg(w http.ResponseWriter, re *http.Request) {
	username := util.GetUsername(re)
	msg := util.GetMessage(re)
	flag := util.GetChanBool()
	var once sync.Once
	server.context.SendMessage(message.CreateBasicMessage("system", username, msg), func(err error, u bool) {
		if u {
			w.Write([]byte("success"))
		} else {
			if err != nil {
				w.Write([]byte(err.Error()))
			} else {
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
func (server *Server) sendMessage(w http.ResponseWriter, re *http.Request) {
	util.HttpCrossChunked(w)
	username := util.GetUsername(re)
	msg := util.GetMessage(re)
	if len(username) == 0 || len(msg) == 0 {
		w.WriteHeader(401)
		w.Write([]byte("username or msg can't blank"))
	} else {
		us := strings.Split(username, ",")
		w.Write([]byte("{"))
		var isStart = true
		server.context.SendMultiMessage("system", us, msg, func(username string, status int) {
			if isStart {
				w.Write([]byte("\"" + username + "\":" + strconv.Itoa(status)))
				isStart = false
			} else {
				w.Write([]byte(",\"" + username + "\":" + strconv.Itoa(status)))
			}
		})
		w.Write([]byte("}"))

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
func (server *Server) queryUser(w http.ResponseWriter, re *http.Request) {

	username := util.GetUsername(re)
	value := server.context.Query("queryUser", username)
	if value != nil {
		data, _ := ffjson.Marshal(value)
		w.Write(data)
	} else {
		w.Write([]byte("queryUser error"))
	}

}
func (server *Server) systemInfo(w http.ResponseWriter, re *http.Request) {
	value := server.context.Query("systemInfo")
	if value != nil {
		data, _ := ffjson.Marshal(value)
		w.Write(data)
	} else {
		w.Write([]byte("queryUser error"))
	}
}

func (server *Server) onlineUser(w http.ResponseWriter, re *http.Request) {
	start := util.GetStart(re)
	size := util.GetSize(re)

	log.InfoF("onlineUser  start:{} size:{}", start, size)

	page := user.NewPage()

	var cSize = 0
	var cStart = 0
	queryPageUser := server.context.GetHandle("queryPageUser")
	if queryPageUser != nil {
		p := queryPageUser(start, size).(*user.Page)
		page.Num = p.Num + page.Num
		cStart = start - p.Num
		cSize = size - p.Size()
		page.List = append(page.List, p.List...)

	}

	log.InfoF("onlineUser2  cStart:{} cSize:{}", cStart, cSize)

	if cStart < 0 {
		cStart = 0
	}
	clusterQueryPageUser := server.context.GetHandle("clusterQueryPageUser")
	if clusterQueryPageUser != nil {
		p := clusterQueryPageUser(cStart, cSize).(*user.Page)
		page.Num = p.Num + page.Num
		page.List = append(page.List, p.List...)
	}
	data, _ := ffjson.Marshal(page)
	w.Write(data)

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
	port := os.Getenv("PORT")
	if port != "" {
		p, err := strconv.Atoi(port)
		if err == nil {
			server.port = p
		}else{
			server.port = context.GetConfig().GetIntOrDefault("rest.server.port", 8080)
		}
	} else {
		server.port = context.GetConfig().GetIntOrDefault("rest.server.port", 8080)
	}
	server.AddRoute("/root_version", server.root)
	server.AddRoute("/sendmsg", server.sendMsg)
	server.AddRoute("/sendMessage", server.sendMessage)
	server.AddRoute("/info_user", server.clusterInfo)
	server.AddRoute("/queryUser", server.queryUser)
	server.AddRoute("/systemInfo", server.systemInfo)
	server.AddRoute("/onlineUser", server.onlineUser)
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
