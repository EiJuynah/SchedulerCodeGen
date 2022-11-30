package conflict_detect

import (
	"CodeGenerationGo/pkg/yaml-process"
	"fmt"
	"testing"
)

func TestPodAffinity2Stringclause(t *testing.T) {
	pod, err := yaml_process.ReadPodYamlFile("E:\\project\\CodeGenerationGo\\files\\out.yaml")
	if err != nil {
		fmt.Println(err)
	}
	cp := PodAffinity2StrClauses(*pod)
	for _, a := range cp {
		for _, b := range a {
			for _, c := range b {
				fmt.Printf(c, ",")
			}
			fmt.Println()
		}
	}

}

func TestStrClauses2CNF(t *testing.T) {
	pod, err := yaml_process.ReadPodYamlFile("E:\\project\\CodeGenerationGo\\files\\out.yaml")
	if err != nil {
		fmt.Println(err)
	}
	cp := PodAffinity2StrClauses(*pod)

	a := StrClauses2CNF(cp)

	fmt.Println(a)
}

func TestCNFExample(t *testing.T) {
	CNFExample()
}

func TestSATPodAffinity(t *testing.T) {
	pod, err := yaml_process.ReadPodYamlFile("E:\\project\\CodeGenerationGo\\files\\out.yaml")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(SATPodAffinity(*pod))
}
