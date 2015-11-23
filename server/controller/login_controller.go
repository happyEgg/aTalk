/*
登陆请求，成功后返回个人信息，同时推送未读单聊消息和群聊信息
*/
package controller

import (
	"aTalk/protocol"
	"aTalk/server/common"
	"aTalk/server/models"
	//"fmt"
	"code.google.com/p/go.net/websocket"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2/bson"
	"net"
	"time"
)

func LoginController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("user")
	collectionMsg := common.DBMongo.C("message")
	collectionDevice := common.DBMongo.C("device")
	var user models.User
	var loginMsg models.Device

	//判断此用户是否在其他地方登陆
	if OnlineCheck(msg.UserInfo.GetUsername()) {
		common.CommonResultReturn(UserMap[msg.UserInfo.GetUsername()], "loginResult", 0, 5)
		delete(ConnMap, UserMap[msg.UserInfo.GetUsername()])
		delete(UserMap, msg.UserInfo.GetUsername())
		return
	}

	//判断用户是否存在
	n, err := collection.Find(bson.M{"user_name": msg.UserInfo.GetUsername()}).Count()
	if err != nil {
		common.CommonResultReturn(conn, "loginResult", msg.GetMsgTypeId(), 2)
		Logger.Error("select user:", err)
		return
	}

	if n != 0 {

		//判断账号密码是否正确
		err = collection.Find(bson.M{"user_name": msg.UserInfo.GetUsername(), "password": msg.UserInfo.GetPassword()}).One(&user)
		if err != nil {
			Logger.Error("select user:", err)
			common.CommonResultReturn(conn, "loginResult", msg.GetMsgTypeId(), 1)
			return
		}
		//登陆成功返回个人信息
		singleInfo(conn, msg, &user)

		//记录此次设备登陆的信息
		nowTime := time.Now()
		loginMsg.Id = bson.NewObjectId()
		loginMsg.UserName = msg.UserInfo.GetUsername()
		loginMsg.Devices = msg.GetDevices()
		loginMsg.LoginTime = nowTime
		if collectionDevice.Insert(&loginMsg) != nil {
			Logger.Error("insert loginMsg failed:", err)
		} else {
			LoginTimeId = loginMsg.Id
		}

		//登陆成功，更新此用户的登陆时间
		err = collection.Update(bson.M{"user_name": msg.UserInfo.GetUsername()}, bson.M{"$set": bson.M{"login_time": nowTime}})
		if err != nil {
			Logger.Error(" login:", err)
		}

		//把用户和链接加入到map中
		UserMap[msg.UserInfo.GetUsername()] = conn
		ConnMap[conn] = msg.UserInfo.GetUsername()

		//用户上线后推送未读的聊天记录

		//查找有多少用户给我的消息未读
		var result []string
		err = collectionMsg.Find(bson.M{"receiver": user.UserName, "group_id": bson.M{"$exists": false},
			"send_time": bson.M{"$gt": user.LogoutTime, "$lte": nowTime}}).Distinct("sender", &result)
		if err != nil {
			Logger.Error("get user member failed: ", err)
			return
		}

		for _, value := range result {

			//开启线程发送未读消息
			go HandleSingleUnreadMsg(conn, &value, &user)
		}

		//开启一个go程返回未读群消息
		go groupUnreadMsg(conn, &user)

	} else {
		common.CommonResultReturn(conn, "loginResult", msg.GetMsgTypeId(), 2)
	}
}

func groupUnreadMsg(conn interface{}, user *models.User) {
	var result []models.Message
	collection := common.DBMongo.C("message")

	nowTime := time.Now()

	//查找此用户所属的群组
	err := collection.Find(bson.M{"receiver": user.UserName, "group_id": bson.M{"$exists": true},
		"send_time": bson.M{"$gte": user.LogoutTime, "$lte": nowTime}}).Distinct("group_id", &result)
	if err != nil {
		Logger.Error("push msg failed: ", err)
	}

	for _, value := range result {

		//发送每一个群的信息
		go HandleGroupUnreadMsg(conn, &value, user)
	}
}

func singleInfo(conn interface{}, msg *protocol.WMessage, user *models.User) {
	//collection := common.DBMongo.C("user")
	returnUser := &protocol.User{}
	returnUser.Username = &user.UserName
	sex := protocol.Gender(user.Gender)
	returnUser.Sex = &sex
	returnUser.Age = &user.Age
	returnUser.Phone = &user.Phone
	returnUser.RealName = &user.RealName
	returnUser.ProfilePhotoid = &user.ProfilePhotoid
	returnUser.ChatGroundid = &user.ChatGroundid
	returnUser.Email = &user.Email
	returnUser.Notice = &user.Notice

	returnMsg := &protocol.WMessage{
		MsgType:   proto.String("loginResult"),
		MsgTypeId: msg.MsgTypeId,
		StataCode: proto.Int32(0),
		UserInfo:  returnUser,
		MsgTime:   proto.Int64(time.Now().Unix()),
	}

	buffer, err := proto.Marshal(returnMsg)
	if err != nil {
		common.CommonResultReturn(conn, "loginResult", msg.GetMsgTypeId(), 3)
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
