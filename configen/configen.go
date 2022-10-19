package configen

import (
	"CodeGenerationGo/template"
	"CodeGenerationGo/util"
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func add() {
	fmt.Println("111")
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

func insertAffinity2Yaml(statelist []string, sourcePath string, outPath string) {
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

// 根据
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

	insertAffinity2Yaml(statements, sourcePath, outPath)

}

func DeletePodStatusFromYaml(sourcePath string, out string) {
	pod, _ := util.ReadPodYamlFile(sourcePath)
	deletePodStatus(pod)
	util.WriteObject2Yaml(pod, out)

}

func deletePodStatus(pod *template.Pod) {
	pod.Status = template.PodStatus{}

}
