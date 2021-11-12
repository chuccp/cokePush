package message

type BasicMessage struct {
	*Message
}
func (bm *BasicMessage)GetMessageType() byte{
	return BasicMessageType
}

func CreateBasicMessage(fromUser string,toUser string,messageText string)* BasicMessage {
	bm:=&BasicMessage{Message:CreateMessage()}
	bm.SetString(FromUser,fromUser)
	bm.SetString(ToUser,toUser)
	bm.SetString(Text,messageText)
	return bm
}