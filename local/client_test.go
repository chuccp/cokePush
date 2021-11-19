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


	//client:=newClient(u)
	//bm := message.CreateBasicMessage("222222","222222" , "444444")
	//client.handleMessage(bm)

}
func TestServer(t *testing.T) {
	var defaultRegister = core.NewRegister()
	server:=NewServer()
    defaultRegister.AddServer(server)
	cokePush:=defaultRegister.Create()
	cokePush.StartSync()
	u:=NewUser("222222",WriteMessage)
	bm := message.CreateBasicMessage("222222","222222" , "444444")
	client:=server.CreateClient(u)
	client.login()
	client.handleMessage(bm)


}