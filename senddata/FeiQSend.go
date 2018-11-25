package senddata

import (
	"bufio"
	"feiQ/config"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

//组装飞鸽传书的数据包
func BuildMsg(command int, optionData string) string {
	msg := config.FeiQVersion + ":" + strconv.FormatInt(time.Now().Unix(), 10) +
		":" + config.FeiQUserName + ":" + config.FeiQHostName + ":" + strconv.Itoa(command) + ":" + optionData
	return msg
}


//发送消息
func SendMsg(msg string, DestIP *net.UDPAddr) {
	config.UDPSocket.WriteToUDP([] byte(msg), DestIP)
}


//发送上线广播消息
func SendBroadcastOnline() {
	msg := BuildMsg(config.IPMSGBrEntry, config.FeiQUserName)

	BroadCastIP := net.UDPAddr{
		IP:   net.IPv4(10, 211, 55, 255),
		Port: 2425,
	}
	SendMsg(msg, &BroadCastIP)
}


//发送下线广播消息
func SendBroadcastOffline() {
	msg := BuildMsg(config.IPMSGBrEXIT, config.FeiQUserName)

	BroadCastIP := net.UDPAddr{
		IP:   net.IPv4(10, 211, 55, 255),
		Port: 2425,
	}
	SendMsg(msg, &BroadCastIP)
}


//向指定ip发送消息
func SendDestIpMsg() {
	var destIpMsg string
	fmt.Print("请输入目标IP：")
	fmt.Scanf("%s", &destIpMsg)
	fmt.Print("发送消息：")
	reader := bufio.NewReader(os.Stdin)
	sendData, _, _ := reader.ReadLine()

	destIp, _ := net.ResolveUDPAddr("udp", destIpMsg+":"+"2425")
	//如果DestIpMsg和SendData都有内容才能发送数据
	if destIpMsg != "" && string(sendData) != "" {
		msg := BuildMsg(config.IPMSGSendMsg, string(sendData))
		SendMsg(msg, destIp)
	}
}