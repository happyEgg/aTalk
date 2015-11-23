/*
群成员备注修改
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

func GroupRemarkController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("group")
	var group models.Group

	remark := msg.Friends[0].GetRemark()
	//判断此成员在不在此群中
	n, err := collection.Find(bson.M{"_id": msg.Group.GetId(), "members.name": msg.UserInfo.GetUsername()}).Count()
	if err != nil || n == 0 {
		Logger.Error("group remark:", err)
		common.CommonResultReturn(conn, "groupRemarkResult", msg.GetMsgTypeId(), 1)
		return
	}

	//更新备注
	err = collection.Update(bson.M{"_id": msg.Group.GetId(), "members.name": msg.UserInfo.GetUsername()}, bson.M{"$set": bson.M{"members.remark": remark}})
	err = collection.Find(bson.M{"_id": msg.Group.GetId(), "members.name": msg.UserInfo.GetUsername()}).One(&group)
	if err != nil {
		Logger.Error("group remark:", err)
		common.CommonResultReturn(conn, "groupRemarkResult", msg.GetMsgTypeId(), 2)
		return
	}

	var friends = make([]*protocol.Friend, 0)
	for _, v := range group.Members {
		friend := &protocol.Friend{}
		friend.Username = &v.Name
		friend.Remark = &v.Remark
		friends = append(friends, friend)
	}

	returnMsg := &protocol.WMessage{
		MsgType:   proto.String("groupRemarkResult"),
		MsgTypeId: msg.MsgTypeId,
		StataCode: proto.Int32(0),
		Group: &protocol.Groups{
			Id:          proto.String(group.Id.String()),
			GroupName:   proto.String(group.GroupName),
			Description: proto.String(group.Description),
			Owner:       proto.String(group.Owner),
			GroupMember: friends,
		},
	}

	buffer, err := proto.Marshal(returnMsg)
	if err != nil {
		common.CommonResultReturn(conn, "groupRemarkResult", msg.GetMsgTypeId(), 3)
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
	return
}
