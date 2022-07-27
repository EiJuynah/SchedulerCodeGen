package util

func Test() {

}

func readInputTxt() {

}

func AnalysisRelationalStatement(statement string) {

}

func Relation2Opera(relat string) string {
	if relat == "&" {
		return "In"
	}
	if relat == "^" {
		return "NotIn"
	}
	return "nil"
}
