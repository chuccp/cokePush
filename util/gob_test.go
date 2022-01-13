package util

import (
	"bytes"
	"encoding/gob"
	"testing"
)

type UU struct {
	AAA string
	BBB string
	CCC string
	DDD string
}

func (u *UU)NewValue()interface{}  {
	var nu UU
	return &nu
}


func Test_endoc(t *testing.T)  {
  var v  UU
  v.AAA = "123"
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err:=enc.Encode(&v)
	if err==nil{
		t.Log(string(data.Bytes()))
		var c  = NewPtr(v)
		dec:=gob.NewDecoder(&data)
		err:=dec.Decode(c)
		if err==nil{
			t.Log(c)
		}else{
			t.Log(err)
		}

	}


}
