package configen

import (
	"CodeGenerationGo/pkg/template"
	"CodeGenerationGo/pkg/yaml-process"
	"bufio"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// 初始化affinity
func AffinityInit() template.Affinity {
	affinity := template.Affinity{}
	return affinity
}

func InsertMatchRes2PodAffinity(affinity *template.Affinity, matchRes template.PodAffinityMatchDto) template.Affinity {
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

//将sclang的语句，插入到现有的yaml文件中
func insertAffinity2Yaml(statelist []string, sourcePath string, outPath string) {
	var affinity template.Affinity
	//将所有的语句串插入该affinity

	//插入podAffinity与podAntiAffinity
	for _, state := range statelist {
		matches := ParseStatement(state)
		InsertMatchRes2PodAffinity(&affinity, matches)
	}

	//读取源yaml文件的pod对象
	//将Affinity插入到pod对象中
	pod, _ := yaml_process.ReadPodYamlFile(sourcePath)
	if pod.Spec.Affinity == nil {
		pod.Spec.Affinity = &affinity
	}
	//再转换成yaml对象
	yamlByte, _ := yaml.Marshal(pod)
	//输出成yaml文件
	if err := os.WriteFile(outPath, yamlByte, 0666); err != nil {
		log.Fatal(err)
	}
}

// 根据scfile的sclang，读取源yaml文件，生成outpath插入过affinity的yaml文件
//代码生成主方法
func InsertYamlbyTxtstatement(scfilePath string, sourcePath string, outPath string) {
	var statements []string
	//读取scfile文件
	file, err := os.Open(scfilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	//读取sclang语句
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		statements = append(statements, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	//将sclang的语句，插入到现有的yaml文件中
	insertAffinity2Yaml(statements, sourcePath, outPath)

}

func DeletePodStatusFromYaml(sourcePath string, out string) {
	pod, _ := yaml_process.ReadPodYamlFile(sourcePath)
	deletePodStatus(pod)
	yaml_process.WriteObject2Yaml(pod, out)

}

func deletePodStatus(pod *template.Pod) {
	pod.Status = template.PodStatus{}

}
