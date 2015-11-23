/*
处理未读的单聊信息和群聊信息
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

//用户上线时推送个人未读的单聊信息
func HandleSingleUnreadMsg(conn interface{}, value *string, user *models.User) {
	collection := (common.DBMongo).C("message")
	var messages []models.Message
	nowTime := time.Now()

	//查找用户上次退出和此次登陆时的未读消息
	err := collection.Find(bson.M{"sender": value, "group_id": bson.M{"$exists": false}, "receiver": user.UserName,
		"send_time": bson.M{"$gte": user.LogoutTime, "$lte": nowTime}}).Sort("-send_time").Limit(PullCount).All(&messages)

	//未读消息的数量
	num, err := collection.Find(bson.M{"sender": value, "group_id": bson.M{"$exists": false}, "receiver": user.UserName,
		"send_time": bson.M{"$gte": user.LogoutTime, "$lte": nowTime}}).Sort("-send_time").Limit(PullCount).Count()
	if err != nil {
		Logger.Error("get unread single msg failed: ", err)
		return
	}

	MsgCount[value] = num

	//循环推送查找到的每一条未读消息
	for _, v := range messages {

		msgUser := &protocol.User{}
		send := &protocol.SendMessage{}
		msgUser.Username = &v.Sender
		send.Receiver = &v.Receiver

		msgType := protocol.Mode(v.MsgType)
		send.Type = &msgType
		send.Msg = &v.Msg
		sendtime := v.MsgTime.Unix()
		send.MsgTime = &sendtime
		returnMsg := &protocol.WMessage{
			MsgType: proto.String("sendMsgResult"),
			//MsgTypeId: proto.Int32(8),
			StataCode: proto.Int32(0),
			SendMsg:   send,
			UserInfo:  msgUser,
		}

		buffer, err := proto.Marshal(returnMsg)
		if err != nil {
			Logger.Warn("protobuf marshal failed: ", err)
			fmt.Println("failed: ", err)
			return
		}
		newbuffer := append(common.IntToBytes(len(buffer)), buffer...)
		switch conn.(type) {
		case net.Conn:
			_, err = conn.(net.Conn).Write(newbuffer)
			common.CheckErr(err)
		case websocket.Conn:
			_, err = conn.(*websocket.Conn).Write(newbuffer)
			common.CheckErr(err)
		}

	}
}

//处理未读的群聊信息
func HandleGroupUnreadMsg(conn interface{}, value *models.Message, user *models.User) {

	//根据群组找其所在的群组的聊天消息
	var messages []models.Message
	collection := (common.DBMongo).C("message")
	nowTime := time.Now()

	//查找用户上次退出和此次登陆时的未读消息
	err := collection.Find(bson.M{"group_id": value.GroupId, "send_time": bson.M{"$gte": user.LogoutTime, "$lte": nowTime}}).
		Sort("-send_time").Limit(PullCount).All(&messages)

	//未读消息的数量
	num, err := collection.Find(bson.M{"sender": value.Sender, "group_id": value.GroupId, "receiver": value.Receiver,
		"send_time": bson.M{"$gte": user.LogoutTime, "$lte": nowTime}}).Sort("-send_time").Limit(PullCount).Count()
	if err != nil {
		Logger.Error("get unread single msg failed: ", err)
		return
	}

	MsgCount[value.GroupId] = num

	//循环推送查找到的每一条未读消息
	for _, v := range messages {
		msgUser := &protocol.User{}
		msgUser.Username = &v.Sender

		group := &protocol.Groups{}
		group.Id = &v.GroupId

		send := &protocol.SendMessage{}
		msgType := protocol.Mode(v.MsgType)
		send.Type = &msgType
		send.Msg = &v.Msg
		sendtime := v.MsgTime.Unix()
		send.MsgTime = &sendtime
		returnMsg := &protocol.WMessage{
			MsgType: proto.String("groupMsgResult"),
			//MsgTypeId: proto.Int32(16),
			StataCode: proto.Int32(0),
			SendMsg:   send,
			UserInfo:  msgUser,
		}

		buffer, err := proto.Marshal(returnMsg)
		if err != nil {
			Logger.Warn("protobuf marshal failed: ", err)
			return
		}
		newbuffer := append(common.IntToBytes(len(buffer)), buffer...)
		switch conn.(type) {
		case net.Conn:
			_, err = conn.(net.Conn).Write(newbuffer)
			common.CheckErr(err)
		case websocket.Conn:
			_, err = conn.(*websocket.Conn).Write(newbuffer)
			common.CheckErr(err)
		}
	}
}
