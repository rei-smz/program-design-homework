package sdsl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var scriptFile *os.File
var script Script
var newStep *Step

//ProcessDefault 处理Default分支
func ProcessDefault(nextStepId string) {
	newStep._default = StepId(nextStepId)
}

//ProcessSilence 处理Silence分支
func ProcessSilence(nextStepId string) {
	newStep.silence = StepId(nextStepId)
}

//ProcessBranch 处理Branch分支
func ProcessBranch(answer, nextStepId string) {
	newStep.hashTable[Answer(answer)] = StepId(nextStepId)
}

//ProcessListen 处理Listen分支
func ProcessListen(start, end string) error {
	startTimer, err1 := strconv.Atoi(start)
	if err1 != nil {
		return err1
	}
	endTimer, err2 := strconv.Atoi(end)
	if err2 != nil {
		return err2
	}
	newStep.listen.beginTimer = int64(startTimer)
	newStep.listen.endTimer = int64(endTimer)
	return nil
}

// ProcessStep 处理Step分支
func ProcessStep(stepId StepId) {
	newStep = new(Step)
	newStep.stepId = stepId
	newStep.hashTable = make(map[Answer]StepId)
	if len(script.stepList) == 0 {
		script.mainStep = newStep
	}
	script.stepList[stepId] = newStep
}

// ProcessExpression 将token转为表达式返回给Speak
func ProcessExpression(tokens []string) Expression {
	var ret Expression
	for _, token := range tokens {
		if token == "+" {
			continue
		}
		ret.item = append(ret.item, token)
		if token[0] == '$' {
			script.vars = append(script.vars, VarName(token))
		}
	}
	return ret
}

// ProcessSpeak 处理Speak分支
func ProcessSpeak(token []string) {
	newStep.speak = ProcessExpression(token)
}

// ProcessTokens 对每一行的token进行处理
func ProcessTokens(tokens []string) error {
	switch tokens[0] {
	case "Step":
		{
			ProcessStep(StepId(tokens[1]))
			break
		}
	case "Speak":
		{
			ProcessSpeak(tokens[1:])
			break
		}
	case "Listen":
		{
			err := ProcessListen(tokens[1], tokens[2])
			if err != nil {
				return err
			}
			break
		}
	case "Branch":
		{
			ProcessBranch(tokens[1], tokens[2])
			break
		}
	case "Silence":
		{
			ProcessSilence(tokens[1])
			break
		}
	case "Default":
		{
			ProcessDefault(tokens[1])
			break
		}
	case "Exit":
		{
			script.exitStep = newStep
			break
		}
	default:
		{
			Log("脚本文件包含无法识别的Token")
			return errors.New("脚本文件包含无法识别的Token")
		}
	}
	return nil
}

// ParseLine 处理脚本文件的每一行
func ParseLine(line []string) error {
	var tokens []string
	for _, token := range line {
		if token[0] == '#' {
			break
		}

		tokens = append(tokens, strings.TrimSpace(token))
	}
	err := ProcessTokens(tokens)
	return err
}

// ParseFile 处理脚本文件
func ParseFile(fileName string) error {
	var err error
	scriptFile, err = os.Open(fileName)
	if err != nil {
		fmt.Println("打开脚本文件失败，请检查文件路径")
		return err
	}

	buf := bufio.NewReader(scriptFile)
	script.stepList = make(map[StepId]*Step)
	for {
		str, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			Log("读取脚本文件时发生错误")
			return err
		}
		if str[0] != '#' {
			line := strings.Split(str, " ")
			err := ParseLine(line)
			if err != nil {
				return err
			}
		}
	}
	defer func(scriptFile *os.File) {
		err := scriptFile.Close()
		if err != nil {
			Log("关闭脚本文件时发生错误")
		}
	}(scriptFile)
	Log("处理脚本文件完成")
	return nil
}
