/*
用户注册
*/
package controller

import (
	"aTalk/protocol"
	"aTalk/server/common"
	"aTalk/server/models"
	"gopkg.in/mgo.v2/bson"
	//"net"
	"time"
)

func RegisterController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("user")
	var user models.User

	//判断敏感词
	if !SensitiveWords(msg.UserInfo.GetUsername()) {

		common.CommonResultReturn(conn, "registerResult", msg.GetMsgTypeId(), 2)
		return
	}

	//判断注册名是否存在，不存在就在数据库里创建
	user.UserName = msg.UserInfo.GetUsername()
	n, err := collection.Find(bson.M{"user_name": user.UserName}).Count()
	if err != nil {
		common.CommonResultReturn(conn, "registerResult", msg.GetMsgTypeId(), 3)
		return
	}

	//证明此用户已存在
	if n != 0 {
		common.CommonResultReturn(conn, "registerResult", msg.GetMsgTypeId(), 1)
		return
	}

	user.Id = bson.NewObjectId()
	user.Password = msg.UserInfo.GetPassword()
	user.CreateTime = time.Now()
	if collection.Insert(&user) == nil {
		common.CommonResultReturn(conn, "registerResult", msg.GetMsgTypeId(), 0)
	} else {
		common.CommonResultReturn(conn, "registerResult", msg.GetMsgTypeId(), 3)
	}
}
