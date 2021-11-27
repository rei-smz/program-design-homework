package main

import (
	"github.com/gorilla/websocket"
	"reflect"
	"strings"
	"time"
	"unsafe"
)


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
	Log(string(message))
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
