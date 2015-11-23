/*配置日志文件*/

package config

import (
	"github.com/astaxie/beego/logs"
)

var Logger *logs.BeeLogger

func init() {
	Logger = logs.NewLogger(10000)
	Logger.EnableFuncCallDepth(true)
	Logger.SetLogger("file", `{"filename":"diary/diary.log"}`)
}
