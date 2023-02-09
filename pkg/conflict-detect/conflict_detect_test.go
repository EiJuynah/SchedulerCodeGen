package conflict_detect

import (
	"CodeGenerationGo/pkg/yaml-process"
	"fmt"
	"testing"
)

func TestPodAffinity2Stringclause(t *testing.T) {
	pjtpath := "D:\\code\\GO\\"
	pod, err := yaml_process.ReadPodYamlFile(pjtpath + "CodeGenerationGo\\files\\out.yaml")
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
	pjtpath := "D:\\code\\GO\\"
	pod, err := yaml_process.ReadPodYamlFile(pjtpath + "CodeGenerationGo\\files\\out.yaml")
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
	pjtpath := "D:\\code\\GO\\"
	pod, err := yaml_process.ReadPodYamlFile(pjtpath + "CodeGenerationGo\\files\\out.yaml")
	if err != nil {
		fmt.Println(err)
	}
	res := SATPodAffinity(*pod)
	fmt.Println(res)
}

func TestAbaAba(t *testing.T) {
	a := make(map[string]string)
	a["e"] = "r"
	a["c"] = "C"
	a["x"] = "X"
	for k, v := range a {
		fmt.Println(k)
		fmt.Println(v)
	}

}
