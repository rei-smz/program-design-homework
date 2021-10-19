//应答服务脚本解释器
package main

import "sync"

var port int
var scriptPath string
var wg sync.WaitGroup

func main() {
	//初始化
	CreateLog()
	ParseFile(scriptPath)

	//运行服务
	wg.Add(1)
	go Server()
	wg.Wait()
}
