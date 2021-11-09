package km

import "github.com/chuccp/cokePush/message"

type km interface {
	ReadMessage() (message.Message,error)
}

type km00001 struct {

}

func NewKm00001()*km00001  {

	return &km00001{}
}
func ( km *km00001)ReadMessage() (message.Message,error){

	return nil,nil
}