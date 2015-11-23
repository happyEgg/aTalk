/*
系统提示消息发送
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

func GroupSystemMsg(msg *protocol.WMessage, group *models.Group, num int) {
	collection := common.DBMongo.C("group")

	//获取所有在这个组的成员,并发送系统提示消息
	for _, v := range group.Members {
		var message models.Message

		//把系统消息加入到数据库
		message.Id = bson.NewObjectId()
		message.GroupId = msg.Group.GetId()
		message.Receiver = v.Name
		message.MsgType = 0

		//遍历到此用户时，发送回执，不发送下面的消息
		if v.Name == msg.UserInfo.GetUsername() && num != 2 {
			common.CommonResultReturn(UserMap[msg.UserInfo.GetUsername()], "groupExitResult", msg.GetMsgTypeId(), 0)
			continue
		}
		//如果num = 0，证明是管理员解散群
		if num == 0 {
			message.Msg = msg.UserInfo.GetUsername() + "管理员解散该群"
		} else if num == 1 {
			message.Msg = msg.UserInfo.GetUsername() + "退出该群"

		} else if num == 2 {
			//如果num=2，管理员踢人通知
			if v.Name == msg.UserInfo.GetUsername() {
				common.CommonResultReturn(UserMap[msg.UserInfo.GetUsername()], "groupKickResult", msg.GetMsgTypeId(), 0)
				continue
			}
			message.Msg = msg.Friends[0].GetUsername() + "被管理员请出群"
		}

		nowTime := time.Now()
		message.MsgTime = nowTime
		if err := collection.Insert(&message); err != nil {
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

			if num == 2 {
				returnMsg := &protocol.WMessage{
					MsgType: proto.String("groupKickResult"),
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
			} else {

				returnMsg := &protocol.WMessage{
					MsgType: proto.String("groupExitResult"),
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
}
