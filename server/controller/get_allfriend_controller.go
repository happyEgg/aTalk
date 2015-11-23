/*
获取所有好友
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

//获取所有好友
func GetAllFriendController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("user")
	var user models.User
	var friends = make([]*protocol.Friend, 0)

	//获取和用户相关联的所有好友
	err := collection.Find(bson.M{"user_name": msg.UserInfo.GetUsername()}).Select(bson.M{"friends": true, "_id": false}).One(&user)
	if err != nil {
		Logger.Error("select friends: ", err)
		common.CommonResultReturn(conn, "allFriendsResult", msg.GetMsgTypeId(), 2)
		return
	}

	//把好友循环写入好友数组中发送给请求者
	for _, v := range user.Friends {
		friend := &protocol.Friend{}
		var user models.User

		friend.Username = &v.Name
		friend.Remark = &v.Remark
		err := collection.Find(bson.M{"user_name": v.Name}).Select(bson.M{"profile_photoid": true, "_id": false}).One(&user)
		if err != nil {
			Logger.Error("select friend failed: ", err)
			common.CommonResultReturn(conn, "allFriendsResult", msg.GetMsgTypeId(), 3)
			return
		}

		//好友头像
		friend.ProfilePhotoid = &user.ProfilePhotoid

		//判断该好友是否在线
		if OnlineCheck(v.Name) {
			x := int32(1)
			friend.Online = &x
		} else {
			x := int32(0)
			friend.Online = &x
		}

		friends = append(friends, friend)
	}

	ReturnMsg := &protocol.WMessage{
		MsgType:   proto.String("allFriendsResult"),
		MsgTypeId: msg.MsgTypeId,
		StataCode: proto.Int32(0),
		Friends:   friends,
	}
	buffer, err := proto.Marshal(ReturnMsg)
	if err != nil {
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
}
