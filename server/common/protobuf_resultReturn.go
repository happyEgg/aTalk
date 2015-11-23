package common

import (
	"aTalk/protocol"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
)

//使用空包作为心跳包
func KeepAliveResult(conn interface{}) {
	n := IntToBytes(0)
	switch conn.(type) {
	case net.Conn:
		_, err := conn.(net.Conn).Write(n)
		CheckErr(err)
	case websocket.Conn:
		_, err := conn.(*websocket.Conn).Write(n)
		CheckErr(err)
	}
}

//通用的服务端返回给客户端的数据协议
func CommonResultReturn(conn interface{}, result string, resultId int32, num int32) {
	returnMsg := &protocol.WMessage{
		MsgType:   proto.String(result),
		MsgTypeId: proto.Int32(resultId),
		StataCode: proto.Int32(num),
	}

	buffer, err := proto.Marshal(returnMsg)
	if err != nil {
		Logger.Warn("protobuf marshal failed: ", err)
		fmt.Println("failed: ", err)
		//	Log.SetLogger("file", `{"filename": "savelog.log"}`)
		return
	}
	length := len(buffer)
	newbuffer := append(IntToBytes(length), buffer...)
	switch conn.(type) {
	case net.Conn:
		_, err := conn.(net.Conn).Write(newbuffer)
		CheckErr(err)
	case websocket.Conn:
		_, err := conn.(*websocket.Conn).Write(newbuffer)
		CheckErr(err)
	}
	return
}
