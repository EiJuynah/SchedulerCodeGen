package util

import (
	"CodeGenerationGo/template"
)

func Test() {

}

func readInputTxt() {

}

func AnalysisRelationalStatement(statement string) {

}

func Relation2Opera(relat string) template.LabelSelectorOperator {
	if relat == "&" {
		return template.LabelSelectorOpIn
	}
	if relat == "^" {
		return template.LabelSelectorOpNotIn
	}
	return "nil"
}

//func ContainsInArray[T any](items []T, item T) bool {
//	for _, eachItem := range items {
//		bytes.Compare(eachItem, item)
//
//	}
//	return false
//}
