package configen

import (
	"CodeGenerationGo/template"
	"CodeGenerationGo/util"
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"regexp"
	"strconv"
)

func add() {
	fmt.Println("111")
}

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
	for _, element := range result2 {
		var match template.LabelSelectorRequirement
		match.Key = element[1]
		match.Values = append(match.Values, element[2])
		match.Operator = matchRes.Relationship

		matches = append(matches, match)

	}
	matchRes.MatchExpressions = matches

	return matchRes
}

// 初始化affinity
func AffinityInit() template.Affinity {
	affinity := template.Affinity{}
	return affinity
}

func InsertMatchRes2PodAffinity(affinity *template.Affinity, matchRes template.MatchRes) template.Affinity {
	var labelSelector template.LabelSelector
	labelSelector.MatchExpressions = matchRes.MatchExpressions
	labelSelector.MatchLabels = make(map[string]string) //分配内存
	labelSelector.MatchLabels[matchRes.LabelKey] = matchRes.Value

	podAffinityTerm := template.PodAffinityTerm{
		LabelSelector: &labelSelector,
		TopologyKey:   template.DEFAULT_TOPOLOGYKRY, //拓扑域采用默认的kubernetes.io/hostname
	}
	if matchRes.Relationship == template.LabelSelectorOpIn || matchRes.Relationship == template.LabelSelectorOpExists {
		if matchRes.Trendrule == "preferred" {

			var preference template.WeightedPodAffinityTerm
			preference = template.WeightedPodAffinityTerm{
				Weight:          matchRes.Weight,
				PodAffinityTerm: podAffinityTerm}

			if affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution == nil {
				affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution = []template.WeightedPodAffinityTerm{preference}

			} else {
				affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(
					affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution, preference)
			}

		}

		if matchRes.Trendrule == "required" {

			if affinity.PodAffinity == nil {
				affinity.PodAffinity = &template.PodAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: []template.PodAffinityTerm{podAffinityTerm}}
			} else {
				affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(
					affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution, podAffinityTerm)
			}

		}

	} else { //使用antiaffinity,同时将notin和notexist改为in与exist

		for index, requirement := range podAffinityTerm.LabelSelector.MatchExpressions {
			if requirement.Operator == template.LabelSelectorOpNotIn {
				podAffinityTerm.LabelSelector.MatchExpressions[index].Operator = template.LabelSelectorOpIn
			} else if requirement.Operator == template.LabelSelectorOpDoesNotExist {
				podAffinityTerm.LabelSelector.MatchExpressions[index].Operator = template.LabelSelectorOpExists
			}
		}

		if matchRes.Trendrule == "preferred" {

			var preference template.WeightedPodAffinityTerm
			preference = template.WeightedPodAffinityTerm{
				Weight:          matchRes.Weight,
				PodAffinityTerm: podAffinityTerm}

			if affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution == nil {
				affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = []template.WeightedPodAffinityTerm{preference}

			} else {
				affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(
					affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution, preference)
			}

		}

		if matchRes.Trendrule == "required" {

			if affinity.PodAntiAffinity == nil {
				affinity.PodAntiAffinity = &template.PodAntiAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: []template.PodAffinityTerm{podAffinityTerm}}
			} else {
				affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(
					affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution, podAffinityTerm)
			}

		}

	}

	return *affinity
}

func InsertAffinity2Yaml(statelist []string, sourcePath string, outPath string) {
	var affinity template.Affinity
	//将所有的语句串插入该affinity
	for _, state := range statelist {
		matches := ParseStatement(state)
		InsertMatchRes2PodAffinity(&affinity, matches)
	}

	pod, _ := util.ReadPodYamlFile(sourcePath)
	if pod.Spec.Affinity == nil {
		pod.Spec.Affinity = &affinity
	}
	yamlByte, _ := yaml.Marshal(pod)

	if err := os.WriteFile(outPath, yamlByte, 0666); err != nil {
		log.Fatal(err)
	}
}

//	func YamlGen(states []string, sourcePath string, outPath string) {
//		affinity := AffinityInit()
//		for _, state := range states {
//			match := ParseStatement(state)
//			affinity = InsertMatchRes2Affinity(affinity, match)
//		}
//		InsertAffinity2Yaml(affinity, sourcePath, outPath)
//
// }
func InsertYamlbyTxtstatement(statesfile string, sourcePath string, outPath string) {
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

	InsertAffinity2Yaml(statements, sourcePath, outPath)

}

func DeletePodStatusFromYaml(sourcePath string) {
	pod, _ := util.ReadPodYamlFile(sourcePath)
	deletePodStatus(pod)
	util.WriteObject2Yaml(pod, sourcePath)

}

func deletePodStatus(pod *template.Pod) {
	pod.Status = template.PodStatus{}

}
