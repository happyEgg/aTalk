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
	URL = config.String("db_host")
)
