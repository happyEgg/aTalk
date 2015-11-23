/*
管理员或成员退出群的操作
*/
package controller

import (
	"aTalk/protocol"
	"aTalk/server/common"
	"aTalk/server/models"
	"gopkg.in/mgo.v2/bson"
	//"net"
)

//退出群
func GroupExitController(conn interface{}, msg *protocol.WMessage) {
	var group models.Group
	collection := common.DBMongo.C("group")

	//得到群信息和群成员
	err := collection.Find(bson.M{"_id": msg.Group.GetId()}).One(&group)
	if err != nil {
		Logger.Error("select group id failed: ", err)
		common.CommonResultReturn(conn, "groupExitResult", msg.GetMsgTypeId(), 1)
		return
	}

	if group.Owner == msg.UserInfo.GetUsername() {
		//证明是群主退群，删除群
		err = collection.Remove(bson.M{"_id": msg.Group.GetId()})
		if err != nil {
			Logger.Error("delete group member: ", err)
			common.CommonResultReturn(conn, "groupExitResult", msg.GetMsgTypeId(), 1)
			return
		}

		//系统通知所有人函数
		GroupSystemMsg(msg, &group, 0)
	} else {
		//证明是普通用户退群,将其从群成员中删除
		err = collection.Update(bson.M{"_id": msg.Group.GetId()}, bson.M{"$pull": bson.M{"friends": bson.M{"name": msg.UserInfo.GetUsername()}}})
		if err != nil {
			Logger.Error("delete group member: ", err)
			common.CommonResultReturn(conn, "groupExitResult", msg.GetMsgTypeId(), 1)
			return
		}
		//系统通知所有人函数
		GroupSystemMsg(msg, &group, 1)
	}

	common.CommonResultReturn(conn, "groupExitResult", msg.GetMsgTypeId(), 0)
}

//群主踢人函数
func GroupKickController(conn interface{}, msg *protocol.WMessage) {
	var group models.Group
	collection := common.DBMongo.C("group")

	//管理员踢人,先判断此用户是不是此群的管理员
	err := collection.Find(bson.M{"_id": msg.Group.GetId()}).One(&group)
	if err != nil {
		common.CommonResultReturn(conn, "groupKickResult", msg.GetMsgTypeId(), 2)
		return
	}

	//如果为真，证明其是管理员
	if group.Owner == msg.UserInfo.GetUsername() {

		//将其从群成员中删除
		err = collection.Update(bson.M{"_id": msg.Group.GetId()}, bson.M{"$pull": bson.M{"friends": bson.M{"name": msg.Group.GroupMember[0].GetUsername()}}})
		if err != nil {
			common.CommonResultReturn(conn, "groupKickResult", msg.GetMsgTypeId(), 2)
			return
		}
		//系统通知所有人函数
		GroupSystemMsg(msg, &group, 2)
	} else {
		common.CommonResultReturn(conn, "groupKickResult", msg.GetMsgTypeId(), 1)
	}
}
