package main

import (
	"CodeGenerationGo/pkg/configen"
	"flag"
	"fmt"
	"log"
	"os/exec"
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
	if *podName != "" {
		SCFilePath := "./SCFile"
		configen.AddAffinityByPodname(*podName, SCFilePath)
	}
	if *labelpair != "" {
		SCFilePath := "./SCFile"
		configen.AddAffinityByLabel(*labelpair, SCFilePath)

	}

}

func Reschedule(podName string) {
	//args := os.Args

	//reg := regexp.MustCompile(`(\w+):(\w+)`)
	//podinfo := reg.FindAllStringSubmatch(args[1], -1)
	//podLabelKey := podinfo[0]
	//podValue := podinfo[1]

	////第一个参数是想要调度的pod的name
	//podName := args[1]
	SCFilePath := "./SCFile"

	//first step
	//获取当前pod的配置信息，将其输出重定向到pod.yaml中
	//kubectl get -o yaml pod {podname} > pod.yaml
	//此时的commmand为Windows的cmd，如果是linux环境，换成 bin/bash
	cmd1 := exec.Command("cmd", "/c", "kubectl", "get", "-o", "yaml", "pod", podName, ">", "pod.yaml")
	if err := cmd1.Run(); err != nil { // 运行命令
		log.Fatal(err)
	}

	//second step
	//删除当前的pod
	//kubectl delete pod <podname> -n <namespace>
	cmd2 := exec.Command("cmd", "/c", "kubectl", "delete", "pod", podName)
	if err := cmd2.Run(); err != nil { // 运行命令
		log.Fatal(err)
	}

	//third step
	//根据SCFile中的sclang，修改当前的配置文件的affinity，并删除status，生成新的调度文件newpod.yaml于当前目录
	configen.InsertYamlbyTxtstatement(SCFilePath, "./pod.yaml", "./newpod.yaml")

	//forth step
	//根据新生成的调度文件，重新启动一个新的pod
	//kubectl apply -f ./newpod.yaml
	cmd := exec.Command("cmd", "/c", "kubectl", "apply", "-f", "./newpod.yaml")
	if err := cmd.Run(); err != nil { // 运行命令
		log.Fatal(err)
	}

	fmt.Println("config yaml is generated")

}
