//应答服务脚本解释器
package main

import "sync"

var port int
var wg sync.WaitGroup

func main() {
	//初始化
	CreateLog()

	//运行服务
	wg.Add(1)
	go Server()
	wg.Wait()
}
