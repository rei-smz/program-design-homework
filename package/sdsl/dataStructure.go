package sdsl

// VarName 表示脚本中的变量名
type VarName string

// StepId 表示步骤ID
type StepId string

// Answer 表示用户的命令
type Answer string

// Expression 表示Speak的表达式
type Expression struct {
	item []string
}

// Listen 是用于控制Listen的变量
type Listen struct {
	beginTimer int64
	endTimer   int64
}

// Step 是步骤的数据结构
type Step struct {
	stepId    StepId
	speak     Expression
	listen    Listen
	hashTable map[Answer]StepId //{answer: nextStepId}
	silence   StepId
	_default  StepId
}

// Script 是脚本的数据结构
type Script struct {
	stepList map[StepId]*Step
	vars     []VarName
	mainStep *Step
	exitStep *Step
}

// UserControlBlock 是控制块的数据结构
type UserControlBlock struct {
	vars        map[VarName]string
	currentStep *Step
}
