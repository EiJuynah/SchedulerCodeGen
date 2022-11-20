package conflict_detect

import (
	"CodeGenerationGo/pkg/util"
	"fmt"
	"testing"
)

func TestPodAffinity2Stringclause(t *testing.T) {
	pod, err := util.ReadPodYamlFile("E:\\project\\CodeGenerationGo\\files\\out.yaml")
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
