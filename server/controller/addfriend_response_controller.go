/*
同意或拒绝添加好友的请求
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
	"log"
	"net"
)

//添加好友应答
func SendFirendResponseController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("user")
	var user models.User

	//判断需要同意的好友是否存在
	exist, err := collection.Find(bson.M{"user_name": msg.AddFriend.GetReceiver()}).Count()
	if err != nil {
		common.CommonResultReturn(conn, "sendFriendResponse", msg.GetMsgTypeId(), 1)
		Logger.Error("get friend relation:", err)
		return
	}
	if exist != 0 {

		//验证是否已经是其好友
		num, err := collection.Find(bson.M{"user_name": msg.AddFriend.GetSender(), "friends.name": msg.AddFriend.GetReceiver(), "friends.mode": 3}).Count()
		if err != nil {
			Logger.Error("get friend relation:", err)
			return
		}

		//其已经是好友关系
		if num != 0 {
			common.CommonResultReturn(conn, "sendFriendResponse", msg.GetMsgTypeId(), 2)
			return
		}

		//如果同意添加好友，就写入到数据库
		n := int32(msg.AddFriend.GetModes())
		if n == 3 {
			friend := &models.Friend{
				Name: msg.AddFriend.GetReceiver(),
				Mode: 3,
			}
			friend1 := &models.Friend{
				Name: msg.AddFriend.GetSender(),
				Mode: 3,
			}
			err := collection.Update(bson.M{"user_name": msg.AddFriend.GetSender()}, bson.M{"$push": bson.M{"friends": friend}})
			err = collection.Update(bson.M{"user_name": msg.AddFriend.GetReceiver()}, bson.M{"$push": bson.M{"friends": friend1}})
			if err != nil {
				Logger.Error("insert friend relation failed: ", err)
				common.CommonResultReturn(conn, "sendFriendResponse", msg.GetMsgTypeId(), 3)
				return
			}

			//添加好友成功，把好友的信息返回
			err = collection.Find(bson.M{"user_name": msg.AddFriend.GetReceiver()}).One(&user)
			if err != nil {
				Logger.Error("select user info: ", err)
				log.Println("好友内容查找错误")
				return
			}
			msg.UserInfo = common.DataToProto(&user)
			n := int32(0)
			msg.StataCode = &n
			buffer, err := proto.Marshal(msg)
			if err != nil {
				common.CommonResultReturn(conn, "sendFriendResponse", msg.GetMsgTypeId(), 3)
				Logger.Warn("protobuf marshal failed: ", err)
				fmt.Println("failed: %s\n", err)
				return
			}
			length := len(buffer)
			newbuffer := append(common.IntToBytes(length), buffer...)
			switch conn.(type) {
			case net.Conn:
				_, err = conn.(net.Conn).Write(newbuffer)
				common.CheckErr(err)
			case websocket.Conn:
				_, err = conn.(*websocket.Conn).Write(newbuffer)
				common.CheckErr(err)
			}
		}

		//给好友发送通知,来自添加好友的应答
		if OnlineCheck(msg.AddFriend.GetReceiver()) {
			buffer, err := proto.Marshal(msg)
			if err != nil {
				fmt.Println("failed: %s\n", err)
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
	} else {
		common.CommonResultReturn(conn, "sendFriendResponse", msg.GetMsgTypeId(), 1)
	}
}
