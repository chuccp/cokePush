package km

import (
	"fmt"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/util"
	"github.com/chuccp/utils/log"
	"testing"
	"time"
	"unsafe"
)

/**
测试message 转 chunk
*/
func TestStream_ReadMessage(t *testing.T) {

	fs, err := util.NewFileStream("D:\\attach\\bb.bin")
	k := NewKm00001(net.NewIOStream2(fs))
	time1 := time.Now()
	if err == nil {
		for  {
			msg, err := k.ReadMessage()
			if err == nil {
				t.Log(msg)
			} else {
				t.Log(err)
				break
			}
		}
	}
	time2 := time.Now()
	t.Log(time2.Sub(time1))

}
func TestStream_WriteMessage(t *testing.T) {

	bm := message.CreateBasicMessage("1", "2", "3")

	back:=message.CreateBackMessage(message.BackMessageClass,message.BackMessageOKType,message.MsgId())

	live:=message.CreateLiveMessage()

	fs, err := util.NewFileStream("D:\\attach\\bb.bin")
	k := NewKm00001(net.NewIOStream2(fs))
	time1 := time.Now()
	if err == nil {
		k.WriteMessage(back)
		k.WriteMessage(bm)
		k.WriteMessage(back)
		k.WriteMessage(bm)
		k.WriteMessage(live)
	}else{
		log.InfoF("{}",err)
	}
	time2 := time.Now()

	t.Log(time2.Sub(time1))
}
func Test_chan(t *testing.T)  {
	bm := message.CreateBasicMessage("333333", "2222222", "444444")
	m:=getMessageQ(bm)
	fa:=true
	go func() {
		time.Sleep(time.Second*2)
		for fa{
			log.Info("=====1")
			m.fa<-true
		}
	}()


	t.Log(<-m.fa)
	time.Sleep(time.Second)
	fa = false
	t.Log(m.fa)

}
func ccc(h map[int]int,t *testing.T)  {
	go func() {
		time.Sleep(time.Second)
		t.Log("ccc",(h)[1],unsafe.Pointer(&h))
		fmt.Printf("实际参数的地址 %p\n", &h)
	}()

}
func TestChan(t *testing.T)  {
	chunkMap:=make(map[int]int)
	var h = chunkMap
	ccc(h,t)

	chunkMap[1]=1
	t.Log("TestChan",h[1], unsafe.Pointer(&h))
	fmt.Printf("实际参数的地址 %p\n", &h)
	time.Sleep(time.Second*10)
}
func TestLength(t *testing.T) {

	data:=lengthToBytes(32_767)
	for _,k:=range data{
		t.Logf("%b\n",k)
	}

	num :=  uint32(data[0])&127
	num = uint32(data[0])
	num = (num << 8) | uint32(data[1])

	log.Info(num)
}