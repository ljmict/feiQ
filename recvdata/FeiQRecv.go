package recvdata

import (
	"feiQ/config"
	"fmt"
)

func RecvMsg() {
	for {
		buf := make([] byte, 1024)
		dataLen, _ := config.UDPSocket.Read(buf)
		fmt.Printf("收到消息：%s\n", string(buf[:dataLen]))
	}
}
