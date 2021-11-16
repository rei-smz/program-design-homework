package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// ShutdownServer 关闭Web服务
func ShutdownServer()  {
	err := logFile.Close()
	if err != nil {
		Log("关闭服务失败，原因是日志文件无法关闭")
		return
	}
	Log("关闭服务")
	os.Exit(0)
}

// Handler 处理前后端通信事务
func Handler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("content-type", "application/json")

	query := r.URL.Query()

	//TODO: 请求处理部分
	switch r.URL.Path {
	case "/api/check-service": {
		_, _ = w.Write([]byte("ok"))
	}
	case "/api/new-user": {
		username := query.Get("username")
		uid, _ := strconv.Atoi(query.Get("uid"))
		go w.Write(AddUser(uint64(uid), username))
	}
		
	}
}

// CheckService 检查Web服务是否正常运行
func CheckService()  {
	for {
		res, err := http.Get("http://localhost:" + strconv.FormatInt(int64(port), 10) + "/api/check-service")
		if err == nil {
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					Log("关闭请求体失败")
				}
			}(res.Body)
			body, _ := ioutil.ReadAll(res.Body)
			if string(body) == "ok" {
				Log("Web服务工作正常")
			}
		}
	}
}

// Server 提供Web服务
func Server()  {
	http.HandleFunc("/api/", Handler)
	http.Handle("/", http.FileServer(http.Dir("static")))
	go CheckService()
	err := http.ListenAndServe(":" + strconv.FormatInt(int64(port), 10), nil)
	if err != nil {
		Log("Web服务未正常工作")
		os.Exit(1)
	}
	wg.Done()
}
