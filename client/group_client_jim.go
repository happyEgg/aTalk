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
	zhanBaoWrite(pTCPConn)
	for {
		fmt.Scanf("%s", &input)
		write(pTCPConn)
	}
	//zhanBaoWrite(pTCPConn)
	//for {
	//fmt.Scanf("%s", input)
	//writebuf := write()
	//Write(pTCPConn)
	//go readServer(pTCPConn)
	//time.Sleep(time.Second * 1 / 1000)
	//}
	time.Sleep(time.Second * 10)
}

func write(conn *net.TCPConn) {

	//n := protocol.Mode(5)
	reg := &protocol.WMessage{
		MsgType:   proto.String("groupMsg"),
		MsgTypeId: proto.Int32(16),
		UserInfo: &protocol.User{
			Username: proto.String("jim"),
			//Password: proto.String("123456"),
		},
		Group: &protocol.Groups{
			//Receiver: proto.String("zhang"),
			GroupName: proto.String("chi fan qun"),
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

func zhanBaoWrite(conn *net.TCPConn) {

	login := &protocol.WMessage{
		MsgType:   proto.String("login"),
		MsgTypeId: proto.Int32(1),
		UserInfo: &protocol.User{
			Username: proto.String("jim"),
			Password: proto.String("123456"),
		},
	}
	loginbuf, err := proto.Marshal(login)
	if err != nil {
		fmt.Println("failed: %s\n", err)
		return
	}
	loginlength := len(loginbuf)
	loginbuffer := append(common.IntToBytes(loginlength), loginbuf...)
	conn.Write(loginbuffer)
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
