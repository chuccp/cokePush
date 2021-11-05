package km

type km interface {
	ReadMessage() (Message,error)
}

type km00001 struct {

}

func NewKm00001()*km00001  {

	return &km00001{}
}
func ( km *km00001)ReadMessage() (Message,error){

	return nil,nil
}