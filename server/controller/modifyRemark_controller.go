/*
修改好友备注
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
)

func ModifyRemarkController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("user")
	var user models.User
	var friends = make([]*protocol.Friend, 0)

	//在好友关系中查找是否有此好友,如果存在,就更新其备注
	//更新成功，把此好友的备注及其他内容返回
	friendName := msg.Friends[0].GetUsername()
	remark := msg.Friends[0].GetRemark()
	n, err := collection.Find(bson.M{"user_name": msg.UserInfo.GetUsername(), "friends.username": friendName}).Count()
	if err != nil {
		common.CommonResultReturn(conn, "modifyRemarkResult", msg.GetMsgTypeId(), 2)
		return
	}
	if n == 0 {
		common.CommonResultReturn(conn, "modifyRemarkResult", msg.GetMsgTypeId(), 1)
		return
	}
	err = collection.Find(bson.M{"user_name": friendName}).Select(bson.M{"profile_photoid": true}).One(&user)
	if err == nil {
		friendReturn := &protocol.Friend{}

		err := collection.Update(bson.M{"user_name": msg.UserInfo.GetUsername(), "friends.name": friendName}, bson.M{"friends.remark": remark})
		if err != nil {
			common.CommonResultReturn(conn, "modifyRemarkResult", msg.GetMsgTypeId(), 2)
			return
		}
		friendReturn.ProfilePhotoid = &user.ProfilePhotoid

		friendReturn.Username = &friendName

		friendReturn.Remark = &remark
		if OnlineCheck(friendName) {
			x := int32(1)
			friendReturn.Online = &x
		} else {
			x := int32(0)
			friendReturn.Online = &x
		}
		friends = append(friends, friendReturn)
		msgReturn := &protocol.WMessage{
			MsgType:   proto.String("modifyRemarkResult"),
			MsgTypeId: msg.MsgTypeId,
			StataCode: proto.Int32(0),
			Friends:   friends,
		}

		buffer, err := proto.Marshal(msgReturn)
		if err != nil {
			common.CommonResultReturn(conn, "modifyRemarkResult", msg.GetMsgTypeId(), 2)
			Logger.Warn("protobuf marshal failed: ", err)
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

	} else {
		common.CommonResultReturn(conn, "modifyRemarkResult", msg.GetMsgTypeId(), 2)
	}
}
