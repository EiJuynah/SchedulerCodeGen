package conflict_detect

import (
	"CodeGenerationGo/pkg/template"
	"fmt"
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
	//针对affinity
	for _, require := range pod.Spec.Affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution {
		//require.LabelSelector.MatchLabels
		var subclause [][]string
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

func CNFExample() {
	// Consider the example formula already in CNF.
	//
	// ( ¬x1 ∨ ¬x3 ∨ ¬x4 ) ∧ ( x2 ∨ x3 ∨ ¬x4 ) ∧
	// ( x1 ∨ ¬x2 ∨ x4 ) ∧ ( x1 ∨ x3 ∨ x4 ) ∧ ( ¬x1 ∨ x2 ∨ ¬x3 )
	// (¬x4)

	// Represent each variable as an int where a negative value means negated
	formula := cnf.NewFormulaFromInts([][]int{
		[]int{-1, -3, -4},
		[]int{2, 3, -4},
		[]int{1, -2, 4},
		[]int{1, 3, 4},
		[]int{-1, 2, -3},
		[]int{-4},
	})

	// Create a solver and add the formulas to solve
	s := sat.New()
	s.AddFormula(formula)

	// For low level details on how a solution came to be:
	// s.Trace = true
	// s.Tracer = log.New(os.Stderr, "", log.LstdFlags)

	// Solve it!
	sat := s.Solve()

	// Solution can be read from Assignments. The key is the variable
	// (always positive) and the value is the value.
	solution := s.Assignments()

	fmt.Printf("Solved: %v\n", sat)
	fmt.Printf("Solution:\n")
	fmt.Printf("  x1 = %v\n", solution[1])
	fmt.Printf("  x2 = %v\n", solution[2])
	fmt.Printf("  x3 = %v\n", solution[3])
	fmt.Printf("  x4 = %v\n", solution[4])
	// Output:
	// Solved: true
	// Solution:
	//   x1 = true
	//   x2 = true
	//   x3 = true
	//   x4 = false
}

func StrClauses2CNF(strclauses [][][]string) cnf.Formula {
	var formulaInt [][]int
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
			if sign_bit == "-1" {
				for i := 0; i < len(clauseInt); i++ {
					clauseInt[i] = -clauseInt[i]
				}
			}

			formulaInt = append(formulaInt, clauseInt)
		}
	}

	formula := cnf.NewFormulaFromInts(formulaInt)

	return formula
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

func CNFSolve(formula cnf.Formula) bool {
	solver := sat.New()
	solver.AddFormula(formula)

	return solver.Solve()
}

func SATPodAffinity(pod template.Pod) bool {
	strclauses := PodAffinity2StrClauses(pod)
	problem := StrClauses2CNF(strclauses)

	return CNFSolve(problem)
}
