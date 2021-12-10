package message

type BasicMessage struct {
	*Message
}
func CreateBasicMessage(fromUser string,toUser string,messageText string)* BasicMessage {
	bm:=&BasicMessage{Message:CreateMessage(OrdinaryMessageClass,BasicMessageType)}
	bm.SetString(FromUser,fromUser)
	bm.SetString(ToUser,toUser)
	bm.SetString(Text,messageText)
	return bm
}