/*
模糊查询搜索好友
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

func SearchUserController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("user")
	var user []models.User
	var searchUser = make([]*protocol.Friend, 0)

	friend := msg.Friends[0].GetUsername()

	//模糊查询用户,并返回全部查找到的好友信息
	err := collection.Find(bson.M{"user_name": bson.M{"$regex": friend, "$options": 'i'}}).
		Select(bson.M{"user_name": true, "gender": true, "profile_photoid": true}).Limit(10).All(&user)
	if err != nil {
		Logger.Error("select user info: ", err)
		common.CommonResultReturn(conn, "searchUserResult", msg.GetMsgTypeId(), 1)
		return
	}
	for _, v := range user {
		friend := &protocol.Friend{}
		friend.Username = &v.UserName
		sex := protocol.Gender(v.Gender)
		friend.Sex = &sex
		friend.ProfilePhotoid = &v.ProfilePhotoid
		searchUser = append(searchUser, friend)
	}

	ReturnMsg := &protocol.WMessage{
		MsgType:   proto.String("searchUserResult"),
		MsgTypeId: msg.MsgTypeId,
		StataCode: proto.Int32(0),
		Friends:   searchUser,
	}
	buffer, err := proto.Marshal(ReturnMsg)
	if err != nil {
		common.CommonResultReturn(conn, "searchUserResult", msg.GetMsgTypeId(), 1)
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
