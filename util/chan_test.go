package util

import (
	log "github.com/chuccp/coke-log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

type cc struct {
	i chan int32
	num int32
}

func (c *cc)wait()int32  {
	 <-c.i
	 return c.num
}

func Test_chan(t *testing.T){


	var ch = make( chan int32)
	for i:=0;i<100;i++{
		time.Sleep(time.Microsecond)
		go func() {
			log.Info("=input==",i)
			c:=&cc{ch, int32(i)}
			log.Info("=output==",c.wait())
		}()
	}
	time.Sleep(time.Second*15)
	for i:=0;i<100;i++{
		time.Sleep(time.Second)
		ch<-1
	}


	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig

}
