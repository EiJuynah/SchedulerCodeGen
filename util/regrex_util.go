package util

import (
	"CodeGenerationGo/template"
	"math/rand"
	"path/filepath"
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
