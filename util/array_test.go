package util

import (
	"testing"
)

func TestEqual(t *testing.T) {


	var code1 = []byte("ansdfkjdoiwdwdhwqdhansdfkjoiwdwdhwqdh")
	var code2 = []byte("ansdfkjdoiwdwdhwqdhansdfkjoiwdwdhwqdh")
	t.Log("sys")
	num:=1000_000_000
	var i = 0
	for  i=0;i<num;i++ {
		Equal(code1,code2)
	}

	t.Log("sys end")


}
