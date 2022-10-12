package util

import (
	"CodeGenerationGo/template"
	"fmt"
	"math/rand"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func Test() {

}

func readInputTxt() {

}

func AnalysisRelationalStatement(statement string) {

}

func Relation2Opera(relat string) template.LabelSelectorOperator {
	if relat == "&" {
		return template.LabelSelectorOpIn
	}
	if relat == "^" {
		return template.LabelSelectorOpNotIn
	}
	return "nil"
}

//	func ContainsInArray[T any](items []T, item T) bool {
//		for _, eachItem := range items {
//			bytes.Compare(eachItem, item)
//
//		}
//		return false
//	}
//
// 获取随机字母+数字组合字符串
func getRandstring(length int) string {
	if length < 1 {
		return ""
	}
	char := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charArr := strings.Split(char, "")
	charlen := len(charArr)
	ran := rand.New(rand.NewSource(time.Now().Unix()))
	var rchar string = ""
	for i := 1; i <= length; i++ {
		rchar = rchar + charArr[ran.Intn(charlen)]
	}
	return rchar
}

// 随机文件名
func RandFileName(fileName string) string {
	randStr := getRandstring(16)
	return randStr + filepath.Ext(fileName)
}

func runInLinux(cmd string) string {
	fmt.Println("Running Linux cmd:", cmd)
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	return strings.TrimSpace(string(result))
}

func runInWindows(cmd string) string {
	fmt.Println("Running Win cmd:", cmd)
	result, err := exec.Command("cmd", "/c", cmd).Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	return strings.TrimSpace(string(result))
}

func RunCommand(cmd string) string {
	if runtime.GOOS == "windows" {
		return runInWindows(cmd)
	} else {
		return runInLinux(cmd)
	}
}
