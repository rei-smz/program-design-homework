package sdsl

import (
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	// 读取存储空间大小
	ReadBufferSize: 1024,
	// 写入存储空间大小
	WriteBufferSize: 1024,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handler 处理前后端通信事务
func Handler(w http.ResponseWriter, r *http.Request) {
	var ws *websocket.Conn
	ws, _ = upgrader.Upgrade(w, r, nil)
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			return
		}
	}(ws)

	Log("新的Websocket连接")

	newControlBlock := new(UserControlBlock)
	err := ControlBlockInit(newControlBlock, DEFAULT, "")
	if err != nil {
		return
	}

	HandleMessage(ws, newControlBlock)
}

// Server 提供Web服务
func Server(port int) error {
	http.HandleFunc("/ws", Handler)
	http.Handle("/", http.FileServer(http.Dir("static")))
	err := http.ListenAndServe(":" + strconv.FormatInt(int64(port), 10), nil)
	Log("Web服务启动")
	if err != nil {
		Log("Web服务未正常工作")
		return err
	}
	return nil
}
