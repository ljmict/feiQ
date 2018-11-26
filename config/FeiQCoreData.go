package config

import "net"

var UDPSocket *net.UDPConn = nil

const (
	FeiQPort     = "2425"
	FeiQVersion  = "1"
	FeiQUserName = "Jack"
	FeiQHostName = "mac-pro"
	Broadcast    = "255.255.255.255"
	IPMSGBrEntry = 0x00000001 //上线提醒消息命令
	IPMSGBrEXIT  = 0x00000002 //下线提醒消息命令
	IPMSGSendMsg = 0x00000020 //表示发送消息
	IPMSGAnsentry = 0x00000003 //对方也在线
	IPMSGRecvMsg = 0x00000021 //确认收到消息
	IPMSGFileAttachOpt = 0x00200000
	IPMSGFileRegular = 0x00000001 //普通文件
)

//保存在线用户列表
var UserSlice []map[string]string

