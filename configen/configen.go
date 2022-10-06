package configen

import (
	"CodeGenerationGo/template"
	"CodeGenerationGo/util"
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"os"
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
	matchRes.Key = result1[0][3]
	matchRes.Value = result1[0][4]
	matchRes.Relationship = util.Relation2Opera(result1[0][5])

	if matchRes.Trendrule == "required" {
		matchRes.Weight = -1

	} else if matchRes.Trendrule == "preferred" {
		matchRes.Weight, _ = strconv.Atoi(result1[0][2]) //字符串转换为数字
	} else {
		fmt.Println("relationship word srr")
	}

	matches := make([]metav1.LabelSelectorRequirement, 0)
	for _, element := range result2 {
		var match metav1.LabelSelectorRequirement
		match.Key = element[1]
		match.Values = append(match.Values, element[2])
		match.Operator = matchRes.Relationship

		matches = append(matches, match)

	}
	matchRes.MatchExpressions = matches

	return matchRes
}

// 初始化affinity
func AffinityInit() v1.Affinity {
	affinity := v1.Affinity{}

	return affinity
}

func InsertMatchRes2Affinity(affinity v1.Affinity, matchRes template.MatchRes) v1.Affinity {
	var labelSelector metav1.LabelSelector
	labelSelector.MatchExpressions = matchRes.MatchExpressions

	if matchRes.Trendrule == "preferred" {
		//TODO
		//var preference template.Perference
		//preference.Weight = matchRes.Weight
		//preference.PodAffinityTerm.LabelSelector = append(
		//	preference.PodAffinityTerm.LabelSelector, labelSelector)
		//affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution.Preference = append(
		//	affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution.Preference, preference)
	}

	if matchRes.Trendrule == "required" {
		affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution.LabelSelector = append(
			affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution.LabelSelector, labelSelector)
	}

	return affinity
}

func InsertAffinity2Yaml(affinity template.Affinity, sourcePath string, outPath string) {
	config, _ := util.ReadConfigYaml(sourcePath)
	config.Spec.Affinity = affinity
	yamlByte, _ := yaml.Marshal(config)

	if err := os.WriteFile(outPath, yamlByte, 0666); err != nil {
		log.Fatal(err)
	}
}

func YamlGen(states []string, sourcePath string, outPath string) {
	affinity := AffinityInit()
	for _, state := range states {
		match := ParseStatement(state)
		affinity = InsertMatchRes2Affinity(affinity, match)
	}
	InsertAffinity2Yaml(affinity, sourcePath, outPath)

}
func YamlGenbyTxt(statesfile string, sourcePath string, outPath string) {
	var statements []string
	file, err := os.Open(statesfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		statements = append(statements, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	YamlGen(statements, sourcePath, outPath)

}
