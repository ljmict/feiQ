package senddata

import (
	"bufio"
	"feiQ/config"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//显示所有在线用户
func DisplayOnlineUser() {
	for index, user := range config.UserSlice {
		fmt.Printf("%d %s\n", index, user["userName"])
	}
}

//组装飞鸽传书的数据包
func BuildMsg(command int, optionData string) string {
	msg := config.FeiQVersion + ":" + strconv.FormatInt(time.Now().Unix(), 10) +
		":" + config.FeiQUserName + ":" + config.FeiQHostName + ":" + strconv.Itoa(command) + ":" + optionData
	fmt.Println(msg)
	return msg
}

//构建文件消息
func BuildFileMsg(fileName string) string {
	//文件序号:文件名:文件大小:修改时间:文件的属性
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		panic("文件不存在")
	}

	fileNo := strconv.Itoa(0)
	fileSize := strconv.FormatInt(fileInfo.Size(), 16)
	fileCTime := strconv.FormatInt(fileInfo.ModTime().Unix(), 16)
	fileType := strconv.FormatInt(config.IPMSGFileRegular, 16)
	buildFileMsg := []string{fileNo, fileName, fileSize, fileCTime, fileType}
	optionStr := strings.Join(buildFileMsg, ":")
	fileStr := string("\x00") + optionStr + ":"
	commandNum := config.IPMSGSendMsg | config.IPMSGFileAttachOpt
	fmt.Println(fileStr)
	return BuildMsg(commandNum, fileStr)
}

//发送消息
func SendMsg(msg string, DestIP *net.UDPAddr) {
	config.UDPSocket.WriteToUDP([] byte(msg), DestIP)
}


//发送上线广播消息
func SendBroadcastOnline() {
	msg := BuildMsg(config.IPMSGBrEntry, config.FeiQUserName)

	BroadCastIP := net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 2425,
	}
	SendMsg(msg, &BroadCastIP)
}


//发送下线广播消息
func SendBroadcastOffline() {
	msg := BuildMsg(config.IPMSGBrEXIT, config.FeiQUserName)

	BroadCastIP := net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 2425,
	}
	SendMsg(msg, &BroadCastIP)
}

//向指定ip发送消息
func getDestIP() *net.UDPAddr {
	var destIpMsg string
	index := 0
	fmt.Print("请输入目标IP（输入d显示在线用户）：")
	fmt.Scanf("%s", &destIpMsg)
	if destIpMsg == "d" {
		DisplayOnlineUser()
		fmt.Print("输入对应的序号选择你要发送的用户：")
		fmt.Scanf("%d", &index)
		destIpMsg = config.UserSlice[index]["ip"]
	}
	destIp, _ := net.ResolveUDPAddr("udp", destIpMsg+":"+"2425")
	return destIp
}

//向指定ip发送消息
func SendDestIpMsg() {
	destIp := getDestIP()
	fmt.Print("发送消息：")
	reader := bufio.NewReader(os.Stdin)
	sendData, _, _ := reader.ReadLine()

	//如果DestIpMsg和SendData都有内容才能发送数据
	if destIp != nil && string(sendData) != "" {
		msg := BuildMsg(config.IPMSGSendMsg, string(sendData))
		SendMsg(msg, destIp)
	}
}

//发送文件消息
func SendFileMsg() {
	fileName := ""
	destIp := getDestIP()
	fmt.Println("请输入要发送的文件名（输入d显示当前路径下文件名）：")
	fmt.Scanf("%s", &fileName)
	if fileName == "d" {
		//显示当前目录下所有文件
		fileSlice, _ := ioutil.ReadDir("./")
		for index, file := range fileSlice {
			if ! file.IsDir() {
				fmt.Printf("%d、%s ", index, file.Name())
			}
		}

		fileNum := 0
		fmt.Println("请输入文件序号：")
		fmt.Scanf("%d", &fileNum)
		fileName = fileSlice[fileNum].Name()
	}

	if destIp != nil && fileName != "" {
		fileMsg := BuildFileMsg(fileName)
		SendMsg(fileMsg, destIp)
	}
}