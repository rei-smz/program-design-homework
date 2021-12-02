package sdsl

import (
	"github.com/gorilla/websocket"
	"testing"
	"time"
)

func tickWriter(connect *websocket.Conn, b *testing.B) {
	messageList := [...]string{"你好", "小笼包", "再见"}
	for _, message := range messageList {
		err := connect.WriteMessage(websocket.TextMessage, []byte(message))
		if nil != err {
			b.Error(err)
		} else {
			b.Logf("sent %s", message)
		}
		time.Sleep(time.Millisecond * 600)
	}
}

func RunWs(conn *websocket.Conn, b *testing.B)  {
	go tickWriter(conn, b)
	for {
		_, data, err := conn.ReadMessage()
		if nil != err {
			b.Error(err)
			break
		}
		b.Log(string(data))
	}
}

func BenchmarkConn(b *testing.B) {
	var connNum int64 = 20000
	dialer := websocket.DefaultDialer
	for i := int64(0); i < connNum; i++ {
		connect, _, err := dialer.Dial("ws://127.0.0.1:7777/ws", nil)
		if err != nil {
			b.Errorf("发生Websocket连接失败")
			continue
		}
		go RunWs(connect, b)
	}
}
