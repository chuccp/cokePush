package ex

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/util"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
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
