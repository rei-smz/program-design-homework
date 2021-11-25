package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"strconv"
)

var upgrader = websocket.Upgrader {
	// 读取存储空间大小
	ReadBufferSize:1024,
	// 写入存储空间大小
	WriteBufferSize:1024,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ws *websocket.Conn


//// ShutdownServer 关闭Web服务
//func ShutdownServer()  {
//	err := logFile.Close()
//	if err != nil {
//		Log("关闭服务失败，原因是日志文件无法关闭")
//		return
//	}
//	Log("关闭服务")
//	os.Exit(0)
//}

// Handler 处理前后端通信事务
func Handler(w http.ResponseWriter, r *http.Request)  {
	ws, _ = upgrader.Upgrade(w, r, nil)
	defer ws.Close()

	newControlBlock := new(UserControlBlock)
	for _, varName := range script.vars {
		newControlBlock.vars[varName] = ""
	}
	newControlBlock.currentStep = script.mainStep

	HandleMessage(ws, newControlBlock)
	//w.Header().Set("Access-Control-Allow-Origin","*")
	//w.Header().Set("content-type", "application/json")
	//
	//query := r.URL.Query()
	//
	////TODO: 请求处理部分
	//switch r.URL.Path {
	//case "/api/check-service": {
	//	_, _ = w.Write([]byte("ok"))
	//}
	//case "/api/new-user": {
	//	username := query.Get("username")
	//	i := 0
	//	for users[uint64(i)] != nil {
	//		i += 1
	//	}
	//	uid := i
	//	go func() {
	//		_, err := w.Write(AddUser(uint64(uid), username))
	//		if err != nil {
	//			Log("向前端发送消息失败")
	//		}
	//	}()
	//}
	//case "/api/send-msg": {
	//	uid, _ := strconv.Atoi(query.Get("uid"))
	//	message := query.Get("msg")
	//	go func() {
	//		_, err := w.Write(HandleMessage(uint64(uid), message))
	//		if err != nil {
	//			Log("向前端发送消息失败")
	//		}
	//	}()
	//}
	//
	//}
}

//// CheckService 检查Web服务是否正常运行
//func CheckService()  {
//	for {
//		res, err := http.Get("http://localhost:" + strconv.FormatInt(int64(port), 10) + "/api/check-service")
//		if err == nil {
//			defer func(Body io.ReadCloser) {
//				err := Body.Close()
//				if err != nil {
//					Log("关闭请求体失败")
//				}
//			}(res.Body)
//			body, _ := ioutil.ReadAll(res.Body)
//			if string(body) == "ok" {
//				Log("Web服务工作正常")
//			}
//		}
//	}
//}

// Server 提供Web服务
func Server()  {
	http.HandleFunc("/ws", Handler)
	http.Handle("/", http.FileServer(http.Dir("static")))
	//go CheckService()
	err := http.ListenAndServe(":" + strconv.FormatInt(int64(port), 10), nil)
	if err != nil {
		Log("Web服务未正常工作")
		os.Exit(1)
	}
	wg.Done()
}
