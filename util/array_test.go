package util

import (
	"testing"
)

func TestEqual(t *testing.T) {

	var num uint32 = 1024
	data := U32TOBytes(num)
	println(data[0])

}
