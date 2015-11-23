/*
群聊天信息转发
*/
package controller

import (
	"aTalk/protocol"
	"aTalk/server/common"
	"aTalk/server/models"
	"code.google.com/p/go.net/websocket"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2/bson"
	"net"
	"time"
)

func GroupMsgController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("message")
	collectionGroup := common.DBMongo.C("group")
	var groupMsg models.Message
	var group models.Group
	nowTime := time.Now()

	//判断此成员是否在此群组中
	n, err := collectionGroup.Find(bson.M{"_id": msg.Group.GetId(), "members.name": msg.UserInfo.GetUsername()}).Count()
	if err != nil {
		common.CommonResultReturn(conn, "groupMsgResult", msg.GetMsgTypeId(), 2)
		Logger.Error("count :", err)
		return
	}
	if n == 0 {
		//此用户不在此群组中
		common.CommonResultReturn(conn, "groupMsgResult", msg.GetMsgTypeId(), 1)
		return
	}

	//保存聊天记录到数据库
	groupMsg.Id = bson.NewObjectId()
	groupMsg.Sender = msg.UserInfo.GetUsername()
	groupMsg.GroupId = msg.Group.GetId()
	groupMsg.MsgType = int32(msg.SendMsg.GetType())
	groupMsg.Msg = msg.SendMsg.GetMsg()
	groupMsg.MsgTime = nowTime

	err = collection.Insert(&groupMsg)
	if err != nil {
		Logger.Error("insert group msg:", err)
	}

	//转发群组信息给群内所有人
	//获取群组员
	err = collectionGroup.Find(bson.M{"_id": msg.Group.GetId()}).Select(bson.M{"members": true}).All(&group)
	if err != nil {
		Logger.Error("select group member: ", err)
		common.CommonResultReturn(conn, "groupMsgResult", msg.GetMsgTypeId(), 3)
		return
	}
	common.CommonResultReturn(conn, "groupMsgResult", msg.GetMsgTypeId(), 0)
	//发送人
	user := &protocol.User{}
	user.Username = &groupMsg.Sender

	for _, v := range group.Members {
		//把发送消息的用户排除
		if v.Name == msg.UserInfo.GetUsername() {
			continue
		}
		//把群聊信息保存到数据库
		var message models.Message
		message.Id = bson.NewObjectId()
		message.GroupId = msg.Group.GetId()
		message.Receiver = v.Name
		message.MsgType = int32(msg.SendMsg.GetType())
		message.MsgTime = nowTime
		if collection.Insert(&message) != nil {
			Logger.Error("insert group message:", err)
		}

		//如果此成员在线，给其发送消息
		if OnlineCheck(v.Name) {
			sendMsg := &protocol.SendMessage{}

			sendMsg.Receiver = &v.Name
			sendMsg.Type = msg.SendMsg.Type
			sendMsg.Msg = &groupMsg.Msg
			msgTime := nowTime.Unix()
			sendMsg.MsgTime = &msgTime

			returnMsg := &protocol.WMessage{
				MsgType: proto.String("groupMsgResult"),
				//MsgTypeId: proto.Int32(16),
				StataCode: proto.Int32(0),
				SendMsg:   sendMsg,
				Group: &protocol.Groups{
					Id: proto.String(msg.Group.GetId()),
				},
				UserInfo: user,
			}

			buffer, err := proto.Marshal(returnMsg)
			if err != nil {
				Logger.Warn("protobuf marshal failed: ", err)
				continue
			}
			newbuffer := append(common.IntToBytes(len(buffer)), buffer...)
			switch UserMap[v.Name].(type) {
			case net.Conn:
				_, err = (UserMap[v.Name]).(net.Conn).Write(newbuffer)
				common.CheckErr(err)
			case websocket.Conn:
				_, err = (UserMap[v.Name]).(*websocket.Conn).Write(newbuffer)
				common.CheckErr(err)
			}
		}
	}
}
