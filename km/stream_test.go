package km

import (
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/util"
	"testing"
	"time"
)

/**
测试message 转 chunk
*/
func TestStream_WriteMessage(t *testing.T) {

	fs, err := util.NewFileStream("D:\\attach\\bb.bin")
	time1 := time.Now()
	var i = 0
	if err == nil {
		for ; i < 1; i++ {
			fs.Seek(0)
			k := NewKm00001(fs)
			msg, err := k.ReadMessage()
			if err == nil {
				msg.GetTimestamp()
			} else {
				t.Log(err)
			}
		}
	}

	time2 := time.Now()

	t.Log(time2.Sub(time1))

}
func TestStream_ReadMessage(t *testing.T) {

	bm := message.CreateBasicMessage("333333", "2222222", "444444")

	fs, err := util.NewFileStream("D:\\attach\\bb.bin")
	time1 := time.Now()
	if err == nil {
		k := NewKm00001(fs)
		k.WriteMessage(bm)
	}
	time2 := time.Now()

	t.Log(time2.Sub(time1))
}
func Test_chan(t *testing.T)  {
	bm := message.CreateBasicMessage("333333", "2222222", "444444")
	m:=newMessageQ(bm)
	go func() {
		time.Sleep(time.Second*2)
		m.fa<-true
	}()

	t.Log(<-m.fa)



}