package configen

import (
	"CodeGenerationGo/util"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func GetPodYaml(podName string) ([]byte, []byte) {
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

func AddAffinityByPodname(podName string, SCFilePath string) {
	GetPodYaml(podName)
	DeletePodStatusFromYaml(".\\pod.yaml", ".\\temp.yaml")
	util.DeleteFile(".\\pod.yaml")
	InsertYamlbyTxtstatement(SCFilePath, ".\\temp.yaml", ".\\newpod.yaml")
	util.DeleteFile(".\\temp.yaml")
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
