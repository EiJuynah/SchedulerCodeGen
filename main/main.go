package main

import (
	"CodeGenerationGo/codegenerator"
	"CodeGenerationGo/util"
	"fmt"
)

func main() {
	match1 := codegenerator.ParseStatement("preferred:100 app:appA & app:apgB,app:appC")
	match2 := codegenerator.ParseStatement("required: app:appA ^ app:apgD,app:appE")
	fmt.Println(match1)
	affinity := codegenerator.AffinityInit()
	affinity = codegenerator.InsertMatchRes2Affinity(affinity, match1)
	affinity = codegenerator.InsertMatchRes2Affinity(affinity, match2)
	fmt.Println(affinity)

	config, _ := util.ReadConfigYaml("D:\\学习\\华师大软件优化\\k8s_ali\\CodeGenerationGo\\files\\test4.yaml")
	fmt.Println(config)

	config.Spec.Affinity = affinity
	fmt.Println(affinity)
}
