package conflict_detect

import (
	"fmt"
	"regexp"
	"strconv"
)

//定位冲突的位置
//定位展示:定位出有冲突的若干个冲突对的位置
//位置表示方法:行号+行数

//SatSolverLogger是为了适配sat.Solver.Tracer接口的实现
//将输出的log信息，保存在内存中logs里面。
//下一步直接分析logs的内容，得到有冲突的语句
type SatSolverLogger struct {
	logs []string
}

func (ssl *SatSolverLogger) PrintLogs() {
	for _, str := range ssl.logs {
		fmt.Println(str)
	}
}

func (ssl *SatSolverLogger) Printf(format string, v ...any) {

	s := fmt.Sprintf(format, v...)
	ssl.logs = append(ssl.logs, s)

}

//扫描每一条log，输出有冲突的Lit的index
func AnalyseLogs(logs []string) []int {
	res := []int{}
	for _, log := range logs {
		index := getConflictLitIndex(log)
		//log中存在冲突的时候，返回不为-1
		if index != -1 {
			res = append(res, index)
		}
	}

	return res

}

//判断该条log是否有冲突，如果有冲突，则输出有冲突的lit的index
//return : 返回有冲突的Lit的标号
// return:没有相应的匹配则返回-1
func getConflictLitIndex(log string) int {
	//针对sat: addClause: not adding literal; literal * false: [*]
	regstr := ".literal -*(\\d) false: \\[(-*\\d)\\]"
	reg_false := regexp.MustCompile(regstr)

	if reg_false == nil {
		fmt.Println("regexp err")
		return -1
	}

	//根据规则提取关键信息
	result1 := reg_false.FindAllStringSubmatch(log, -1)

	if result1 == nil {
		return -1
	}
	//current 有冲突的 lit标号
	conflictLitIndexStr := result1[0][1]

	conflictLitIndex, err := strconv.Atoi(conflictLitIndexStr)
	if err != nil {
		fmt.Println("string output err")
		return -1
	}
	return conflictLitIndex
}

func ConflictLocate(clauseMap map[int]string, conflictClauseIndexs []int) []string {
	labelpair := []string{}
	for _, index := range conflictClauseIndexs {
		labelpair = append(labelpair, clauseMap[index])
	}

	return labelpair

}
