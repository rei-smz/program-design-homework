package main

import (
	"github.com/gorilla/websocket"
	"strings"
	"time"
)

////var users map[uint64]*UserControlBlock
//
//func HandleStep(controlBlock *UserControlBlock)  {
//
//}
//
//// AddUser 处理前端添加用户的请求，返回新增用户的uid
//func AddUser(uid uint64, name string) {
//	//newControlBlock := new(UserControlBlock)
//	//newControlBlock.vars["$name"] = name
//	//newControlBlock.uid = uid
//	//newControlBlock.currentStep = script.mainStep
//	//users[uid] = newControlBlock
//	//callbackInfo := "{\"uid\":" + strconv.FormatUint(uid, 10)
//	//go HandleStep(newControlBlock)
//	//return []byte(callbackInfo)
//}
//
//func HandleMessage(uid uint64, message string) []byte {
//	if users[uid] == nil {
//		return []byte("{\"success\":false}")
//	}
//
//}

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

func HandleListen(conn *websocket.Conn, controlBlock *UserControlBlock) StepId {
	var (
		message []byte
		err error
		timer time.Duration
	)

	timer = time.Second * time.Duration(controlBlock.currentStep.listen.beginTimer)
	err = conn.SetReadDeadline(time.Now().Add(timer))
	if err != nil {
		Log("设置Listen定时器失败")
		return "error"
	}
	_, message, err = conn.ReadMessage()
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return controlBlock.currentStep.silence
		} else {
			Log("WebSocket通信出错")
			return "error"
		}
	}
	nextStep, ok := controlBlock.currentStep.hashTable[Answer(message)]
	if ok {
		return nextStep
	} else {
		return controlBlock.currentStep._default
	}
}

func HandleMessage(conn *websocket.Conn, controlBlock *UserControlBlock) {
	var err error

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
		controlBlock.currentStep = script.stepList[nextStep]
	}
}
