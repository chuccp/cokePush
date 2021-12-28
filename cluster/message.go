package cluster

import (
	"bytes"
	"errors"
	"github.com/chuccp/cokePush/message"
	"strconv"
	"strings"
)

type queryMachineMessage struct {
	*message.Message
}



func newQueryMachineMessage(localPort int,machineId string)*queryMachineMessage  {
	bm := &queryMachineMessage{Message: message.CreateMessage(message.FunctionMessageClass, message.QueryMachineType)}
	bm.SetValue(message.LocalMachineAddress,toBytes(localPort,machineId))
	return bm
}
func newQueryMachineInfo(localPort int,machineId string)*queryMachineMessage  {
	bm := &queryMachineMessage{Message: message.CreateMessage(message.FunctionMessageClass, message.QueryMachineInfoType)}
	bm.SetValue(message.LocalMachineAddress,toBytes(localPort,machineId))
	return bm
}

func backQueryMachine(data []byte,msgId uint32)*queryMachineMessage  {
	bm := &queryMachineMessage{Message: message.CreateBackMessage(message.BackMessageClass, message.BackMessageOKType,msgId)}
	bm.SetValue(message.BackMachineAddress,data)
	return bm
}

func toBytes(localPort int,machineId string)[]byte  {
	var localHost = "0.0.0.0"
	return toBytes2(localHost, localPort,machineId)
}
func toBytes2(host string,port int,machineId string)[]byte  {
	var buff bytes.Buffer
	buff.WriteString(host)
	buff.WriteByte(':')
	buff.WriteString(strconv.Itoa(port))
	buff.WriteByte('|')
	buff.WriteString(machineId)
	return buff.Bytes()
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