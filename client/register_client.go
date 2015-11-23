package main

import (
	"aTalk/protocol"
	"aTalk/server/common"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"time"
)

func main() {

	pTCPAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:7979")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		return
	}
	fmt.Println("正在连接....", pTCPAddr)
	pTCPConn, err := net.DialTCP("tcp", nil, pTCPAddr)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		return
	}
	defer pTCPConn.Close()
	input := ""
	fmt.Println("连接成功...")
	go readServer(pTCPConn)
	//write(pTCPConn)
	time.Sleep(time.Second * 1)
	write(pTCPConn)
	for {
		fmt.Scanf("%s", &input)
		//addFriendWrite(pTCPConn)
	}
	time.Sleep(time.Second * 10)
}

func write(conn *net.TCPConn) {
	reg := &protocol.WMessage{
		MsgType:   proto.String("register"),
		MsgTypeId: proto.Int32(1),
		UserInfo: &protocol.User{
			Username: proto.String("zhang"),
			Password: proto.String("123456"),
		},
	}
	buf, err := proto.Marshal(reg)
	if err != nil {
		fmt.Println("failed: %s\n", err)
		return
	}
	fmt.Println("buf: ", len(buf))
	length := len(buf)
	buffer := append(common.IntToBytes(length), buf...)
	conn.Write(buffer)
	//return buffer
	//conn.Write(common.IntToBytes(length))
}

func readServer(pTCPConn *net.TCPConn) {
	//buffer := make([]byte, 1024)
	tempBuf := make([]byte, 4)
	for {
		bufferSize, err := pTCPConn.Read(tempBuf)
		if err != nil {
			fmt.Println("read error: ", err)
			break
		}
		if bufferSize >= 4 {
			bodyLen := BytesToInt(tempBuf)
			fmt.Println("bodyLen: ", bodyLen)
			bodyBuffer := make([]byte, bodyLen)
			_, err := pTCPConn.Read(bodyBuffer)
			if err != nil {
				fmt.Println("err : ", err)
				break
			}
			msg := &protocol.WMessage{}
			err = proto.Unmarshal(bodyBuffer, msg)
			if err != nil {
				fmt.Println("unmarshal：", err)
				return
			}
			fmt.Println("解码：", msg.String())
		}
	}
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}
