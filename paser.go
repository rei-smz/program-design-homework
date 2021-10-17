package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var scriptFile *os.File

// ParseLine 处理脚本文件的每一行
func ParseLine(line []string)  {

}

// ParseFile 处理脚本文件
func ParseFile(fileName string)  {
	var err error
	scriptFile, err = os.Open(fileName)
	if err != nil {
		fmt.Println("打开脚本文件失败，请检查文件路径")
		os.Exit(1)
	}

	defer func(scriptFile *os.File) {
		err := scriptFile.Close()
		if err != nil {
			Log("关闭脚本文件时发生错误")
		}
	}(scriptFile)

	buf := bufio.NewReader(scriptFile)
	for {
		str, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			Log("读取脚本文件时发生错误，退出程序")
			os.Exit(1)
		}
		if str[0] != '#' {
			strings.Trim(str, " ")
			line := strings.Split(str, " ")
			ParseLine(line)
		}
	}
	Log("处理脚本文件完成")
}
