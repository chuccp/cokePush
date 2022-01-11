package cluster

import (
	"bytes"
	"sync"
)

type machine struct {
	remotePort int
	remoteHost string
	machineId string
	isLocal bool
}

type machineStore struct {
	machineMap *sync.Map
}

func newMachineStore() *machineStore {
	return &machineStore{machineMap: new(sync.Map)}
}
func (machineStore *machineStore)add(machine *machine)bool  {
	mid:=machine.machineId
	_,fa:=machineStore.machineMap.LoadOrStore(mid,machine)
	return !fa
}
func (machineStore *machineStore)getMachine(machineId string)(*machine,bool)  {
	v,fa:=machineStore.machineMap.Load(machineId)
	if fa{
		return v.(*machine),true
	}
	return nil,fa
}
func (machineStore *machineStore)each(f func(machineId string,machine * machine)bool)  {
	machineStore.machineMap.Range(func(key, value interface{}) bool {
		return f(key.(string),value.(* machine))
	})
}
func (machineStore *machineStore)getMachines(buff *bytes.Buffer)  {
	machineStore.machineMap.Range(func(key, value interface{}) bool {
		m:=value.(* machine)
		toBytes3(buff,m)
		buff.WriteString(";")
		return true
	})
}
func (machineStore *machineStore)eachAddress(f func(remoteHost string,remotePort int))  {
	machineStore.machineMap.Range(func(key, value interface{}) bool {
		mv:=value.(* machine)
		f(mv.remoteHost,mv.remotePort)
		return true
	})
}
func has(m *machine,maine * machine) bool {
	if maine.remoteHost==m.remoteHost && maine.remotePort==m.remotePort{
		return true
	}
	return false
}