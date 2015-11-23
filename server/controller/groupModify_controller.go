/*
修改群信息
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

//
func GroupModifyController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("group")
	var group models.Group
	var members = make([]*protocol.Friend, 0)
	arr := make(map[string]interface{})

	//判断其是否是群主
	n, err := collection.Find(bson.M{"_id": msg.Group.GetId(), "owner": msg.UserInfo.GetUsername()}).Count()
	if err != nil {
		common.CommonResultReturn(conn, "groupModifyResult", msg.GetMsgTypeId(), 1)
		Logger.Error("select group owner: ", err)
		return
	}

	if n != 0 {
		if msg.Group.GroupName != nil {
			arr["group_name"] = msg.Group.GetGroupName()
		}

		if msg.Group.Description != nil {
			arr["description"] = msg.Group.GetDescription()
		}

		err = collection.Update(bson.M{"_id": msg.Group.GetId()}, bson.M{"$set": arr})
		err = collection.Find(bson.M{"_id": msg.Group.GetId()}).One(&group)
		if err != nil {
			common.CommonResultReturn(conn, "groupModifyResult", msg.GetMsgTypeId(), 2)
			Logger.Error("group info:", err)
			return
		}

		//把群信息及成员信息更新发送给用户
		for _, v := range group.Members {
			member := &protocol.Friend{}
			var user models.User
			member.Username = &v.Name
			member.Remark = &v.Remark

			//获取成员头像
			err := collection.Find(bson.M{"user_name": v.Name}).Select(bson.M{"profile_photoid": true, "_id": false}).One(&user)
			if err != nil {
				common.CommonResultReturn(conn, "groupModifyResult", msg.GetMsgTypeId(), 2)
				Logger.Error("select friend failed: ", err)
				return
			}
			member.ProfilePhotoid = &user.ProfilePhotoid

			members = append(members, member)
		}

		proGroup := &protocol.Groups{}
		num := group.Id.String()
		proGroup.Id = &num
		proGroup.GroupName = &group.GroupName
		proGroup.Owner = &group.Owner
		proGroup.GroupMember = members
		proGroup.Description = &group.Description

		returnMsg := &protocol.WMessage{
			MsgType:   proto.String("groupModifyResult"),
			MsgTypeId: msg.MsgTypeId,
			StataCode: proto.Int32(0),
			Group:     proGroup,
		}

		buffer, err := proto.Marshal(returnMsg)
		if err != nil {
			common.CommonResultReturn(conn, "groupModifyResult", msg.GetMsgTypeId(), 2)
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
		common.CommonResultReturn(conn, "groupModifyResult", msg.GetMsgTypeId(), 1)
	}
}
