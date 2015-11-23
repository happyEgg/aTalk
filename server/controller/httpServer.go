package controller

import (
	"aTalk/protocol"
	"aTalk/server/common"
	"aTalk/server/config"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

func Httpserver() {
	addr := config.String("http_addr")
	log.Printf(">>>>>>>websocket端口:%s<<<<<<<<<<", addr)
	http.HandleFunc("/chat", handler)
	http.Handle("/", websocket.Handler(connHandler))
	err := http.ListenAndServe(":"+addr, nil)
	if err != nil {
		fmt.Println("httpserver:", err)
		config.Logger.Error("http server:", err)
	}

	// log.Println(">>>>>>>>websocket服务器启动<<<<<<<<")
	// log.Println(">>>>>>>>等待客户端连接<<<<<<<<<")
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintf(w, "这是一个测试")
}

func connHandler(ws *websocket.Conn) {
	defer ws.Close()
	headBuffer := make([]byte, 4)
	//c.ws.SetReadDeadline(time.Now().Add(pongWait))

	//循环处理接收数据
	for {
		bufSize, err := ws.Read(headBuffer)
		if err != nil {
			break
		}

		if bufSize < 4 {
			continue
		}

		messager := make(chan byte)
		//心跳计时
		go common.HeartBeating(ws, messager, 20)
		//检测是client否有数据传来
		go common.GravelChannel(headBuffer, messager)

		bodyLen := common.BytesToInt(headBuffer)
		if bodyLen == 0 {
			common.KeepAliveResult(ws)
			continue
		}
		bodyBuffer := make([]byte, bodyLen)
		bodySize, err := ws.Read(bodyBuffer)
		if err != nil {
			break
		}

		//如果接收的包大小小于包头内容的长度，舍弃
		if bodySize < bodyLen {
			continue
		}

		//类型判断
		msg := &protocol.WMessage{}
		err = proto.Unmarshal(bodyBuffer, msg)
		fmt.Println("msg:", msg.String())
		if err != nil {
			config.Logger.Error("http server unmarshal:", err)
			continue
		}

		switch msg.GetMsgType() {
		case "register":
			RegisterController(ws, msg)
			goto Break
		//只允许一方登陆
		case "login":
			LoginController(ws, msg)

		case "logout":

			fmt.Println(msg.UserInfo.GetUsername(), " logout")
			goto Break
		case "modifyInfo":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go ModifyInfoController(ws, msg)
			}

		case "sendFriendRequest":
			if OnlineCheck(msg.AddFriend.GetSender()) {
				go SendFirendRequestController(ws, msg)
			}
		case "sendFriendResponse":
			if OnlineCheck(msg.AddFriend.GetSender()) {
				go SendFirendResponseController(ws, msg)
			}

		case "allFriends":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GetAllFriendController(ws, msg)
			}

		case "searchUser":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go SearchUserController(ws, msg)
			}
		case "sendMsg":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go SendMsgController(ws, msg)
			}
		case "modifyRemark":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go ModifyRemarkController(ws, msg)
			}
		case "msgRecord":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go SingleMsgRecordController(ws, msg)
			}
		case "groupCreate":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupCreateController(ws, msg)
			}
		case "groupInvite":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupInviteController(ws, msg)
			}
		case "groupExit":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupExitController(ws, msg)
			}
		case "groupKick":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupKickController(ws, msg)
			}
		case "groupMsgRecord":
			//在msgRecord_controller这个文件中
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupMsgRecordController(ws, msg)
			}
		case "groupMsg":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupMsgController(ws, msg)
			}
		case "groupModify":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupModifyController(ws, msg)
			}
		case "groupRemark":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go GroupRemarkController(ws, msg)
			}
		case "sendfile":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go SendFile(ws, msg)
			}
		case "recvfile":
			if OnlineCheck(msg.UserInfo.GetUsername()) {
				go RecvFile(ws, msg)
			}
		case "friendInfo":
			go GetFriendInfoController(ws, msg)
		case "groupInfo":
			go GetGroupInfoController(ws, msg)
		}
	}

Break:
	if (ConnMap[*ws]) != "" {
		collectionDevice := common.DBMongo.C("device")
		collectionUser := common.DBMongo.C("user")
		nowTime := time.Now()
		// _, err := o.QueryTable("login_message").Filter("id", LoginTimeId).Update(orm.Params{"logout_time": nowTime})
		// _, err = o.QueryTable("user").Filter("user_name", ConnMap[conn]).Update(orm.Params{"logout_time": nowTime})
		err := collectionDevice.Update(bson.M{"_id": LoginTimeId}, bson.M{"$set": bson.M{"logout_time": nowTime}})
		err = collectionUser.Update(bson.M{"user_name": ConnMap[*ws]}, bson.M{"$set": bson.M{"logout_time": nowTime}})
		if err != nil {
			Logger.Error("logout time update failed: ", err)
			fmt.Println("logout time update failed: ", err)
		}
		delete(UserMap, ConnMap[*ws])
		delete(ConnMap, *ws)
	}
	fmt.Println("退出： ", ws.RemoteAddr().String())
}
