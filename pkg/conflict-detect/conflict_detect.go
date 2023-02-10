package conflict_detect

import (
	"CodeGenerationGo/pkg/template"
	"github.com/mitchellh/go-sat"
	"github.com/mitchellh/go-sat/cnf"
	"github.com/mitchellh/go-sat/dimacs"
	"strconv"
)

//输入为kubernetes pod对象
//@title
//@description 将pod对象中的affinity约束，转换成形式化的约束
//@auth

func PodAffinity2StrClauses(pod template.Pod) [][][]string {
	//name := pod.ObjectMeta.Name
	//

	var clause [][][]string
	//针对PodAffinity
	for _, require := range pod.Spec.Affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution {
		//require.LabelSelector.MatchLabels
		var subclause [][]string

		//分析MatchLabels
		//MatchLabel等与opt为In的MatchExpresssions,而且其中的values的数量只有一个
		for k, v := range require.LabelSelector.MatchLabels {
			key := k
			values := v
			subsubclause := []string{}

			subsubclause = append(subsubclause, "1")

			str := key + ":" + values
			subsubclause = append(subsubclause, str)

			subclause = append(subclause, subsubclause)

		}

		//分析MatchExpressions
		for _, expression := range require.LabelSelector.MatchExpressions {

			key := expression.Key
			opt := expression.Operator
			values := expression.Values
			subsubclause := []string{}

			//字句中第一个字符为符号位，若为notin，则取非，为-1，若为in，则取正，为1
			if opt == template.LabelSelectorOpNotIn {
				subsubclause = append(subsubclause, "-1")
				for _, str := range values {
					str = key + ":" + str
					subsubclause = append(subsubclause, str)
				}

			} else if opt == template.LabelSelectorOpIn {
				subsubclause = append(subsubclause, "1")
				for _, str := range values {
					str = key + ":" + str
					subsubclause = append(subsubclause, str)
				}

			}

			subclause = append(subclause, subsubclause)

		}
		clause = append(clause, subclause)
	}

	//针对antiaffinity
	for _, require := range pod.Spec.Affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution {
		//require.LabelSelector.MatchLabels
		var subclause [][]string
		for _, expression := range require.LabelSelector.MatchExpressions {

			key := expression.Key
			opt := expression.Operator
			values := expression.Values
			subsubclause := []string{}

			//antiaffinity与affinity相反
			//字句中第一个字符为符号位，若为notin，则取非的非，为1，若为in，则取非，为-1
			if opt == template.LabelSelectorOpNotIn {
				subsubclause = append(subsubclause, "1")
				for _, str := range values {
					str = key + ":" + str
					subsubclause = append(subsubclause, str)
				}

			} else if opt == template.LabelSelectorOpIn {
				subsubclause = append(subsubclause, "-1")
				for _, str := range values {
					str = key + ":" + str
					subsubclause = append(subsubclause, str)
				}

			}

			subclause = append(subclause, subsubclause)

		}
		clause = append(clause, subclause)
	}

	return clause

}

// return：
// 1：clause集formula
// 2：str与index的对应map
func StrClauses2CNF(strclauses [][][]string) (cnf.Formula, map[int]string) {
	var formulaIntMap [][]int
	//该map存储将value与x1,x2,x3...等cnf中的变量的映射关系
	valueName2LitMap := make(map[string]int)
	for _, scheduleState := range strclauses {
		for _, matchExpression := range scheduleState {
			//切片第一为符号位
			sign_bit := matchExpression[0]
			valueName2LitMap["sign_bit"], _ = strconv.Atoi(sign_bit)
			//将每句MacthExpression映射
			for i := 1; i < len(matchExpression); i++ {

				value := matchExpression[i]
				_, ok := valueName2LitMap[value]
				if !ok {
					valueName2LitMap[value] = len(valueName2LitMap)
				}
			}
			var clauseInt []int
			for i := 1; i < len(matchExpression); i++ {
				clauseInt = append(clauseInt, valueName2LitMap[matchExpression[i]])
			}

			if sign_bit == "1" {
				formulaIntMap = append(formulaIntMap, clauseInt)
			}

			//若为NotIn，根据德摩根定律将非结合进去。 ^(A | B) = ^A & ^B
			if sign_bit == "-1" {
				for i := 0; i < len(clauseInt); i++ {
					clauseInt[i] = -clauseInt[i]
				}
				for _, ci := range clauseInt {
					newClause := []int{ci}
					formulaIntMap = append(formulaIntMap, newClause)
				}

			}

		}
	}

	//获得的x1，x2与对应字符串的key，vale反转字典，用与后续定位
	Lit2StrMap := map[int]string{}
	for k, v := range valueName2LitMap {
		Lit2StrMap[v] = k
	}

	formula := cnf.NewFormulaFromInts(formulaIntMap)

	return formula, Lit2StrMap
}

func CNF2Dimacs(formula cnf.Formula) dimacs.Problem {
	set := make(map[cnf.Lit]bool)
	for _, clause := range formula {
		for _, lit := range clause {
			set[lit] = true
		}
	}
	variables := len(set)
	clauses := len(formula)
	problem := dimacs.Problem{
		Formula:   formula,
		Variables: variables,
		Clauses:   clauses}

	return problem
}

//input:mitchellh/go-sat提供的cnf.Formula
// 使用go提供的sat solver去求解满足性问题
//当前的求解器是go提供的mitchellh/go-sat
//求解器可替换
//

func CNFSolve(formula cnf.Formula) (bool, []int) {
	solver := sat.New()
	//开启日志，从日志中获取信息
	//Trace和Tracer两个参数需要设置
	//tracer使用log的Printf
	solver.Trace = true
	//sslogger := log.New(os.Stdout, "SAT solver Log:", 0)
	sslogger := &SatSolverLogger{}
	solver.Tracer = sslogger
	solver.AddFormula(formula)

	ifsat := solver.Solve()
	conflictClauseIndexs := []int{}
	//如果有冲突，输出有冲突的标号
	if ifsat == false {
		conflictClauseIndexs = AnalyseLogs(sslogger.logs)
	}

	return ifsat, conflictClauseIndexs
}

func SATPodAffinity(pod template.Pod) (bool, []string) {
	strclauses := PodAffinity2StrClauses(pod)
	problem, clauseMap := StrClauses2CNF(strclauses)

	ifsat, conflictClauseIndexs := CNFSolve(problem)
	if ifsat {
		return true, nil
	} else {
		//fmt.Println(clauseMap)
		conflictlabels := ConflictLocate(clauseMap, conflictClauseIndexs)

		return false, conflictlabels

	}
}
