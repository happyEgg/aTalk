package main

import (
	"aTalk/protocol"
	"aTalk/server/common"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
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
	login(pTCPConn)
	time.Sleep(time.Second * 1)
	recvfile(pTCPConn)
	time.Sleep(time.Second * 10)
}

func login(conn *net.TCPConn) {

	login := &protocol.WMessage{
		MsgType:   proto.String("login"),
		MsgTypeId: proto.Int32(1),
		UserInfo: &protocol.User{
			Username: proto.String("zhang"),
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

	tempBuf := make([]byte, 4)
	for {
		size, err := pTCPConn.Read(tempBuf)
		if err != nil {
			panic(err)
		}

		if size < 4 {
			return
		}

		bodysize := BytesToInt(tempBuf)
		bodybuffer := make([]byte, bodysize)
		_, err = pTCPConn.Read(bodybuffer)
		if err != nil {
			fmt.Println("err : ", err)
			break
		}

		pb := &protocol.WMessage{}
		err = proto.Unmarshal(bodybuffer, pb)
		if err != nil {
			fmt.Println("unmarshal：", err)
			return
		}

		fd, err := os.OpenFile("./upload/"+pb.UserInfo.GetUsername()+"/"+pb.SendMsg.GetFile(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("err:", err)
			break
		}

		defer fd.Close()
		//filesize := pb.SendMsg.GetFileSize()
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
				//fmt.Println("解码：", msg.String())
				_, err = fd.Write([]byte(msg.SendMsg.GetFile()))
				if err != nil {
					fmt.Println("err:", err)
					return
				}
			}
		}
	}
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}
