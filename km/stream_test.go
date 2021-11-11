package km

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/util"
	"testing"
)

func TestStream_ReadMessage(t *testing.T) {

	var num byte = 98
	var i = 0
	for ; i < 10_000_000_000; i++ {
		num = num <<2>>2
	}


}
func TestStream_WriteMessage(t *testing.T) {

	bm := message.CreateBasicMessage("123", "456", "789")

	fs, err := util.NewFileStream("")
	if err == nil {
		k := NewKm00001(fs)
		k.WriteMessage(bm)
	}
}
