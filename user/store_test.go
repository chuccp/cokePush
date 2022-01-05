package user

import (
	"github.com/chuccp/cokePush/message"
	"log"
	"strconv"
	"testing"
	"time"
	"unsafe"
)

type testUser struct {
	username string
}

func newTestUser(username string)*testUser  {
	tu:= &testUser{}
	tu.SetUsername(username)
	return tu
}
func (testUser*testUser)SetUsername(username string){
	testUser.username = username
}
func (testUser*testUser)GetUsername() string{
	return "username"
}
func (testUser*testUser)GetUserId() string{
	return testUser.GetUsername()+strconv.FormatUint(uint64(uintptr(unsafe.Pointer(testUser))),36)
}
func (testUser*testUser)WriteMessage(iMessage message.IMessage) error{
	return nil
}


func Test_user2(t *testing.T) {


	var sss1 = newUserStore()

	 num:=  10_000

	nin:= time.Now().Nanosecond()

	for i:=0;i<num;i++{
		strconv.FormatUint(uint64(uintptr(unsafe.Pointer(sss1))),10)
	}
	log.Println(time.Now().Nanosecond()-nin)

	nin= time.Now().Nanosecond()
	for i:=0;i<num;i++ {

		strconv.FormatUint(uint64(uintptr(unsafe.Pointer(sss1))),36)
	}
	log.Println(time.Now().Nanosecond()-nin)

}