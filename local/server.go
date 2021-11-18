package local

type Server struct {


}

func NewServer() *Server {
	return &Server{}
}

func (server *Server)Start()error{
	return nil
}
func (server *Server)Init(){

}
func (server *Server)Name()string{
	return "local"
}