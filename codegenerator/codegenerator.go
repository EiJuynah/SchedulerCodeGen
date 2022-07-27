package codegenerator

import (
	"CodeGenerationGo/Template"
	"CodeGenerationGo/util"
	"fmt"
	"regexp"
	"strconv"
)

func ParseStatement(statement string) Template.MatchRes {
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

	var matchRes Template.MatchRes

	matchRes.Trendrule = result1[0][1]
	matchRes.Key = result1[0][3]
	matchRes.Value = result1[0][4]
	matchRes.Relationship = result1[0][5]
	matchRes.Relationship = util.Relation2Opera(matchRes.Relationship)

	if matchRes.Trendrule == "required" {
		matchRes.Weight = -1

	} else if matchRes.Trendrule == "preferred" {
		matchRes.Weight, _ = strconv.Atoi(result1[0][2])
	} else {
		fmt.Println("relationship word srr")
	}

	matches := make([]Template.MatchExpression, 0)
	for _, element := range result2 {
		var match Template.MatchExpression
		match.Key = element[1]
		match.Values = append(match.Values, element[2])
		match.Operator = matchRes.Relationship

		matches = append(matches, match)

	}
	matchRes.MatchExpressions = matches
	//fmt.Println("result1 = ", result1[0])
	//fmt.Println("result2 = ", result2)

	return matchRes
}

func AffinityInit() Template.Affinity {
	affinity := Template.Affinity{}

	return affinity
}
