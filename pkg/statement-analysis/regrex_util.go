package statement_analysis

import (
	"CodeGenerationGo/pkg/template"
	"fmt"
	"math/rand"
	"path/filepath"
	"regexp"
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

// 正则表达式分割字符串
// 将 aaaOPbbb 字符串，分割为aaa，OP， bbb三个字段
// return {key,value}
func SplitStringByOP(statement string, OP string) []string {
	//res := make([]string, 3)
	regstr := `(\w+)` + OP + `(\w+)`
	reg := regexp.MustCompile(regstr)
	if reg == nil {
		fmt.Println("syntax err")
		return nil
	}
	res := reg.FindAllStringSubmatch(statement, -1)
	if res[0][1] == "" && res[0][2] == "" {
		fmt.Println("statement input err")
	}
	return []string{res[0][1], res[0][2]}

}
