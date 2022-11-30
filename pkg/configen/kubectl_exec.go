package configen

import (
	"CodeGenerationGo/pkg/statement-analysis"
	"CodeGenerationGo/pkg/yaml-process"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func GetPodYamlbyName(podName string) ([]byte, []byte) {
	//获取当前pod的配置信息，将其输出重定向到pod.yaml中
	//kubectl get -o yaml pod {podname} > pod.yaml
	//此时的commmand为Windows的cmd，如果是linux环境，换成 bin/bash

	cmd := exec.Command("cmd", "/c", "kubectl", "get", "-o", "yaml", "pod", podName, ">", "pod.yaml")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return stdout.Bytes(), stderr.Bytes()
}

func GetPodYamlbylabel(labelpair string) ([]byte, []byte) {
	//获取当前pod的配置信息，将其输出重定向到pod.yaml中
	//kubectl get -o yaml pod {podname} > pod.yaml
	//此时的commmand为Windows的cmd，如果是linux环境，换成 bin/bash
	label := statement_analysis.SplitStringByOP(labelpair, ":")
	key := label[0]
	value := label[1]
	pair := key + "=" + value
	cmd := exec.Command("cmd", "/c", "kubectl", "get", "-o", "yaml", "pod", "-l", pair, ">", "pod.yaml")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return stdout.Bytes(), stderr.Bytes()
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
	InsertYamlbyTxtstatement(SCFilePath, "./pod.yaml", "./newpod.yaml")

	//forth step
	//根据新生成的调度文件，重新启动一个新的pod
	//kubectl apply -f ./newpod.yaml
	cmd := exec.Command("cmd", "/c", "kubectl", "apply", "-f", "./newpod.yaml")
	if err := cmd.Run(); err != nil { // 运行命令
		log.Fatal(err)
	}

	fmt.Println("config yaml is generated")

}

func AddAffinityByPodname(podName string, SCFilePath string) {
	GetPodYamlbyName(podName)
	DeletePodStatusFromYaml(".\\pod.yaml", ".\\temp.yaml")
	yaml_process.DeleteFile(".\\pod.yaml")
	InsertYamlbyTxtstatement(SCFilePath, ".\\temp.yaml", ".\\newpod.yaml")
	yaml_process.DeleteFile(".\\temp.yaml")
}

func AddAffinityByLabel(labelpair string, SCFilePath string) {
	GetPodYamlbylabel(labelpair)
	DeletePodStatusFromYaml(".\\pod.yaml", ".\\temp.yaml")
	yaml_process.DeleteFile(".\\pod.yaml")
	InsertYamlbyTxtstatement(SCFilePath, ".\\temp.yaml", ".\\newpod.yaml")
	yaml_process.DeleteFile(".\\temp.yaml")
}

func DeletePod(podName string) {
	//second step
	//删除当前的pod
	//kubectl delete pod <podname> -n <namespace>
	cmd := exec.Command("cmd", "/c", "kubectl", "delete", "pod", podName)
	if err := cmd.Run(); err != nil { // 运行命令
		log.Fatal(err)
	}

}

func GenNewYaml(SCFilePath string) {
	InsertYamlbyTxtstatement(SCFilePath, "./pod.yaml", "./newpod.yaml")
}

func StartNewPod() {
	//根据新生成的调度文件，重新启动一个新的pod
	//kubectl apply -f ./newpod.yaml
	cmd := exec.Command("cmd", "/c", "kubectl", "apply", "-f", "./newpod.yaml")
	if err := cmd.Run(); err != nil { // 运行命令
		log.Fatal(err)
	}
}

func runInLinux(cmd string) string {
	fmt.Println("Running Linux cmd:", cmd)
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	return strings.TrimSpace(string(result))
}

func runInWindows(cmd string) string {
	fmt.Println("Running Win cmd:", cmd)
	result, err := exec.Command("cmd", "/c", cmd).Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	return strings.TrimSpace(string(result))
}

func RunCommand(cmd string) string {
	if runtime.GOOS == "windows" {
		return runInWindows(cmd)
	} else {
		return runInLinux(cmd)
	}
}
