/*
修改个人信息，如果修改了密码
*/
package controller

import (
	"aTalk/protocol"
	"aTalk/server/common"
	"code.google.com/p/go.net/websocket"
	"gopkg.in/mgo.v2/bson"
	"net"
)

func ModifyInfoController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("user")
	arr := common.ProtoToData(msg)

	//判断敏感词
	if !SensitiveWords(msg.UserInfo.GetRealName()) {
		common.CommonResultReturn(conn, "modifyInfoResult", msg.GetMsgTypeId(), 2)
		return
	}

	//判断新密码是否为空,如果是空,证明更新的是其他内容。不为空,证明是要修改密码
	_, exist := arr["new_password"]
	if !exist {

		//如果新密码为空，证明不是修改密码
		err := collection.Update(bson.M{"user_name": msg.UserInfo.GetUsername()}, bson.M{"$set": arr})
		if err != nil {
			Logger.Error("update user info: ", err)
			common.CommonResultReturn(conn, "modifyInfoResult", msg.GetMsgTypeId(), 3)
		} else {
			common.CommonResultReturn(conn, "modifyInfoResult", msg.GetMsgTypeId(), 0)
		}
	} else {

		//验证此用户的旧密码是否正确
		n, err := collection.Find(bson.M{"user_name": msg.UserInfo.GetUsername(), "password": msg.UserInfo.GetPassword()}).Count()
		if err != nil {
			Logger.Error("old password: ", err)
			common.CommonResultReturn(conn, "modifyInfoResult", msg.GetMsgTypeId(), 3)
			return
		}

		//如果n=0证明旧密码不正确
		if n == 0 {
			//Logger.Error("update user password: ", err)
			common.CommonResultReturn(conn, "modifyInfoResult", msg.GetMsgTypeId(), 1)
		} else {
			//旧密码正确，修改密码
			pwd := arr["new_password"]
			err = collection.Update(bson.M{"user_name": msg.UserInfo.GetUsername()}, bson.M{"$set": bson.M{"password": pwd}})
			if err != nil {
				Logger.Error("old password: ", err)
				common.CommonResultReturn(conn, "modifyInfoResult", msg.GetMsgTypeId(), 3)
				return
			}
			common.CommonResultReturn(conn, "modifyInfoResult", msg.GetMsgTypeId(), 0)
			delete(UserMap, msg.UserInfo.GetUsername())
			delete(ConnMap, conn)
			switch conn.(type) {
			case net.Conn:
				conn.(net.Conn).Close()
			case websocket.Conn:
				conn.(*websocket.Conn).Close()
			}
		}
	}
}
