package define

const (
	// auth user
	OP_AUTH             = int32(1)     //发起验证（连接时）
	OP_AUTH_REPLY       = int32(2)     //验证结果
	// heartbeat
	OP_HEARTBEAT        = int32(3)    //心跳
	OP_HEARTBEAT_REPLY  = int32(4)    //心跳回复
	// send text messgae
	OP_SEND_SMS         = int32(5)     //发消息，暂无使用，目前使用http接口发消息
	OP_SEND_SMS_REPLY   = int32(6)     //回复消息
	
	//replies
	OP_REPLIES          = int32(7)     //回复列表
	
	// kick user
	OP_DISCONNECT_REPLY = int32(8)    //踢用户下线
	
	
	// raw message
	OP_RAW              = int32(11)    //原始数据，暂未使用，用于消息合并
	// room
	OP_ROOM_READY       = int32(12)       //room就绪信号
	// proto
	OP_PROTO_READY      = int32(13)     
	OP_PROTO_FINISH     = int32(14)
	
	//lottery
	OP_LOTTERY_NOTIFY   = int32(21)   //中奖结果通知
	OP_LOTTERY_WINERS   = int32(22)   //中奖列表
	OP_LOTTERY_UPDATE   = int32(23)   //更新抽奖
	
	// for test
	OP_TEST             = int32(254)
	OP_TEST_REPLY       = int32(255)
)
