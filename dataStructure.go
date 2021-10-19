package main

type Expression struct {
	varName []string
	other []string
}

type Listen struct {
	beginTimer int64
	endTimer int64
}

type Step struct {
	stepId string
	speak Expression
	listen Listen
	hashTable map[string]string //{answer: nextStepId}
	silence string
	_default string
}

type Script struct {
	stepList map[string]*Step
	vars map[string]string
	mainStep *Step
	exitStep *Step
}
