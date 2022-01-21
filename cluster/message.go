package cluster

import (
	"bytes"
	"errors"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/util"
	"strconv"
	"strings"
)

type queryMachineMessage struct {
	*message.Message
}
func newAddUserMessage(machineId string,username string)*message.Message  {
	bm := message.CreateMessage(message.FunctionMessageClass, message.AddUserType)
	bm.SetString(message.USERNAME,username)
	bm.SetString(message.MaChineId,machineId)
	return bm
}
func newDeleteUserMessage(machineId string,username string)*message.Message  {
	bm := message.CreateMessage(message.FunctionMessageClass, message.DeleteUserType)
	bm.SetString(message.USERNAME,username)
	bm.SetString(message.MaChineId,machineId)
	return bm
}
func newQueryMachineMessage(localPort int,machineId string)*queryMachineMessage  {
	bm := &queryMachineMessage{Message: message.CreateMessage(message.FunctionMessageClass, message.QueryMachineType)}
	bm.SetValue(message.LocalMachineAddress,toBytes(localPort,machineId))
	return bm
}
func newQueryMachineBasic(localPort int,machineId string)*queryMachineMessage  {
	bm := &queryMachineMessage{Message: message.CreateMessage(message.FunctionMessageClass, message.QueryMachineBasicType)}
	bm.SetValue(message.LocalMachineAddress,toBytes(localPort,machineId))
	return bm
}
//获取集群信息
func newQueryMachineInfo() message.IMessage {
	bm := &queryMachineMessage{Message: message.CreateMessage(message.FunctionMessageClass, message.QueryMachineInfoType)}
	return bm
}
func newQuery(queryName string,value ...string) message.IMessage {
	bm := &queryMachineMessage{Message: message.CreateMessage(message.FunctionMessageClass, message.QueryType)}
	bm.SetString(message.QueryName,queryName)
	for i,v:=range value{
		bm.SetString(byte(i),v)
	}
	return bm
}

func backQueryMachine(data []byte,msgId uint32)*queryMachineMessage  {
	bm := &queryMachineMessage{Message: message.CreateBackMessage(message.BackMessageClass, message.BackMessageOKType,msgId)}
	bm.SetValue(message.BackMachineAddress,data)
	return bm
}

func backQueryOk(data []byte,msgId uint32)*queryMachineMessage  {
	bm := &queryMachineMessage{Message: message.CreateBackMessage(message.BackMessageClass, message.BackMessageOKType,msgId)}
	bm.SetValue(message.QueryData,data)
	return bm
}
func backQueryError(msgId uint32)*queryMachineMessage  {
	bm := &queryMachineMessage{Message: message.CreateBackMessage(message.BackMessageClass, message.BackMessageErrorType,msgId)}
	return bm
}

func backQueryInfoMachine(data []byte,msgId uint32)*queryMachineMessage  {
	bm := &queryMachineMessage{Message: message.CreateBackMessage(message.BackMessageClass, message.BackMessageOKType,msgId)}
	bm.SetValue(message.QueryMachineInfo,data)
	return bm
}
func toBytes(localPort int,machineId string)[]byte  {
	var localHost = "0.0.0.0"
	return toBytes2(localHost, localPort,machineId)
}
func toBytes2(host string,port int,machineId string)[]byte  {
	var buff  = util.NewBuff()
	buff.WriteString(host)
	buff.WriteByte(':')
	buff.WriteString(strconv.Itoa(port))
	buff.WriteByte('|')
	buff.WriteString(machineId)
	data:= buff.Bytes()
	util.FreeBuff(buff)
	return data
}
func toBytes3(buff *bytes.Buffer,machine *machine)  {
	buff.WriteString(machine.remoteHost)
	buff.WriteByte(':')
	buff.WriteString(strconv.Itoa(machine.remotePort))
	buff.WriteByte('|')
	buff.WriteString(machine.machineId)

}
func toMachine(machineAddress string)(*machine,error)  {
	m:= strings.Split(machineAddress,"|")
	if len(m)<2{
		return nil, errors.New("machineAddress error")
	}
	address := m[0]
	addresses := strings.Split(address,":")
	host:= addresses[0]
	port,err:=strconv.Atoi(addresses[1])
	if err!=nil{
		return nil, err
	}
	var machine = newMachine(port,host)
	machine.machineId = m[1]
	return machine,nil
}