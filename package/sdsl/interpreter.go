package sdsl

import (
	"bufio"
	"github.com/gorilla/websocket"
	"io"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	// DEFAULT 代表初始化控制块的默认模式（填充随机浮点数）
	DEFAULT = 0
	// FILE 代表使用文件初始化控制块
	FILE = 1
	// SQL 代表使用数据库初始化控制块
	SQL = 2
)

// ControlBlockInit 初始化用户控制块
func ControlBlockInit(controlBlock *UserControlBlock, mode int, address string) error {
	controlBlock.vars = make(map[VarName]string)
	for _, varName := range script.vars {
		controlBlock.vars[varName] = ""
	}
	controlBlock.currentStep = script.mainStep
	switch mode {
	case DEFAULT:
		for _, varName := range script.vars {
			controlBlock.vars[varName] = strconv.FormatFloat(rand.Float64(), 'f', -1, 64)
		}
		break
	case FILE:
		{
			file, err := os.Open(address)
			if err != nil {
				return err
			}
			buf := bufio.NewReader(file)
			for {
				str, err := buf.ReadString('\n')
				if err == io.EOF {
					break
				}
				line := strings.Split(str, " ")
				if _, ok := controlBlock.vars[VarName(line[0])]; ok {
					controlBlock.vars[VarName(line[0])] = line[1]
				}
			}
			file.Close()
			break
		}
	case SQL:
		break
	}
	return nil
}

//HandleSpeak 完成脚本中Speak语句的功能
func HandleSpeak(controlBlock *UserControlBlock) []byte {
	expression := controlBlock.currentStep.speak
	ret := ""
	for _, item := range expression.item {
		if item[0] == '$' {
			ret += controlBlock.vars[VarName(item)]
		} else {
			ret += item
		}
	}
	return []byte(ret)
}

//HandleListen 完成脚本中Listen语句的功能
func HandleListen(conn *websocket.Conn, controlBlock *UserControlBlock) StepId {
	var (
		message []byte
		err     error
		timer   time.Duration
	)

	timer = time.Second * time.Duration(controlBlock.currentStep.listen.beginTimer)

	//设置静默超时
	err = conn.SetReadDeadline(time.Now().Add(timer))
	if err != nil {
		Log("设置Listen定时器失败")
		return "error"
	}

	//设置命令长度限制
	conn.SetReadLimit(controlBlock.currentStep.listen.endTimer)

	//这部分将之前产生的超时错误清除
	connChanger := reflect.Indirect(reflect.ValueOf(conn))
	rdErr := connChanger.FieldByName("readErr")
	ptr := unsafe.Pointer(rdErr.UnsafeAddr())
	*(*error)(ptr) = nil

	_, message, err = conn.ReadMessage()
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return controlBlock.currentStep.silence
		} else {
			Log("WebSocket通信出错")
			return "error"
		}
	}
	Log("收到消息：" + string(message))
	nextStep, ok := controlBlock.currentStep.hashTable[Answer(message)]
	if ok {
		return nextStep
	} else {
		return controlBlock.currentStep._default
	}
}

//HandleMessage 处理前端命令
func HandleMessage(conn *websocket.Conn, controlBlock *UserControlBlock) {
	var err error
	var ok bool

	for {
		err = conn.WriteMessage(websocket.TextMessage, HandleSpeak(controlBlock))
		if err != nil {
			Log("WebSocket通信出错")
			return
		}
		if controlBlock.currentStep == script.exitStep {
			return
		}
		nextStep := HandleListen(conn, controlBlock)
		if nextStep == "error" {
			return
		}
		if controlBlock.currentStep, ok = script.stepList[nextStep]; !ok {
			Log("脚本文件存在语法错误：未定义的Step")
			return
		}
	}
}
