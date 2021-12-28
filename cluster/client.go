package cluster

import (
	"bytes"
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/km"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
)

type Client struct {
	stream *km.Stream
	server *Server
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

func (client *Client) handleMessage(msg message.IMessage) {
	switch msg.GetClassId() {

	case message.FunctionMessageClass:
		messageType := msg.GetMessageType()
		if messageType == message.QueryMachineInfoType {
			rAddress := msg.GetString(message.LocalMachineAddress)
			log.DebugF("收到数据 rAddress:{}  msgId:{}",rAddress,msg.GetMessageId())
			m, err := toMachine(rAddress)
			if err == nil {
				addr := client.stream.RemoteAddr()
				m.remoteHost = addr.IP.String()
				if m.machineId != client.server.machineId {
					if client.server.machineMap.add(m) {
						log.InfoF("QueryMachineInfoType 添加新的机器连接 :{}|{}", addr.String(), m.machineId)
					}
				}
			}
			data := toBytes(client.server.port, client.server.machineId)
			msg := backQueryMachine(data)
			client.stream.WriteMessage(msg)
		} else if messageType == message.QueryMachineType {
			rAddress := msg.GetString(message.LocalMachineAddress)
			log.DebugF("收到数据 rAddress:{}",rAddress)
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
			data:=buff.Bytes()
			log.DebugF("发送数据：{}",string(data))
			msg := backQueryMachine(data)
			client.stream.WriteMessage(msg)

		}
	case message.LiveMessageClass:
		lm:=message.CreateLiveMessage()
		client.stream.WriteMessage(lm)
	}

}
func NewClient(stream *net.IONetStream, server *Server) (*Client, error) {
	kmStream, err := km.NewStream(stream)
	return &Client{stream: kmStream, server: server}, err
}
