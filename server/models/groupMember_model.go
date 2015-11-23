package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Group struct {
	Id          bson.ObjectId `bson:"_id"`
	GroupName   string        `bson:"group_name"`
	Owner       string        `bson:"owner"`
	Description string        `bson:"description"`
	Members     []*Member     `bson:"members"`
	//MemberName  string `orm:"size(32)"`
}

type Member struct {
	Remark string `bson:"remark"`
	Name   string `bson:"name"`
}
