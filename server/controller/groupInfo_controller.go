/*
获取组信息和组成员
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

func GetGroupInfoController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("group")
	collectionUser := common.DBMongo.C("user")
	var group models.Group
	var members = make([]*protocol.Friend, 0)
	err := collection.Find(bson.M{"_id": msg.Group.GetId()}).All(&group)
	if err != nil {
		Logger.Error("group info:", err)
		common.CommonResultReturn(conn, "groupInfoResult", msg.GetMsgTypeId(), 1)
		return
	}

	for _, v := range group.Members {
		member := &protocol.Friend{}
		var user models.User
		member.Username = &v.Name
		member.Remark = &v.Remark
		err := collectionUser.Find(bson.M{"user_name": v.Name}).Select(bson.M{"profile_photoid": true}).One(&user)
		if err != nil {
			Logger.Error("group info: ", err)
			common.CommonResultReturn(conn, "groupInfoResult", msg.GetMsgTypeId(), 1)
			return
		}
		//好友头像
		member.ProfilePhotoid = &user.ProfilePhotoid
		members = append(members, member)
	}
	returnMsg := &protocol.WMessage{
		MsgType:   proto.String("groupInfoResult"),
		MsgTypeId: msg.MsgTypeId,
		StataCode: proto.Int32(0),
		Group: &protocol.Groups{
			Id:          proto.String(group.Id.String()),
			GroupName:   proto.String(group.GroupName),
			Owner:       proto.String(group.Owner),
			Description: proto.String(group.Description),
			GroupMember: members,
		},
	}
	buffer, err := proto.Marshal(returnMsg)
	if err != nil {
		Logger.Error("get friendinfo:", err)
		common.CommonResultReturn(conn, "groupInfoResult", msg.GetMsgTypeId(), 1)
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
	return
}
