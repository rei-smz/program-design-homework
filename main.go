//应答服务脚本解释器
package main

import "sync"

var port int
var scriptPath string
var wg sync.WaitGroup

func main() {
	//初始化
	port = 7777
	scriptPath = "script.txt"
	CreateLog()
	ParseFile(scriptPath)

	//运行服务
	wg.Add(1)
	go Server()
	wg.Wait()
}
