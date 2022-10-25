package main

import (
	"CodeGenerationGo/pkg/configen"
	"fmt"
)

func main() {
	a := configen.SplitStringByOP("aaaa:bbbb", ":")
	fmt.Println(a)
}
