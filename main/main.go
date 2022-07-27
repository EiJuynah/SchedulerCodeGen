package main

import (
	"CodeGenerationGo/codegenerator"
	"fmt"
)

func main() {
	a := codegenerator.ParseStatement("preferred:100 app:appA & apq:apgB,apv:appC")
	fmt.Println(a)
}
