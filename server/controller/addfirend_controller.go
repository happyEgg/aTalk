/*
添加好友的请求
*/
package controller

import (
	"aTalk/protocol"
	"aTalk/server/common"
	"aTalk/server/models"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2/bson"
	"net"
	"time"
)

//发送好友请求
func SendFirendRequestController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("user")
	nowTime := time.Now()

	//判断需要添加的好友是否存在
	exist, err := collection.Find(bson.M{"user_name": msg.AddFriend.GetReceiver()}).Count()
	if err != nil {

		common.CommonResultReturn(conn, "addFriendResult", msg.GetMsgTypeId(), 3)
		Logger.Error("get friend relation:", err)
		return
	}
	if exist != 0 {

		//验证是否已经是其好友
		n, err := collection.Find(bson.M{"user_name": msg.AddFriend.GetSender(), "friends.name": msg.AddFriend.GetReceiver(), "friends.mode": 3}).Count()
		if err != nil {
			common.CommonResultReturn(conn, "addFriendResult", msg.GetMsgTypeId(), 3)
			Logger.Error("get friend relation:", err)
			return
		}

		//已经是好友关系
		if n != 0 {
			common.CommonResultReturn(conn, "addFriendResult", msg.GetMsgTypeId(), 2)
			return
		}

		//把添加好友信息保存到message表
		if int32(msg.AddFriend.GetModes()) == 2 {
			var saveMsg models.Message
			saveMsg.Id = bson.NewObjectId()
			saveMsg.Sender = msg.AddFriend.GetSender()
			saveMsg.Receiver = msg.AddFriend.GetReceiver()
			saveMsg.MsgType = int32(msg.AddFriend.GetModes())
			saveMsg.Msg = msg.UserInfo.GetUsername() + "请求添加您为好友"
			saveMsg.MsgTime = nowTime

			err = collection.Insert(&saveMsg)
			if err != nil {
				common.CommonResultReturn(conn, "addFriendResult", msg.GetMsgTypeId(), 3)
				Logger.Error("add friend request insert database error:", err)
				return
			}

			//把添加好友的信息发给被添加人
			if OnlineCheck(msg.AddFriend.GetReceiver()) {
				num := nowTime.Unix()
				msg.SendMsg.MsgTime = &num
				buffer, err := proto.Marshal(msg)
				if err != nil {
					common.CommonResultReturn(conn, "addFriendResult", msg.GetMsgTypeId(), 3)
					Logger.Warn("protobuf marshal failed: ", err)
					fmt.Println("failed: %s\n", err)
					//Log.SetLogger("file", `{"filename": "savelog.log"}`)
					return
				}
				length := len(buffer)
				newbuffer := append(common.IntToBytes(length), buffer...)
				switch (UserMap[msg.AddFriend.GetReceiver()]).(type) {
				case net.Conn:
					_, err = (UserMap[msg.AddFriend.GetReceiver()]).(net.Conn).Write(newbuffer)
					common.CheckErr(err)
				case websocket.Conn:
					_, err = (UserMap[msg.AddFriend.GetReceiver()]).(*websocket.Conn).Write(newbuffer)
					common.CheckErr(err)
				}
			}
			common.CommonResultReturn(conn, "addFriendResult", msg.GetMsgTypeId(), 0)
		} else {
			//发送的数据不正确
			common.CommonResultReturn(conn, "addFriendResult", msg.GetMsgTypeId(), 4)
		}

	} else {
		//被添加人不存在
		common.CommonResultReturn(conn, "addFriendResult", msg.GetMsgTypeId(), 1)
	}
}
