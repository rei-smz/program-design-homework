package main

type VarName string
type StepId string
type Answer string

type Expression struct {
	item []string
}

type Listen struct {
	beginTimer int64
	endTimer int64
}

type Step struct {
	stepId StepId
	speak Expression
	listen Listen
	hashTable map[Answer]StepId //{answer: nextStepId}
	silence StepId
	_default StepId
}

type Script struct {
	stepList map[string]*Step
	vars []VarName
	mainStep *Step
	exitStep *Step
}
