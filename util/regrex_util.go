package util

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test() {

}

func readInputTxt() {

}

func AnalysisRelationalStatement(statement string) {

}

func Relation2Opera(relat string) metav1.LabelSelectorOperator {
	if relat == "&" {
		return metav1.LabelSelectorOpIn
	}
	if relat == "^" {
		return metav1.LabelSelectorOpNotIn
	}
	return "nil"
}
