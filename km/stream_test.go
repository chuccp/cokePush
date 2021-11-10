package km

import (
	"github.com/chuccp/cokePush/message"
	"testing"
)

func TestStream_ReadMessage(t *testing.T) {

}
func TestStream_WriteMessage(t *testing.T) {

	bm:= message.CreateBasicMessage("123","456","789")
	k:=NewKm00001()
	k.WriteMessage(bm)

}