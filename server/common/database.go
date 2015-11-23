/*
初始化mongodb数据库，同时初始化日志文件
*/

package common

import (
	"aTalk/server/config"
	"github.com/astaxie/beego/logs"
	"gopkg.in/mgo.v2"
)

//声明
var (
	URL     = config.String("db_host")
	db_name = config.String("db_name")
	DBMongo *mgo.Database
	Logger  *logs.BeeLogger
)

//配置日志文件
func initLog() {
	Logger = logs.NewLogger(10000)
	Logger.EnableFuncCallDepth(true)
	Logger.SetLogger("file", `{"filename":"diary/diary.log"}`)
}

//初始化数据库,并添加索引
func InitDB() {
	initLog()
	session, err := mgo.Dial(URL) //连接数据库
	CheckErr(err)
	//defer session.Close() //不应该关闭数据库
	session.SetMode(mgo.Monotonic, true)

	DBMongo = session.DB(db_name) //创建数据库名

	//
	DBMongo.C("message").EnsureIndex(mgo.Index{
		Key:        []string{"send_time"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	DBMongo.C("user").EnsureIndex(mgo.Index{
		Key:        []string{"user_name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	})
	DBMongo.C("friend").EnsureIndex(mgo.Index{
		Key:        []string{"name"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	//	DBMongo.C("group")
	DBMongo.C("member").EnsureIndex(mgo.Index{
		Key:        []string{"name"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	DBMongo.C("device").EnsureIndex(mgo.Index{
		Key:        []string{"user_name", "devices"},
		Unique:     false,
		DropDups:   false,
		Background: true,
	})
}
