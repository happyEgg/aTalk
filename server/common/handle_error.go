/*
处理错误
time: 2015-11-17
*/

package common

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	//"strconv"
// )

// func HandleError(filename string, err error) {
// 	//fmt.Println("Err:", err)
// 	//flienameOnly := GetCurFilename()
// 	logFilename := "diary/" + filename + ".log"
// 	logFile, fileErr := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
// 	if fileErr != nil {
// 		fmt.Printf("open file error=%s\r\n", err.Error())
// 		os.Exit(-1)
// 	}
// 	defer logFile.Close()
// 	logger := log.New(logFile, "\r\n", log.Ldate|log.Ltime)
// 	logger.Println(err)

// 	os.Exit(-1)

// }

func CheckErr(err error) {
	if err != nil {
		Logger.Error("write error: ", err)
		panic(err)
	}
}
