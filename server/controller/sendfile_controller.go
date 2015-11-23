/*
发送文件
*/
package controller

import (
	"aTalk/protocol"
	"aTalk/server/common"
	"code.google.com/p/go.net/websocket"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

func SendFile(conn interface{}, msg *protocol.WMessage) {

	//判断接收方是否在线
	if !OnlineCheck(msg.SendMsg.GetReceiver()) {
		common.CommonResultReturn(conn, "fileResult", msg.GetMsgTypeId(), 1)
		return
	}

	FileMark = msg.GetMsgTypeId()
	//是否同意接收文件
	now := time.Now().Unix()
	msg.SendMsg.MsgTime = &now
	if FileSytle == 0 && int32(msg.SendMsg.GetType()) == 2 {
		fileType := "fileResult"
		msg.MsgType = &fileType
		buffer, err := proto.Marshal(msg)
		if err != nil {
			common.CommonResultReturn(conn, "fileResult", FileMark, 2)
			Logger.Error("sendfile:", err)
			return
		}
		newbuffer := append(common.IntToBytes(len(buffer)), buffer...)
		switch (UserMap[msg.SendMsg.GetReceiver()]).(type) {
		case net.Conn:
			_, err = (UserMap[msg.SendMsg.GetReceiver()]).(net.Conn).Write(newbuffer)
			common.CheckErr(err)
		case websocket.Conn:
			_, err = (UserMap[msg.SendMsg.GetReceiver()]).(*websocket.Conn).Write(newbuffer)
			common.CheckErr(err)
		}
		return
	}

	// if FileSytle == 2 {
	// 	//接收文件名
	// 	fd, err := os.OpenFile("./uploads/"+msg.UserInfo.GetUsername()+"/"+msg.SendMsg.GetMsg(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	// 	if err != nil {
	// 		common.CommonResultReturn(conn, "fileResult", msg.GetMsgTypeId(), 2)
	// 		Logger.Error("send file", err.Error())
	// 		return
	// 	}
	// 	defer fd.Close()

	// 	//写文件
	// 	_, err = fd.Write(msg.SendMsg.GetFile())
	//给接收方发送文件内容
	buffer, err := proto.Marshal(msg)
	if err != nil {
		common.CommonResultReturn(conn, "fileResult", msg.GetMsgTypeId(), 3)
		Logger.Warn("protobuf marshal failed: ", err)
		return
	}
	newbuffer := append(common.IntToBytes(len(buffer)), buffer...)
	switch (UserMap[msg.SendMsg.GetReceiver()]).(type) {
	case net.Conn:
		_, err = (UserMap[msg.SendMsg.GetReceiver()]).(net.Conn).Write(newbuffer)
		common.CheckErr(err)
	case websocket.Conn:
		_, err = (UserMap[msg.SendMsg.GetReceiver()]).(*websocket.Conn).Write(newbuffer)
		common.CheckErr(err)
	}
}

func RecvFile(conn interface{}, msg *protocol.WMessage) {
	//判断接收方是否在线
	if !OnlineCheck(msg.SendMsg.GetReceiver()) {
		common.CommonResultReturn(conn, "recvResult", msg.GetMsgTypeId(), 1)
		return
	}
	if int32(msg.SendMsg.GetType()) == 3 {
		FileSytle = 2
		common.CommonResultReturn(UserMap[msg.SendMsg.GetReceiver()], "recvResult", FileMark, 3)
	} else if int32(msg.SendMsg.GetType()) == 4 {
		common.CommonResultReturn(UserMap[msg.SendMsg.GetReceiver()], "recvResult", FileMark, 0)
	}
}
