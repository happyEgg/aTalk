/*
创建群组
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
)

//创建群组
func GroupCreateController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("group")
	collectionUser := common.DBMongo.C("user")
	var group models.Group
	var members = make([]*protocol.Friend, 0)

	group.Id = bson.NewObjectId()
	group.Owner = msg.UserInfo.GetUsername()
	group.GroupName = msg.Group.GetGroupName()
	group.Description = msg.Group.GetDescription()
	for k, v := range msg.Group.GetGroupMember() {
		group.Members[k].Name = v.GetUsername()
	}
	//把新创建的组加入到数据库
	err := collection.Insert(&group)
	if err != nil {
		Logger.Error("group create:", err)
		common.CommonResultReturn(conn, "groupCreateResult", msg.GetMsgTypeId(), 1)
		return
	}

	for _, v := range group.Members {
		member := &protocol.Friend{}
		var user models.User
		member.Username = &v.Name
		member.Remark = &v.Remark
		err := collectionUser.Find(bson.M{"user_name": v.Name}).Select(bson.M{"profile_photoid": true}).One(&user)
		if err != nil {
			Logger.Error("group create:", err)
			common.CommonResultReturn(conn, "groupCreateResult", msg.GetMsgTypeId(), 1)
			return
		}

		//好友头像
		member.ProfilePhotoid = &user.ProfilePhotoid
		members = append(members, member)
	}
	//把群信息返回给创建者
	returnMsg := &protocol.WMessage{
		MsgType:   proto.String("groupCreateResult"),
		MsgTypeId: msg.MsgTypeId,
		StataCode: proto.Int32(0),
		Group: &protocol.Groups{
			Id:          proto.String(group.Id.String()),
			GroupName:   proto.String(group.GroupName),
			Description: proto.String(group.Description),
			Owner:       proto.String(group.Owner),
			GroupMember: members,
		},
	}

	buffer, err := proto.Marshal(returnMsg)
	if err != nil {
		Logger.Warn("protobuf marshal failed: ", err)
		common.CommonResultReturn(conn, "groupCreateResult", msg.GetMsgTypeId(), 1)
		fmt.Println("failed: ", err)
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
