//SDSL 意为 Simple Domain Specific Language，SDSL Interpreter是一个简单领域定义语言应答服务脚本解释器

package main

import (
	"main/package/sdsl"
	"os"
)

var port int
var scriptPath string

func main() {
	//初始化
	port = 7777
	scriptPath = "script.txt"
	err := sdsl.CreateLog()
	if err != nil {
		os.Exit(1)
	}
	err = sdsl.ParseFile(scriptPath)
	if err != nil {
		os.Exit(1)
	}

	//运行服务
	err = sdsl.Server(port)
	if err != nil {
		os.Exit(1)
	}
}
