package sdsl

import (
	"fmt"
	"os"
	"time"
)

var logFile *os.File

// CreateLog 创建日志文件
func CreateLog() error {
	var err error
	logFile, err = os.Create("log/" + time.Now().Format("20060102_150405") + ".log")
	if err != nil {
		fmt.Println("创建日志文件失败，退出程序")
		return err
	}
	return nil
}

// Log 写入日志
func Log(message string) {
	var err error
	message = "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + message
	_, err = logFile.WriteString(message + "\n")
	if err != nil {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " 发生Log文件写入错误")
	} else {
		fmt.Println(message)
	}
}
