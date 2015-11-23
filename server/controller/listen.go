/*
监听，处理逻辑请求
*/
package controller

import (
	"aTalk/protocol"
	"aTalk/server/common"
	"aTalk/server/config"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
	"os"
	"time"
)

var (
	UserMap = make(map[string]interface{})
	//	WebUserMap = make(map[string]websocket.Conn)
	ConnMap = make(map[interface{}]string)
	// WebConnMap = make(map[interface{}]string)
	MsgCount = make(map[interface{}]int)
	//	FileMark    = make(chan int)
	LoginTimeId bson.ObjectId
	FileMark    int32
	FileSytle   int32
)

var (
	HeadLen     = config.Int("headLen")     //4字节包头
	PullCount   = config.Int("pullCount")   //最多返回多少条数据
	OnceRequent = config.Int("onceRequent") //单次请求返回条数
	TimeLimit   = config.Int64("timeLimit") //天

)

//初始化log
var Logger *logs.BeeLogger

func initLog() {
	Logger = logs.NewLogger(10000)
	Logger.EnableFuncCallDepth(true)
	Logger.SetLogger("file", `{"filename":"diary/diary.log"}`)
}

func Start(addr string) error {
	initLog()

	//websocket用户连接入口
	go Httpserver()

	log.Printf(">>>>>>>>端口:%s<<<<<<<<<<<<<<", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		//Log.SetLogger("file", `{"filename": "savelog.log"}`)
		os.Exit(-1)
	}
	defer listener.Close()

	log.Println(">>>>>>>>>>服务器启动<<<<<<<<<<")
	log.Println(">>>>>>>>等待客户端连接<<<<<<<<<")

	//死循环等待客户端连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			Logger.Error("accept: ", err)
			continue
		}

		log.Println(conn.RemoteAddr().String(), " 连接成功!")
		//开一个goroutines处理客户端消息
		go handleConnection(conn)
	}
}

//处理连接的用户操作
func handleConnection(conn net.Conn) {

	defer conn.Close()
	headBuffer := make([]byte, 4)

	//循环处理接收数据
	for {
		bufSize, err := conn.Read(headBuffer)
		if err != nil {
			//Logger.Error("conn.read:", err)
			break
		}
		if bufSize < 4 {
			continue
		}

		messager := make(chan byte)
		//心跳计时
		go common.HeartBeating(conn, messager, 20)
		//检测是client否有数据传来
		go common.GravelChannel(headBuffer, messager)

		bodyLen := common.BytesToInt(headBuffer)
		if bodyLen == 0 {
			common.KeepAliveResult(conn)
			continue
		}
		bodyBuffer := make([]byte, bodyLen)
		bodySize, err := conn.Read(bodyBuffer)
		if err != nil {
			break
		}

		//如果接收的包大小小于包头内容的长度，舍弃
		if bodySize < bodyLen {
			continue
		}

		msg := &protocol.WMessage{}
		err = proto.Unmarshal(bodyBuffer, msg)
		fmt.Println("msg:", msg.String())
		if err != nil {
			Logger.Warn("protobuf解包：", err)
			continue
		}

		//根据包的类型来进行处理
		switch msg.GetMsgType() {
		case "register":
			RegisterController(conn, msg)
			goto Break
		//只允许一方登陆
		case "login":
			LoginController(conn, msg)

		case "logout":

			fmt.Println(msg.UserInfo.GetUsername(), " logout")
			goto Break
		case "modifyInfo":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go ModifyInfoController(conn, msg)
			}

		case "sendFriendRequest":
			if OnlineCheck(msg.AddFriend.GetSender()) {
				go SendFirendRequestController(conn, msg)
			}
		case "sendFriendResponse":
			if OnlineCheck(msg.AddFriend.GetSender()) {
				go SendFirendResponseController(conn, msg)
			}
		case "allFriends":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GetAllFriendController(conn, msg)
			}

		case "searchUser":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go SearchUserController(conn, msg)
			}
		case "sendMsg":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go SendMsgController(conn, msg)
			}
		case "modifyRemark":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go ModifyRemarkController(conn, msg)
			}
		case "msgRecord":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go SingleMsgRecordController(conn, msg)
			}
		case "groupCreate":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupCreateController(conn, msg)
			}
		case "groupInvite":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupInviteController(conn, msg)
			}
		case "groupExit":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupExitController(conn, msg)
			}
		case "groupKick":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupKickController(conn, msg)
			}
		case "groupMsgRecord":
			//在msgRecord_controller这个文件中
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupMsgRecordController(conn, msg)
			}
		case "groupMsg":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupMsgController(conn, msg)
			}
		case "groupModify":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupModifyController(conn, msg)
			}
		case "groupRemark":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupRemarkController(conn, msg)
			}
		case "sendfile":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go SendFile(conn, msg)
			}
		case "recvfile":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go RecvFile(conn, msg)
			}
		case "friendInfo":
			go GetFriendInfoController(conn, msg)
		case "groupInfo":
			go GetGroupInfoController(conn, msg)

		case "delFriend":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go DelFriendController(conn, msg)
			}
		}
	}

Break:
	if ConnMap[conn] != "" {
		collectionDevice := common.DBMongo.C("device")
		collectionUser := common.DBMongo.C("user")
		nowTime := time.Now()
		// _, err := o.QueryTable("login_message").Filter("id", LoginTimeId).Update(orm.Params{"logout_time": nowTime})
		// _, err = o.QueryTable("user").Filter("user_name", ConnMap[conn]).Update(orm.Params{"logout_time": nowTime})
		err := collectionDevice.Update(bson.M{"_id": LoginTimeId}, bson.M{"$set": bson.M{"logout_time": nowTime}})
		err = collectionUser.Update(bson.M{"user_name": ConnMap[conn]}, bson.M{"$set": bson.M{"logout_time": nowTime}})
		if err != nil {
			Logger.Error("logout time update failed: ", err)
			fmt.Println("logout time update failed: ", err)
		}
		delete(UserMap, ConnMap[conn])
		delete(ConnMap, conn)
	}
	fmt.Println("退出： ", conn.RemoteAddr().String())
}
