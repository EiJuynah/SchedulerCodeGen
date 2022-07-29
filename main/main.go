package main

import (
	"CodeGenerationGo/configen"
)

func main() {
	sourceyamlpath := ".\\files\\source.yaml"
	outputyamlpath := ".\\files\\target.yaml"
	configpath := ".\\files\\input.txt"
	configen.YamlGenbyTxt(configpath, sourceyamlpath, outputyamlpath)
}
