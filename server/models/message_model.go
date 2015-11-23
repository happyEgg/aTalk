package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

//信息表，包括单聊和群聊
type Message struct {
	Id       bson.ObjectId `bson:"_id"`
	Sender   string        `bson:"sender"`
	GroupId  string        `bson:"group_id,omitempty"`
	Receiver string        `bson:"receiver"` //用来对系统通知消息进行处理
	MsgType  int32         `bson:"msg_type"`
	Msg      string        `bson:"msg"`
	MsgTime  time.Time     `bson:"send_time"`
}

//设备表
type Device struct {
	Id         bson.ObjectId `bson:"_id"`
	UserName   string        `bson:"user_name"`
	Devices    string        `bson:"devices"`
	LoginTime  time.Time     `bson:"login_time"`
	LogoutTime time.Time     `bson:"logout_time"`
}
