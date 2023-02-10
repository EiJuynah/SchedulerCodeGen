package conflict_detect

import (
	"CodeGenerationGo/pkg/yaml-process"
	"fmt"
	"testing"
)

func TestPodAffinity2Stringclause(t *testing.T) {
	pjtpath := "E://project//"
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
	pjtpath := "D://code//GO//"
	pod, err := yaml_process.ReadPodYamlFile(pjtpath + "CodeGenerationGo\\files\\out.yaml")
	if err != nil {
		fmt.Println(err)
	}
	cp := PodAffinity2StrClauses(*pod)

	a, fmap := StrClauses2CNF(cp)

	fmt.Println(a)

	fmt.Println(fmap)
}

func TestSATPodAffinity(t *testing.T) {
	pjtpath := "D://code//GO//"
	pod, err := yaml_process.ReadPodYamlFile(pjtpath + "CodeGenerationGo\\files\\out.yaml")
	if err != nil {
		fmt.Println(err)
	}
	res, cfs := SATPodAffinity(*pod)

	if res {
		fmt.Println("the result of conflict detect is : ", res)
	} else {
		fmt.Println("the result of conflict detect is : ", res)
		fmt.Println("the label which maybe has conflict is :", cfs)
	}

}
