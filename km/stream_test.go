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
/**
测试message 转 chunk
 */
func TestStream_WriteMessage(t *testing.T) {

	bm := message.CreateBasicMessage("333", "4444", "55555")

	fs, err := util.NewFileStream("D:\\attach\\bb.bin")
	if err == nil {
		k := NewKm00001(fs)
		k.WriteMessage(bm)
	}
}
