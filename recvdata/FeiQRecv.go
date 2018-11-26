package recvdata

import (
	"feiQ/config"
	"feiQ/senddata"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strconv"
	"strings"
)

//处理接收到的数据
func dealFeiQData(buf []byte, dataLen int) map[string]string {
	strData := string(buf[:dataLen])
	strSlice := strings.Split(strData, ":")

	feiQData := map[string]string{}
	feiQData["feiQVersion"] = strSlice[0]
	feiQData["packetID"] = strSlice[1]
	feiQData["userName"] = strSlice[2]
	feiQData["hostName"] = strSlice[3]
	feiQData["commandStr"] = strSlice[4]

	reader := transform.NewReader(strings.NewReader(strSlice[5]), simplifiedchinese.GBK.NewDecoder())
	byteData, _ := ioutil.ReadAll(reader)
	feiQData["option"] = string(byteData)
	return feiQData
}

//处理命令选项
func dealCommandOptionNum(commandStr string) (int, int) {
	//提取命令字中的命令及选项
	commandNum, _ := strconv.Atoi(commandStr)
	command := commandNum & 0x000000ff
	commandOption := commandNum & 0xffffff00
	return command, commandOption
}

//添加在线用户
func addOnlineUser(userName, hostName string, destIP string) {
	//判断用户是否已经存在UserSlice中，如果没有则添加
	for _, user := range config.UserSlice {
		if user["ip"] == destIP {
			return
		}
	}
	newOnlineUser := map[string]string{}
	newOnlineUser["ip"] = destIP
	newOnlineUser["userName"] = userName
	newOnlineUser["hostName"] = hostName
	config.UserSlice = append(config.UserSlice, newOnlineUser)
}

//删除下线用户
func delOfflineUser(IP string) {
	for index, user := range config.UserSlice {
		if user["ip"] == IP {
			start := config.UserSlice[:index]
			end := config.UserSlice[index+1:]
			config.UserSlice = append(start, end...)
			break
		}
	}
}

//接收数据
func RecvMsg() {
	for {
		buf := make([]byte, 1024)
		dataLen, addr, _ := config.UDPSocket.ReadFromUDP(buf)
		feiQData := dealFeiQData(buf, dataLen)
		command, _ := dealCommandOptionNum(feiQData["commandStr"])
		if command == config.IPMSGBrEntry {
			//有用户上线
			fmt.Printf("%s上线\n", feiQData["option"])
			addOnlineUser(feiQData["option"], feiQData["hostName"], addr.IP.String())

			//通告对方我也在线
			answerOnlineMsg := senddata.BuildMsg(config.IPMSGAnsentry, "")
			senddata.SendMsg(answerOnlineMsg, addr)
		} else if command == config.IPMSGBrEXIT {
			//有用户下线
			fmt.Printf("%s下线\n", feiQData["userName"])
			delOfflineUser(addr.IP.String())
		} else if command == config.IPMSGAnsentry {
			//对方通告在线
			fmt.Printf("%s在线\n", feiQData["userName"])
			addOnlineUser(feiQData["option"], feiQData["hostName"], addr.IP.String())
		} else if command == config.IPMSGSendMsg {
			//接收到消息
			fmt.Printf("%s：%s\n", feiQData["userName"], feiQData["option"])
			//给对方发送消息确认
			msg := senddata.BuildMsg(config.IPMSGRecvMsg, "")
			senddata.SendMsg(msg, addr)
		}
	}
}
