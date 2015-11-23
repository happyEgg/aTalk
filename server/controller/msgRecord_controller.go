/*
单聊信息和群聊信息历史记录请求，每次返回特定的数量
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
	"time"
)

func SingleMsgRecordController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("message")
	collectionUser := common.DBMongo.C("user")
	var messages []models.Message
	var user models.User
	var oldTime time.Time

	//返回此用户的登陆时间
	err := collectionUser.Find(bson.M{"user_name": msg.UserInfo.GetUsername()}).Select(bson.M{"login_time": true}).One(&user)
	if err != nil {
		Logger.Error("single msg record:", err)
		common.CommonResultReturn(conn, "msgRecordResult", msg.GetMsgTypeId(), 1)
		return
	}

	now := time.Now().Unix()
	//只返回3天以内的聊天记录
	limitTime := now - TimeLimit*24*3600
	if msg.GetMsgTime() < limitTime {
		oldTime = time.Unix(limitTime, 0)
	} else {
		oldTime = time.Unix(msg.GetMsgTime(), 0)
	}

	//用来接收发发送人给自己的和自己发送给发送人的消息
	var cond = make([]interface{}, 0)
	cond1 := &models.Message{
		Sender:   msg.UserInfo.GetUsername(),
		Receiver: msg.Friends[0].GetUsername(),
	}
	cond = append(cond, cond1)
	cond2 := &models.Message{
		Sender:   msg.Friends[0].GetUsername(),
		Receiver: msg.UserInfo.GetUsername(),
	}
	cond = append(cond, cond2)
	err = collection.Find(bson.M{"$or": cond, "send_time": bson.M{"$gte": oldTime, "lte": user.LoginTime}}).
		Sort("-send_time").Limit(OnceRequent).Skip(int(MsgCount[msg.UserInfo.GetUsername()])).All(&messages)
	if err != nil {
		Logger.Error("msgRecord failed:", err)
		common.CommonResultReturn(conn, "msgRecordResult", msg.GetMsgTypeId(), 1)
		return
	}

	//记录每次需要偏移的值
	MsgCount[msg.UserInfo.GetUsername()] = MsgCount[msg.UserInfo.GetUsername()] + OnceRequent

	//循环把数据发送出去
	for _, v := range messages {
		msgUser := &protocol.User{}
		send := &protocol.SendMessage{}
		msgUser.Username = &v.Sender
		send.Receiver = &v.Receiver
		msgType := protocol.Mode(v.MsgType)
		send.Type = &msgType
		send.Msg = &v.Msg
		sendtime := v.MsgTime.Unix()
		send.MsgTime = &sendtime
		returnMsg := &protocol.WMessage{
			MsgType:   proto.String("msgRecordResult"),
			MsgTypeId: msg.MsgTypeId,
			StataCode: proto.Int32(0),
			SendMsg:   send,
			UserInfo:  msgUser,
		}

		buffer, err := proto.Marshal(returnMsg)
		if err != nil {
			common.CommonResultReturn(conn, "msgRecordResult", msg.GetMsgTypeId(), 1)
			Logger.Warn("protobuf marshal failed: ", err)
			continue
		}
		//fmt.Println("msgRecord: ", returnMsg.String())
		newbuffer := append(common.IntToBytes(len(buffer)), buffer...)
		switch conn.(type) {
		case net.Conn:
			_, err = conn.(net.Conn).Write(newbuffer)
			common.CheckErr(err)
		case websocket.Conn:
			_, err = conn.(*websocket.Conn).Write(newbuffer)
			common.CheckErr(err)
		}
	}
}

func GroupMsgRecordController(conn interface{}, msg *protocol.WMessage) {
	collection := common.DBMongo.C("message")
	var messages []models.Message
	var oldTime time.Time

	now := time.Now().Unix()

	//只返回3天以内的聊天记录
	if msg.GetMsgTime() < now-(TimeLimit*24*3600) {
		oldTime = time.Unix(now-(TimeLimit*24*3600), 0)
	} else {
		oldTime = time.Unix(msg.GetMsgTime(), 0)
	}

	// 每次获取一定条数的群消息
	err := collection.Find(bson.M{"group_id": msg.Group.GetId(), "send_time": bson.M{"gte": oldTime}}).
		Sort("-send_time").Limit(OnceRequent).Skip(int(MsgCount[msg.Group.GetId()])).All(&messages)
	if err != nil {
		Logger.Error("group msgRecord failed:", err)
		common.CommonResultReturn(conn, "groupMsgRecordResult", msg.GetMsgTypeId(), 1)
		return
	}

	//统计需要的偏移量
	MsgCount[msg.Group.GetId()] = MsgCount[msg.Group.GetId()] + OnceRequent

	for _, v := range messages {
		msgUser := &protocol.User{}
		send := &protocol.SendMessage{}
		msgUser.Username = &v.Sender
		send.Receiver = &v.Receiver
		msgType := protocol.Mode(v.MsgType)
		send.Type = &msgType
		send.Msg = &v.Msg
		sendTime := v.MsgTime.Unix()
		send.MsgTime = &sendTime
		returnMsg := &protocol.WMessage{
			MsgType:   proto.String("groupMsgRecordResult"),
			MsgTypeId: msg.MsgTypeId,
			StataCode: proto.Int32(0),
			SendMsg:   send,
			UserInfo:  msgUser,
		}

		buffer, err := proto.Marshal(returnMsg)
		if err != nil {
			Logger.Warn("protobuf marshal failed: ", err)
			continue
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
	}
}
