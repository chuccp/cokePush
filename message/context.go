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
	// ExMessageId 补充的msgId 兼容老的协议/**
	ExMessageId byte = 5

	// FunctionMessageClass 功能信息，用于查询信息，修改配置
	FunctionMessageClass byte = 2

	LoginType byte = 1

	Username byte = 1

	//JoinGroupType 加入群
	JoinGroupType byte = 2

	QueryMachineType = 3

	QueryMachineBasicType = 4

	QueryMachineInfoType = 5

	QueryMachineInfo = 5

	LocalMachineAddress = 1

	DeleteUserType = 6

	AddUserType = 7

	USERNAME = 1

	MaChineId = 2

	QueryType = 8

	QueryName = 1

	// LiveMessageClass 心跳信息，用于心跳测试以及集群信息交流
	LiveMessageClass byte = 3
	//BlankLiveType 简单心跳检查，发送空包
	BlankLiveType = 1

	// BackMessageClass 返回消息
	BackMessageClass byte = 4

	// BackMessageOKType 返回成功
	BackMessageOKType byte = 1

	BackMachineAddress = 2

	QueryData = 3

	// BackMessageErrorType 返回异常
	BackMessageErrorType byte = 2

	FunctionBackMessageClass byte = 5

	OrdinaryMessageBackMessageClass byte = 6
)
