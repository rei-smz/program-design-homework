package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var scriptFile *os.File
var script Script
var newStep *Step

func ProcessDefault(nextStepId string)  {
	newStep._default = nextStepId
}

func ProcessSilence(nextStepId string)  {
	newStep.silence = nextStepId
}

func ProcessBranch(answer, nextStepId string)  {
	newStep.hashTable[answer] = nextStepId
}

func ProcessListen(start, end string)  {
	startTimer, _ := strconv.Atoi(start)
	endTimer, _ := strconv.Atoi(end)
	newStep.listen.beginTimer = int64(startTimer)
	newStep.listen.endTimer = int64(endTimer)
}

func ProcessStep(stepId string)  {
	newStep = new(Step)
	newStep.stepId = stepId
	if len(script.stepList) == 0 {
		script.mainStep = newStep
	}
	script.stepList[stepId] = newStep
}

func ProcessExpression(token []string)  {
	
}

func ProcessSpeak(token []string)  {
	ProcessExpression(token)
}

// ProcessTokens 对每一行的token进行处理
func ProcessTokens(tokens []string)  {
	switch tokens[0] {
	case "Step": {
		ProcessStep(tokens[1])
	}
	case "Speak": {
		ProcessSpeak(tokens[1:])
	}
	case "Listen": {
		ProcessListen(tokens[1], tokens[2])
	}
	case "Branch": {
		ProcessBranch(tokens[1], tokens[2])
	}
	case "Silence": {
		ProcessSilence(tokens[1])
	}
	case "Default": {
		ProcessDefault(tokens[1])
	}
	case "Exit": {
		script.exitStep = newStep
	}
	default: {
		Log("脚本文件包含无法识别的Token,退出程序")
		os.Exit(1)
	}
	}
}

// ParseLine 处理脚本文件的每一行
func ParseLine(line []string)  {
	var tokens []string
	for _, token := range line {
		if token[0] == '#' {
			break
		}
		tokens = append(tokens, token)
	}
	ProcessTokens(tokens)
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
