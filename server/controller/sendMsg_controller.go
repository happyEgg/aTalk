/*
好友之间的消息转发
*/
package controller

import (
	"aTalk/protocol"
	"aTalk/server/common"
	"aTalk/server/models"
	"code.google.com/p/go.net/websocket"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
	"time"
)

func SendMsgController(conn interface{}, msg *protocol.WMessage) {
	collectionUser := common.DBMongo.C("user")
	collection := common.DBMongo.C("message")

	//好友之间消息的转发

	//判断两者是不是好友
	n, err := collectionUser.Find(bson.M{"user_name": msg.UserInfo.GetUsername(), "friends.name": msg.SendMsg.GetReceiver()}).Count()
	if err != nil {
		common.CommonResultReturn(conn, "sendMsgResult", msg.GetMsgTypeId(), 2)
		Logger.Error("send msg:", err)
		return
	}

	if n != 0 {
		nowTime := time.Now()
		log.Println("time:", nowTime)
		time := nowTime.Unix()
		msg.SendMsg.MsgTime = &time
		buffer, err := proto.Marshal(msg)
		if err != nil {
			common.CommonResultReturn(conn, "sendMsgResult", msg.GetMsgTypeId(), 2)
			Logger.Warn("protobuf marshal failed: ", err)
			return
		}
		newbuffer := append(common.IntToBytes(len(buffer)), buffer...)
		if OnlineCheck(msg.SendMsg.GetReceiver()) {
			switch (UserMap[msg.SendMsg.GetReceiver()]).(type) {
			case net.Conn:
				_, err = (UserMap[msg.SendMsg.GetReceiver()]).(net.Conn).Write(newbuffer)
				common.CheckErr(err)
			case websocket.Conn:
				_, err = (UserMap[msg.SendMsg.GetReceiver()]).(*websocket.Conn).Write(newbuffer)
				common.CheckErr(err)
			}
		}
		common.CommonResultReturn(conn, "sendMsgResult", msg.GetMsgTypeId(), 0)
		var message models.Message
		message.Id = bson.NewObjectId()
		message.Sender = msg.UserInfo.GetUsername()
		message.MsgType = int32(msg.SendMsg.GetType())
		message.Receiver = msg.SendMsg.GetReceiver()
		message.Msg = msg.SendMsg.GetMsg()
		message.MsgTime = nowTime

		err = collection.Insert(&message)
		if err != nil {
			Logger.Error("insert message: ", err)
			log.Println("message: ", err)
		}

	} else {
		common.CommonResultReturn(conn, "sendMsgResult", msg.GetMsgTypeId(), 1)
	}
}
