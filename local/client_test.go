package local

import (
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
	"log"
	"testing"
)

func WriteMessage(iMessage message.IMessage)error  {

	log.Println(iMessage.GetString(message.Text))

	return nil
}

func TestEqual(t *testing.T) {

	u:=NewUser("222222",WriteMessage)
	client:=newClient(u)
	bm := message.CreateBasicMessage("222222","222222" , "444444")
	client.handleMessage(bm)

}
func TestServer(t *testing.T) {
	var defaultRegister = core.NewRegister()
    defaultRegister.AddServer(NewServer())
	cokePush:=defaultRegister.Create()
	cokePush.StartSync()

}