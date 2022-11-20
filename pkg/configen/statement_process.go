package configen

import (
	"CodeGenerationGo/pkg/template"
	"CodeGenerationGo/pkg/util"
	"fmt"
	"regexp"
	"strconv"
)

// 将输入的语句转变成MatchRes格式
// 一条输入语句转换成一个MatchRes
func ParseStatement(statement string) template.MatchRes {
	//statement example:  required: app:appA & app:appB
	//result[0] : required: app:appA
	//Parse the regular expression and return the interpreter if successful
	//reg1 parses the whole statement
	//reg2 parses the sub
	reg1 := regexp.MustCompile(`(required|preferred):(\d*)\s(\w+):(\w+)\s(&|\^)\s(.+)`)
	reg2 := regexp.MustCompile(`(\w+):(\w+)`)
	if reg1 == nil || reg2 == nil {
		fmt.Println("syntax err")
	}

	result1 := reg1.FindAllStringSubmatch(statement, -1)
	result2 := reg2.FindAllStringSubmatch(result1[0][6], -1)

	var matchRes template.MatchRes

	matchRes.Trendrule = result1[0][1]
	matchRes.LabelKey = result1[0][3]
	matchRes.Value = result1[0][4]
	matchRes.Relationship = util.Relation2Opera(result1[0][5])

	if matchRes.Trendrule == "required" {
		matchRes.Weight = -1

	} else if matchRes.Trendrule == "preferred" {
		matchRes.Weight, _ = strconv.Atoi(result1[0][2]) //字符串转换为数字
	} else {
		fmt.Println("relationship word srr")
	}

	matches := make([]template.LabelSelectorRequirement, 0)

	var match template.LabelSelectorRequirement
	//key 是第一个元素的key值
	match.Key = result1[0][1]
	match.Operator = matchRes.Relationship
	for _, element := range result2 {
		match.Values = append(match.Values, element[2])
	}

	matches = append(matches, match)
	matchRes.MatchExpressions = matches

	return matchRes
}
