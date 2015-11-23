/*
邀请好友加入群
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

func GroupInviteController(conn interface{}, msg *protocol.WMessage) {

	collection := common.DBMongo.C("group")
	var group models.Group
	var arr = make([]string, 0)

	//判断被邀请人是否已在数据表中
	//name := msg.Group.GroupMember[0].GetUsername()
	// for _, v := msg.Group.GetGroupMember(){

	// }

	err := collection.Find(bson.M{"_id": msg.Group.GetId()}).All(&group)
	if err != nil {
		Logger.Error("insert group member: ", err)
		common.CommonResultReturn(conn, "groupInviteResult", msg.GetMsgTypeId(), 2)
		return
	}

	for _, v := range msg.Group.GetGroupMember() {
		for _, value := range group.Members {
			if v.GetUsername() == value.Name {
				continue
			}
		}
		arr = append(arr, v.GetUsername())
	}

	//把被邀请人加入到组关系数据库中
	err = collection.Update(bson.M{"_id": msg.Group.GetId()}, bson.M{"$push": bson.M{"members.name": arr}})
	if err != nil {
		Logger.Error("insert group member: ", err)
		common.CommonResultReturn(conn, "groupInviteResult", msg.GetMsgTypeId(), 2)
		return
	}
	common.CommonResultReturn(conn, "groupInviteResult", msg.GetMsgTypeId(), 0)
	for _, v := range arr {
		var message models.Message
		//把系统提示消息加入到组聊信息数据库中
		message.Id = bson.NewObjectId()
		message.Receiver = v
		message.MsgType = 0
		message.GroupId = msg.Group.GetId()
		message.Msg = msg.UserInfo.GetUsername() + "邀请你加入群聊"
		nowTime := time.Now()
		message.MsgTime = nowTime
		err = collection.Insert(&message)
		if err != nil {
			Logger.Error("insert group message: ", err)
			//return
		}

		//如果被邀请人在线,系统消息发送
		if OnlineCheck(message.Receiver) {
			sendMessage := &protocol.SendMessage{}
			group := &protocol.Groups{}

			group.Id = &message.GroupId
			sendMessage.Receiver = &message.Receiver
			typeSystem := protocol.Mode_SYSTEM
			sendMessage.Type = &typeSystem

			sendMessage.Msg = &message.Msg
			nowtime := nowTime.Unix()
			sendMessage.MsgTime = &nowtime

			returnMsg := &protocol.WMessage{
				MsgType: proto.String("groupMsgResult"),
				//MsgTypeId: proto.Int32(13),
				StataCode: proto.Int32(0),
				Group:     group,
				SendMsg:   sendMessage,
			}

			buffer, err := proto.Marshal(returnMsg)
			if err != nil {
				common.CommonResultReturn(conn, "groupInviteResult", msg.GetMsgTypeId(), 2)
				Logger.Warn("protobuf marshal failed: ", err)
				return
			}
			newbuffer := append(common.IntToBytes(len(buffer)), buffer...)
			switch UserMap[message.Receiver].(type) {
			case net.Conn:
				_, err = (UserMap[message.Receiver]).(net.Conn).Write(newbuffer)
				common.CheckErr(err)
			case websocket.Conn:
				_, err = (UserMap[message.Receiver]).(*websocket.Conn).Write(newbuffer)
				common.CheckErr(err)
			}

		}
	}

	//获取所有在这个组的成员,并发送系统提示消息
	var groupMember models.Group
	err = collection.Find(bson.M{"_id": msg.Group.GetId()}).One(&groupMember)
	if err != nil {
		common.CommonResultReturn(conn, "groupInviteResult", msg.GetMsgTypeId(), 2)
		Logger.Error("insert group member: ", err)
		return
	}

	for _, v := range groupMember.Members {
		var message models.Message

		//把系统消息加入到数据库
		message.Id = bson.NewObjectId()
		message.GroupId = msg.Group.GetId()
		message.Receiver = v.Name
		message.MsgType = 0
		str1, str3 := "", ""
		if len(arr) > 2 {
			str3 = "等加入群聊"
		} else {
			str3 = "加入群聊"
		}
		if len(arr) == 1 {
			str1 = msg.UserInfo.GetUsername() + "邀请" + arr[0] + str3
		} else {
			str1 = msg.UserInfo.GetUsername() + "邀请" + arr[0] + "，" + arr[1] + str3
		}

		message.Msg = str1
		nowTime := time.Now()
		message.MsgTime = nowTime
		if collection.Insert(&message) != nil {
			Logger.Error("insert group message: ", err)
			continue
		}

		//如果此成员在线，给其发送系统消息
		if OnlineCheck(v.Name) {
			sendMsg := &protocol.SendMessage{}
			sendMsg.Receiver = &v.Name
			typeSystem := protocol.Mode_SYSTEM
			sendMsg.Type = &typeSystem
			sendMsg.Msg = &message.Msg
			msgTime := nowTime.Unix()
			sendMsg.MsgTime = &msgTime

			returnMsg := &protocol.WMessage{
				MsgType: proto.String("groupMsgResult"),
				//MsgTypeId: proto.Int32(13),
				StataCode: proto.Int32(0),
				Group: &protocol.Groups{
					Id: proto.String(msg.Group.GetId()),
				},

				SendMsg: sendMsg,
			}

			buffer, err := proto.Marshal(returnMsg)
			if err != nil {
				Logger.Warn("protobuf marshal failed: ", err)
				return
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
