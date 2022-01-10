package ex

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/util"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

func Test_time2(t *testing.T) {

		for i:=0;i<5;i++{

			go func() {
				timer := util.GetTimer(time.Second*2)
				go func() {
					fa:=timer.Wait()
					log.InfoF("============== {}",fa)
					if !fa{
						log.InfoF("22222222")
					}else{
						log.InfoF("33333333")
					}
				}()
				time.Sleep(time.Second*2)
				log.InfoF("121212121212")
				timer.End()
				util.FreeTimer(timer)
			}()

			time.Sleep(time.Second)

			go func() {
				timer := util.GetTimer(time.Second*2)
				go func() {
					fa:=timer.Wait()
					log.InfoF("============== {}",fa)
					if !fa{
						log.InfoF("!!!!!!22222222")
					}else{
						log.InfoF("!!!!!!!33333333")
					}
				}()
				time.Sleep(time.Second*2)
				log.InfoF("121212121212")
				timer.End()
				util.FreeTimer(timer)
			}()
		}





	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig
}

type ttt struct {
	ti *time.Time
}

func (ttt *ttt)tm() *time.Time {
	return ttt.ti
}
func Test_time1(t *testing.T){

	ti:=time.Now()
	log.Info(unsafe.Pointer(&ti))

	tt:=&ttt{ti:&ti}
	tn:=tt.tm()
	log.Info(unsafe.Pointer(&ti))
	log.Info(unsafe.Pointer(tn))
}
func Test_time3(t *testing.T){

	c:=make(chan bool)


	go func() {

		log.InfoF("========")
		c<-true
		log.InfoF("!!!!!")
		c<-false
		log.InfoF("=====")

	}()
	time.Sleep(time.Second*5)
	log.InfoF("!!!!!!{}",<-c)
	time.Sleep(time.Second*5)
	//log.InfoF("!!!!!!{}",<-c)

}
func Test_time4(t *testing.T){

	for i:=0;i<1000000;i++{
		go func() {
			time.Sleep(time.Minute)
			log.InfoF("Test_time4")
		}()
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig
}
func Test_time5(t *testing.T){

	for i:=0;i<1000000;i++{
		go func() {
			<-time.After(time.Minute)
			log.InfoF("Test_time5")
		}()
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig
}
func Test_time6(t *testing.T){

	log.InfoF("Test_time6")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig
}