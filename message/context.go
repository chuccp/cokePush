package message

const (

	// OrdinaryMessageClass 普通信息，用于用户之间信息交流
	OrdinaryMessageClass byte = 1
	// BasicMessageType 一对一发送方信息
	BasicMessageType byte = 1
	// MultiMessageType 一对多发送信息
	MultiMessageType = 2
	// GroupMessageType 群消息
	GroupMessageType byte = 3
	// BroadcastMessageType 广播信息，给所有人发信息
	BroadcastMessageType byte = 4

	FromUser byte = 1
	ToUser   byte = 2
	Text     byte = 3
	GroupId  byte = 4

	// FunctionMessageClass 功能信息，用于查询信息，修改配置
	FunctionMessageClass byte = 2

	LoginType byte = 1

	Username byte = 1

	JoinGroupType byte = 2
)
