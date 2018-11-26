package main

import (
	"feiQ/config"
	"feiQ/recvdata"
	"feiQ/senddata"
	"fmt"
	"net"
)

func CreateUDPSocket() {
	//解析UDPAddr
	addr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:2425")
	//创建udp套接字
	config.UDPSocket, _ = net.ListenUDP("udp", addr)
}

//命令菜单
func CommandMenu() int {
	var CommandNum int
	fmt.Println("飞鸽传书v1.0")
	fmt.Println("1.发送上线广播")
	fmt.Println("2.发送下线广播")
	fmt.Println("3.发送消息")
	fmt.Println("4.显示在线用户")
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
			senddata.SendBroadcastOnline()
		} else if CommandNum == 2 {
			//发送广播下线消息
			senddata.SendBroadcastOffline()
		} else if CommandNum == 3 {
			//发送消息
			senddata.SendDestIpMsg()
		} else if CommandNum == 4 {
			//显示所有在线用户
			senddata.DisplayOnlineUser()
		} else if CommandNum == 0 {
			//先发送下线消息，再退出程序
			senddata.SendBroadcastOffline()
			break
		}
	}
}
