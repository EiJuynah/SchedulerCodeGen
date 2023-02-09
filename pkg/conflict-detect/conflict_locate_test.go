package conflict_detect

import (
	"fmt"
	"testing"
)

func TestAnalyseLogs(t *testing.T) {
	strs := []string{"[TRACE] sat: addClause: not adding literal; literal -1 false: [-1]"}
	AnalyseLogs(strs)
}

func TestGetConflictLitIndex(t *testing.T) {
	strs := "[TRACE] sat: addClause: not adding literal; literal -4 false: [-1]"
	res := getConflictLitIndex(strs)
	fmt.Println(res)
}
