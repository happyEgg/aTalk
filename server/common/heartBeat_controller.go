package common

import (
	"code.google.com/p/go.net/websocket"
	"net"
	//"reflect"
	"time"
)

//心跳超时时间限制
func HeartBeating(conn interface{}, readerChannel chan byte, timeout int) {
	select {
	case <-readerChannel:
		switch conn.(type) {
		case net.Conn:
			conn.(net.Conn).SetReadDeadline(time.Now().Add(time.Minute * time.Duration(timeout)))
		}
		break
	case <-time.After(time.Second * 5):
		switch conn.(type) {
		case net.Conn:

			conn.(net.Conn).Close()
		case websocket.Conn:

			conn.(*websocket.Conn).Close()
		}
	}
}

func GravelChannel(n []byte, mess chan byte) {
	for _, v := range n {
		mess <- v
	}
	close(mess)
}
