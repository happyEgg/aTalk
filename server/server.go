package main

//服务器端
import (
	"aTalk/server/common"
	"aTalk/server/config"
	"aTalk/server/controller"
	"runtime"
)

var PORT = config.String("port")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	common.InitDB()
	controller.Start(":" + PORT)

}
