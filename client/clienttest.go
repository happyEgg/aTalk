package main

import (
	"aTalk/protocol"
	"aTalk/server/common"
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
	//time.Sleep(time.Second * 10)
}

func addFriendWrite(conn *net.TCPConn) {

	n := protocol.Mode(2)
	reg := &protocol.WMessage{
		MsgType:   proto.String("sendFriendRequest"),
		MsgTypeId: proto.Int32(1),
		System:    proto.String("IOS"),
		AddFriend: &protocol.AddFriendRequest{
			Sender:   proto.String("jim"),
			Modes:    &n,
			Receiver: proto.String("jj"),
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
}

func write(conn *net.TCPConn) {
	reg := &protocol.WMessage{
		MsgType:   proto.String("modifyInfo"),
		MsgTypeId: proto.Int32(3),
		UserInfo: &protocol.User{
			Username: proto.String("jim"),
			RealName: proto.String("张三"),
			Age:      proto.Int32(20),
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
	fmt.Println("loginbuf: ", len(loginbuf))
	loginlength := len(loginbuf)
	fmt.Println("head:", loginlength)
	loginbuffer := append(common.IntToBytes(loginlength), loginbuf...)
	conn.Write(loginbuffer)
	//return buff
}

func readServer(pTCPConn *net.TCPConn) {
	buffer := make([]byte, 1024)
	tempBuf := make([]byte, 4)
	for {
		bufferSize, err := pTCPConn.Read(buffer)
		if err != nil {
			return
		}
		if bufferSize >= 4 {
			copy(tempBuf, buffer)
			length := common.BytesToInt(tempBuf)

			//fmt.Println("length:", length)
			newbuffer := make([]byte, length)
			copy(newbuffer, buffer[4:])

			msg := &protocol.WMessage{}
			err = proto.Unmarshal(newbuffer, msg)
			if err != nil {
				fmt.Println("unmarshal：", err)
				return
			}
			fmt.Println("解码：", msg.String())
		}

	}
}
