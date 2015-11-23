/*
删除好友
*/
package controller

import (
	"aTalk/protocol"
	"aTalk/server/common"
	"gopkg.in/mgo.v2/bson"
	//"net"
)

func DelFriendController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("user")
	exist, err := collection.Find(bson.M{"user_name": msg.UserInfo.GetUsername(), "friends.name": msg.Friends[0].GetUsername(), "friends.mode": 3}).Count()
	if err != nil {
		common.CommonResultReturn(conn, "delFriendResult", msg.GetMsgTypeId(), 3)
		Logger.Error("get friend relation:", err)
		return
	}
	if exist == 0 {
		common.CommonResultReturn(conn, "delFriendResult", msg.GetMsgTypeId(), 2)
		return
	}

	//把好友从自己好友列表中删除
	err = collection.Update(bson.M{"user_name": msg.UserInfo.GetUsername()}, bson.M{"$pull": bson.M{"friends.name": msg.Friends[0].GetUsername()}})
	//把自己从对方好友列表中删除
	err = collection.Update(bson.M{"user_name": msg.Friends[0].GetUsername()}, bson.M{"$pull": bson.M{"friends.name": msg.UserInfo.GetUsername()}})
	if err != nil {
		common.CommonResultReturn(conn, "delFriendResult", msg.GetMsgTypeId(), 3)
		Logger.Error("get friend relation:", err)
		return
	}

	common.CommonResultReturn(conn, "delFriendResult", msg.GetMsgTypeId(), 0)
}
