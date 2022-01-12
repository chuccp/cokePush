package cluster

import (
	log "github.com/chuccp/coke-log"
	"reflect"
	"testing"
)

type uuu struct {

	Aaaa string

}

func TestStream_ReadMessage(t *testing.T){
	var aa uuu

	var hhh interface{}

	hhh=&aa

	type_:=reflect.TypeOf(hhh)

	switch type_.Kind(){

	case reflect.Ptr:

	case reflect.Struct:
	}

	a_:=reflect.TypeOf(aa)
	log.Info(a_.Kind())

	u:=(reflect.New(type_)).Elem().Interface()

	t.Log(reflect.TypeOf(u)==nil)



}
func TestStream_ReadMessage2(t *testing.T){
	var s ="wwww"
	var s1 interface{}
	s1 = s
	var sss  = []interface{}{s1}
	var ss interface{}
	ss = sss
	vs:=ss.([]interface{})

	t.Log(vs)

}