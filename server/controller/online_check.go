package controller

//判断此用户是否在线
func OnlineCheck(name string) bool {

	_, exist := UserMap[name]

	return exist
}
