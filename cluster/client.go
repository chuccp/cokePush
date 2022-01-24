package cluster

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/km"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/util"
	"github.com/pquerna/ffjson/ffjson"
	"strings"
	"time"
)

type Client struct {
	stream  *km.Stream
	server  *Server
	context *core.Context
}

func (client *Client) Start() {
	for {
		msg, err := client.stream.ReadMessage()
		if err == nil {
			client.handleMessage(msg)
		} else {
			log.ErrorF("服务器端连接断开")
			break
		}
	}
}
/** 测试io性能
 */
func (client *Client)test(){
	time.Sleep(time.Second*5)

	var max  = 10000
	ti:=time.Now().UnixNano()
	for i:=0;i<max;i++{
		live:=message.CreateLiveMessage()
		client.stream.WriteMessage(live)
	}
	log.InfoF("时间：平均{}",(time.Now().UnixNano()-ti)/(int64(max)))
}
// queryMachineInfoType 获取当前服务器信息
func (client *Client) queryMachineBasicType(iMsg message.IMessage) {
	rAddress := iMsg.GetString(message.LocalMachineAddress)
	log.InfoF("收到数据 queryMachineBasicType rAddress:{}  msgId:{}", rAddress, iMsg.GetMessageId())
	m, err := toMachine(rAddress)
	if err == nil {
		addr := client.stream.RemoteAddr()
		m.remoteHost = addr.IP.String()
		if m.machineId != client.server.machineId {
			if client.server.machineMap.add(m) {
				log.InfoF("queryMachineBasicType 添加新的机器连接 :{}|{}", addr.String(), m.machineId)
			}
		}
	}
	data := toBytes(client.server.port, client.server.machineId)
	msg := backQueryMachine(data, iMsg.GetMessageId())
	client.stream.WriteMessage(msg)
}
func (client *Client) queryType(iMsg message.IMessage) {
	qname := iMsg.GetString(message.QueryName)
	iv := make([]interface{}, 0)
	keys := iMsg.GetKeys()
	for _, v := range keys[1:] {
		iv = append(iv, iMsg.GetString(v))
	}
	handle := client.context.GetHandle(qname)
	if handle == nil {
		log.ErrorF("can not find handle:{}", qname)
		qm := backQueryError(iMsg.GetMessageId())
		client.stream.WriteMessage(qm)
		return
	}
	v := handle(iv...)
	if v != nil {
		data, err := ffjson.Marshal(v)
		if err == nil {
			qm := backQueryOk(data, iMsg.GetMessageId())
			client.stream.WriteMessage(qm)
		} else {
			qm := backQueryError(iMsg.GetMessageId())
			client.stream.WriteMessage(qm)
		}

	} else {
		qm := backQueryError(iMsg.GetMessageId())
		client.stream.WriteMessage(qm)
	}
}
func (client *Client) queryMachineInfoType(iMsg message.IMessage) {
	mi := client.server.queryMachineInfo()
	data, err := ffjson.Marshal(mi)
	if err == nil {
		msg := backQueryInfoMachine(data, iMsg.GetMessageId())
		client.stream.WriteMessage(msg)
	} else {
		msg := backQueryInfoMachine([]byte(`{"Address":"`+client.server.machineId+`@`+err.Error()+`"}`), iMsg.GetMessageId())
		client.stream.WriteMessage(msg)
	}
}

// QueryMachineType 获取当前服务器集群列表
func (client *Client) QueryMachineType(iMsg message.IMessage) {
	rAddress := iMsg.GetString(message.LocalMachineAddress)
	log.DebugF("收到数据 QueryMachineType rAddress:{}", rAddress)
	m, err := toMachine(rAddress)
	if err == nil {
		addr := client.stream.RemoteAddr()
		m.remoteHost = addr.IP.String()
		if m.machineId != client.server.machineId {
			if client.server.machineMap.add(m) {
				log.InfoF("QueryMachineType 添加新的机器连接 :{}|{}", addr.String(), m.machineId)
			}
		}
	}
	var buff  = util.NewBuff()
	client.server.machineMap.getMachines(buff)
	buff.WriteString(client.stream.LocalAddr().String())
	buff.WriteString("|")
	buff.WriteString(client.server.machineId)
	data :=util.BuffToBytes(buff)
	log.DebugF("发送数据：{}", string(data))
	msg := backQueryMachine(data, iMsg.GetMessageId())
	client.stream.WriteMessage(msg)
}
func (client *Client) handleMessage(msg message.IMessage) {
	log.DebugF("请求来了 class:{}   type:{} msgId:{}", msg.GetClassId(), msg.GetMessageType(), msg.GetMessageId())
	switch msg.GetClassId() {
	case message.FunctionMessageClass:
		log.DebugF("FunctionMessageClass：", msg.GetMessageId())
		messageType := msg.GetMessageType()
		if messageType == message.QueryMachineBasicType {
			client.queryMachineBasicType(msg)
		} else if messageType == message.QueryMachineType {
			client.QueryMachineType(msg)
		} else if messageType == message.QueryMachineInfoType {
			client.queryMachineInfoType(msg)
		} else if messageType == message.AddUserType {
			username := msg.GetString(message.USERNAME)
			machineId := msg.GetString(message.MaChineId)
			client.server.addUser(username, machineId)
		} else if messageType == message.DeleteUserType {
			username := msg.GetString(message.USERNAME)
			machineId := msg.GetString(message.MaChineId)
			client.server.delete(username, machineId)
		} else if messageType == message.QueryType {
			client.queryType(msg)
		}
	case message.LiveMessageClass:
		log.DebugF("LiveMessageClass：", msg.GetMessageId())
		lm := message.CreateLiveMessage()
		client.stream.WriteMessage(lm)
	case message.OrdinaryMessageClass:
		messageType := msg.GetMessageType()
		if messageType == message.BasicMessageType {
			client.context.SendMessageNoForward(msg, func(err error, hasUser bool) {
				nMsg := message.CreateBackBasicMessage(hasUser, msg.GetMessageId())
				err2 := client.stream.WriteMessage(nMsg)
				log.DebugF("收到普通文本信息 msgId:{} 处理信息1:{}   {}", msg.GetMessageId(), err, err2)
			})
		} else if messageType == message.MultiMessageType {
			from := msg.GetString(message.FromUser)
			to := msg.GetString(message.ToUser)
			text := msg.GetString(message.Text)
			ids:=strings.Split(to, ";")
			go client.context.SendMultiMessageNoReplay(from, &ids, text)
		}
	}
}
func NewClient(stream *net.IONetStream, server *Server, context *core.Context) (*Client, error) {
	kmStream, err := km.NewStream(stream)
	if err != nil {
		return nil, err
	}
	return &Client{stream: kmStream, server: server, context: context}, err
}
