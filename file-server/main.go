package main

import (
	"aTalk/protocol"
	"aTalk/server/common"
	"aTalk/server/config"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
	"os"
)

var (
	file_host = config.String("file_host")
)

func main() {
	listener, err := net.Listen("tcp", file_host)
	if err != nil {
		common.Logger.Error("listen:", err)
		os.Exit(-1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			common.Logger.Warn("accept:", err)
			continue
		}

		go connection(conn)
	}
}

func connection(conn net.Conn) {
	defer conn.Close()
	headBuffer := make([]byte, 4)

	//处理文件传输
	bufSize, err := conn.Read(headBuffer)
	if err != nil {
		common.Logger.Error("conn.Read:", err)
		break
	}

	if bufSize < 4 {
		break
	}
	// messager := make(chan byte)
	// //心跳计时
	// go HeadBeating(conn, messager, 2)
	// //检测client是否有数据传来
	// go GravelChannel(headBuffer, messager)

	bodyLen := common.BytesToInt(headBuffer)
	bodyBuffer := make([]byte, bodyLen)
	bodySize, err := conn.Read(bodyBuffer)
	if err != nil {
		break
	}
	if bodySize < bodyLen {
		break
	}

	msg := &protocol.WMessage{}
	err = proto.Unmarshal(bodyBuffer, msg)
	if err != nil {
		common.Logger.Warn("proto:", err)
		break
	}

	switch msg.GetMsgType() {
	case "sendfile":
		handleFile(conn, msg)
	}
}

func handleFile(conn net.Conn, msg *protocol.WMessage) {
	buffer := make([]byte, 1024)
	fd, err := os.OpenFile("./uploads/"+msg.UserInfo.GetUsername()+"/"+msg.SendMsg.GetFile(),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		common.Logger.Error("os.openfile:", err)
		return
	}
	defer fd.Close()
	for {
		c, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			common.Logger.Error("read:", err)
			panic(err)
		}
		if c == 0 {
			break
		}

		_, err = fd.Write(buffer[:c])
		if err != nil {
			common.Logger.Error("write:", err)
			break
		}

	}
}
