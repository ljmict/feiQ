package main

import (
	"feiQ/config"
	"feiQ/recvdata"
	"fmt"
	"net"
	"strconv"
	"time"
)

func CreateUDPSocket() {
	//解析UDPAddr
	addr, _ := net.ResolveUDPAddr("udp", "10.211.55.2:2425")
	//创建udp套接字
	config.UDPSocket, _ = net.ListenUDP("udp", addr)
}


//组装飞鸽传书的数据包
func BuildMsg(command int, optionData string) string {
	msg := config.FeiQVersion + ":" + strconv.FormatInt(time.Now().Unix(), 10) +
		":" + config.FeiQUserName + ":" + config.FeiQHostName + ":" + strconv.Itoa(command) + ":" + optionData
	return msg
}


//发送消息
func SendMsg(msg string, BroadCastIP *net.UDPAddr) {
	config.UDPSocket.WriteToUDP([] byte(msg), BroadCastIP)
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
	var DestIpMsg, SendData string
	fmt.Print("请输入目标IP：")
	fmt.Scanf("%s", &DestIpMsg)
	fmt.Print("发送消息：")
	fmt.Scanf("%s", &SendData)

	DestIp, _ := net.ResolveUDPAddr("udp", DestIpMsg+":"+"2425")
	//如果DestIpMsg和SendData都有内容才能发送数据
	if DestIpMsg != "" && SendData != "" {
		msg := BuildMsg(config.IPMSGSendMsg, SendData)
		SendMsg(msg, DestIp)
	}
}


//命令菜单
func CommandMenu() int {
	var CommandNum int
	fmt.Println("飞鸽传书v1.0")
	fmt.Println("1.发送上线广播")
	fmt.Println("2.发送下线广播")
	fmt.Println("3.发送消息")
	fmt.Println("0.退出程序")
	fmt.Print("请输入数字：")
	fmt.Scanf("%d", &CommandNum)
	return CommandNum
}

func main() {
	//创建udp套接字
	CreateUDPSocket()

	defer config.UDPSocket.Close()

	//循环接收数据
	go recvdata.RecvMsg()

	//命令菜单
	for {
		CommandNum := CommandMenu()
		if CommandNum == 1 {
			//发送广播上线消息
			SendBroadcastOnline()
		} else if CommandNum == 2 {
			//发送广播下线消息
			SendBroadcastOffline()
		} else if CommandNum == 3 {
			//发送消息
			SendDestIpMsg()
		} else if CommandNum == 0 {
			//先发送下线消息，再退出程序
			SendBroadcastOffline()
			break
		}
	}
}
