/*
把protobuf和修改的数据进行转换
*/
package common

import (
	"aTalk/protocol"
	"aTalk/server/models"
)

//客户端传送给服务端需要修改的字段
func ProtoToData(msg *protocol.WMessage) map[string]interface{} {
	var arr = make(map[string]interface{})

	if msg.UserInfo.Password != nil {
		arr["password"] = msg.UserInfo.GetPassword()
	}
	if msg.UserInfo.NewPassword != nil {
		arr["new_password"] = msg.UserInfo.GetNewPassword()
	}
	if msg.UserInfo.RealName != nil {
		arr["real_name"] = msg.UserInfo.GetRealName()
	}
	if msg.UserInfo.Sex != nil {
		arr["gender"] = int32(msg.UserInfo.GetSex())
	}
	if msg.UserInfo.Age != nil {
		arr["age"] = msg.UserInfo.GetAge()
	}
	if msg.UserInfo.Phone != nil {
		arr["phone"] = msg.UserInfo.GetPhone()
	}
	if msg.UserInfo.Email != nil {
		arr["email"] = msg.UserInfo.GetEmail()
	}
	if msg.UserInfo.ProfilePhotoid != nil {
		arr["profile_photoid"] = msg.UserInfo.GetProfilePhotoid()
	}
	if msg.UserInfo.ChatGroundid != nil {
		arr["chat_groundid"] = msg.UserInfo.GetChatGroundid()
	}
	if msg.UserInfo.Notice != nil {
		arr["notice"] = msg.UserInfo.GetNotice()
	}

	return arr
}

//服务端返回给客户端的数据
func DataToProto(user *models.User) *protocol.User {
	var msg *protocol.User
	gender := protocol.Gender(user.Gender)

	msg.Username = &user.UserName
	msg.RealName = &user.RealName
	msg.Sex = &gender
	msg.Age = &user.Age
	msg.Phone = &user.Phone
	msg.Email = &user.Email
	msg.ProfilePhotoid = &user.ProfilePhotoid
	msg.ChatGroundid = &user.ChatGroundid
	msg.Notice = &user.Notice
	return msg
}
