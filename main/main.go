package main

import (
	"CodeGenerationGo/configen"
	"CodeGenerationGo/template"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func main() {
	//sourceyamlpath := ".\\files\\source.yaml"
	//outputyamlpath := ".\\files\\target.yaml"
	////configpath := ".\\files\\input.txt"
	//configen.InsertAffinity2Yaml("required: app:appA & app:appB", sourceyamlpath, outputyamlpath)
	var affinity template.Affinity
	matches := configen.ParseStatement("required: app:appA ^ app:appB")
	affinity = configen.InsertMatchRes2PodAffinity(&affinity, matches)
	fmt.Println(affinity)
	yamlByte, _ := yaml.Marshal(affinity)

	if err := os.WriteFile(".\\files\\target2.yaml", yamlByte, 0666); err != nil {
		log.Fatal(err)
	}

}
