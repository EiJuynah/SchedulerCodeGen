package template

type MatchRes struct {
	Trendrule        string
	Weight           int
	Relationship     LabelSelectorOperator
	LabelKey         string
	Value            string
	MatchExpressions []LabelSelectorRequirement
}

const (
	DEFAULT_TOPOLOGYKRY = "kubernetes.io/hostname"
)
