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
	sendfile(pTCPConn)
	time.Sleep(time.Second * 10)
}

func sendfile(conn *net.TCPConn) {
	buff := make([]byte, 1024)
	fd, err := os.Open("./123.jpg")
	common.CheckErr(err)
	defer fd.Close()

	fdinfo, err := fd.Stat()
	fmt.Println("这个文件的大小为:", fiinfo.Size(), "字节")
	m := protocol.Mode_FILE
	sendfile := &protocol.WMessage{
		MsgType:   proto.String("sendfile"),
		MsgTypeId: proto.Int32(8),
		SendMsg: &protocol.SendMessage{
			Receiver: proto.String("zhang"),
			MsgType:  &m,
			File:     proto.String("123.jpg"),
			FileSize: proto.Int64(fdinfo.Size()),
		},
	}
	n := protocol.Mode_FILE
	for {
		num, err := fd.Read(buff)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if num == 0 {
			break
		}

		reg := &protocol.WMessage{
			MsgType:   proto.String("sendfile"),
			MsgTypeId: proto.Int32(8),
			UserInfo: &protocol.User{
				Username: proto.String("jim"),
			},
			SendMsg: &protocol.SendMessage{
				Receiver: proto.String("zhang"),
				MsgType:  &n,
				File:     proto.String(buff[:num]),
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
}

func login(conn *net.TCPConn) {

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
