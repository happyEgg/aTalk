package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id             bson.ObjectId `bson:"_id"`
	UserName       string        `bson:"user_name"`
	Password       string        `bson:"password"`
	NewPassword    string        `bson:"new_password"`
	RealName       string        `bson:"real_name"`
	Gender         int32         `bson:"gender"`
	Age            int32         `bson:"age"`
	Phone          string        `bson:"phone"`
	Email          string        `bson:"email"`
	ProfilePhotoid int32         `bson:"profile_photoid"`
	ChatGroundid   int32         `bson:"chat_groundid"`
	Notice         string        `bson:"notice"`
	Friends        []*Friend     `bson:"friends"`
	CreateTime     time.Time     `bson:"create_time"`
	LoginTime      time.Time     `bson:"login_time"`
	LogoutTime     time.Time     `bson:"logout_time"`
}

type Friend struct {
	Name   string `bson:"name"`   //好友登录名
	Mode   int32  `bson:"mode"`   //好友关系,3代表两者是好友关系
	Remark string `bson:"remark"` //备注
}
