package main

import (
	"CodeGenerationGo/configen"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	configen.DeletePodStatusFromYaml(".\\file\\test1.yaml")
}

func a() {
	args := os.Args
	////sourceyamlpath := ".\\files\\source.yaml"
	////outputyamlpath := ".\\files\\target.yaml"
	////configpath := ".\\files\\input.txt"

	if args[1] != "" {
		fmt.Println(args[1])
	}
	//reg := regexp.MustCompile(`(\w+):(\w+)`)
	//podinfo := reg.FindAllStringSubmatch(args[1], -1)
	//podLabelKey := podinfo[0]
	//podValue := podinfo[1]
	podName := args[1] //第一个参数是想要调度的pod的name
	SCFilePath := "./SCFile"

	//kubectl get -o yaml pod {podname} > pod.yaml
	//获取当前pod的配置信息，将其输出重定向到pod.yaml中
	cmd1 := exec.Command("cmd", "/c", "kubectl", "get", "-o", "yaml", "pod", podName, ">", "pod.yaml")
	if err := cmd1.Run(); err != nil { // 运行命令
		log.Fatal(err)
	}

	//根据SCFile中的sclang，修改当前的配置文件的affinity，并删除status，生成新的调度文件
	configen.InsertYamlbyTxtstatement(SCFilePath, "./pod.yaml", "./newpod.yaml")
	fmt.Println("config yaml is generated")
	//var affinity template.Affinity
	//matches := configen.ParseStatement("required: app:appA ^ app:appB,app:appB")
	//affinity = configen.InsertMatchRes2PodAffinity(&affinity, matches)
	//fmt.Println(affinity)
	//yamlByte, _ := yaml.Marshal(affinity)
	//
	//if err := os.WriteFile(".\\files\\target2.yaml", yamlByte, 0666); err != nil {
	//	log.Fatal(err)
	//}
}
