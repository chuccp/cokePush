package util

import (
	"bytes"
	"sync"
)

func Equal(a []byte, b []byte) bool {

	return len(a) == len(b) && bytes.Equal(a, b)
}

var poolBuff = sync.Pool{New: func() interface{} {
	return new(bytes.Buffer)
}}

func NewBuff() *bytes.Buffer {
	buff := poolBuff.Get().(*bytes.Buffer)
	return buff
}
func BuffToBytes(buff *bytes.Buffer)[]byte {
	data:=buff.Bytes()
	len:=buff.Len()
	var values  = make([]byte,len)
	copy(values,data)
	return values
}
func FreeBuff(buff *bytes.Buffer) {
	buff.Reset()
	poolBuff.Put(buff)
}
