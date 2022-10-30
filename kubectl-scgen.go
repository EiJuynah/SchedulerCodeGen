package main

import (
	"CodeGenerationGo/pkg/configen"
	"flag"
	"fmt"
)

var VERSION = "1.04"

func main() {
	var podName = flag.String("name", "", "Input Pod Name")
	var version = flag.Bool("v", false, "output the project's version")
	var labelpair = flag.String("label", "", "assgin a pod with key:value")
	flag.Parse()
	//configen.DeletePodStatusFromYaml(".\\pod.yaml", ".\\newpod.yaml")

	//如果输入-v，输出代码生成器的version
	if *version {
		fmt.Println("Scheduler CodeGeneration version ", VERSION)
	}

	//通过podname选择pod
	if *podName != "" {
		SCFilePath := "./SCFile"
		configen.AddAffinityByPodname(*podName, SCFilePath)
	}

	//通过label选择pod
	if *labelpair != "" {
		SCFilePath := "./SCFile"
		configen.AddAffinityByLabel(*labelpair, SCFilePath)

	}

}
