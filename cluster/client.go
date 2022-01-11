package cluster

import (
	"bytes"
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/km"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
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

// queryMachineInfoType 获取当前服务器信息
func (client *Client) queryMachineBasicType(iMsg message.IMessage) {
	rAddress := iMsg.GetString(message.LocalMachineAddress)
	log.DebugF("收到数据 queryMachineBasicType rAddress:{}  msgId:{}", rAddress, iMsg.GetMessageId())
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
func (client *Client) queryMachineInfoType(iMsg message.IMessage) {

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
	var buff bytes.Buffer
	client.server.machineMap.getMachines(&buff)
	buff.WriteString(client.stream.LocalAddr().String())
	buff.WriteString("|")
	buff.WriteString(client.server.machineId)
	data := buff.Bytes()
	log.DebugF("发送数据：{}", string(data))
	msg := backQueryMachine(data, iMsg.GetMessageId())
	client.stream.WriteMessage(msg)
}
func (client *Client) handleMessage(msg message.IMessage) {
	log.InfoF("请求来了 class:{}   type:{} msgId:{}", msg.GetClassId(), msg.GetMessageType(),msg.GetMessageId())
	switch msg.GetClassId() {
	case message.FunctionMessageClass:
		messageType := msg.GetMessageType()
		if messageType == message.QueryMachineBasicType {
			client.queryMachineBasicType(msg)
		} else if messageType == message.QueryMachineType {
			client.QueryMachineType(msg)
		} else if messageType == message.QueryMachineInfoType {
			client.queryMachineInfoType(msg)
		}
	case message.LiveMessageClass:
		lm := message.CreateLiveMessage()
		client.stream.WriteMessage(lm)
	case message.OrdinaryMessageClass:
		messageType := msg.GetMessageType()
		if messageType==message.BasicMessageType{
			log.InfoF("收到普通文本信息 msgId:{}",msg.GetMessageId())
			client.context.SendMessageNoForward(msg, func(err error, hasUser bool) {

				nMsg:=message.CreateBackBasicMessage(hasUser,msg.GetMessageId())
				err2:=client.stream.WriteMessage(nMsg)
				log.InfoF("收到普通文本信息 msgId:{} 处理信息1:{}   {}",msg.GetMessageId(),err,err2)
			})
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
