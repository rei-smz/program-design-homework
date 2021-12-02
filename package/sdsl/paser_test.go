package sdsl

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseLine(t *testing.T) {
	var testStr string
	fmt.Scanf("输入一行脚本：%s", testStr)
	testLine := strings.Split(testStr, " ")
	err := ParseLine(testLine)
	if err != nil {
		t.Errorf("PaseLine出现错误，错误为%s", err.Error())
	}
}

func TestParseFile(t *testing.T) {
	err := ParseFile("../../script.txt")
	if err != nil {
		t.Errorf("PaseFile出现错误，错误为%s", err.Error())
	}
}
