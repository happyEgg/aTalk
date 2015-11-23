/*
查看好友详细信息
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

func GetFriendInfoController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("user")
	var user models.User
	n, err := collection.Find(bson.M{"user_name": ConnMap[conn], "friends.name": msg.UserInfo.GetUsername()}).Count()
	if err != nil {
		Logger.Error("get friendinfo:", err)
		common.CommonResultReturn(conn, "friendInfoResult", msg.GetMsgTypeId(), 2)
		return
	}
	if n == 0 {
		common.CommonResultReturn(conn, "friendInfoResult", msg.GetMsgTypeId(), 1)
		return
	}
	err = collection.Find(bson.M{"user_name": msg.UserInfo.GetUsername()}).One(&user)
	if err != nil {
		Logger.Error("get friendinfo:", err)
		common.CommonResultReturn(conn, "friendInfoResult", msg.GetMsgTypeId(), 2)
		return
	}

	userInfo := &protocol.User{}
	userInfo = common.DataToProto(&user)
	returnMsg := &protocol.WMessage{
		MsgType:   proto.String("friendInfoResult"),
		MsgTypeId: proto.Int32(msg.GetMsgTypeId()),
		StataCode: proto.Int32(0),
		UserInfo:  userInfo,
	}

	buffer, err := proto.Marshal(returnMsg)
	if err != nil {
		Logger.Error("get friendinfo:", err)
		common.CommonResultReturn(conn, "friendInfoResult", msg.GetMsgTypeId(), 2)
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
