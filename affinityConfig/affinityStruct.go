package affinityConfig

type config struct {
	apiVersion string
	kind       string
	metadata   Metadata
	spec       Spec
}

type Affinity struct {
	podAffinity     PodAffinity
	podAntiAffinity PodAntiAffinity
}

type Container struct {
	name  string
	image string
}

type LabelSelector struct {
	matchExpressions []MatchExpression
}

type MatchExpression struct {
	key      string
	operator string
	values   []string
}

type Metadata struct {
	name  string
	label map[string]string
}

type PodAffinity struct {
	requiredDuringSchedulingIgnoredDuringExecution  RequiredDuringSchedulingIgnoredDuringExecution
	preferredDuringSchedulingIgnoredDuringExecution PreferredDuringSchedulingIgnoredDuringExecution
}

type PodAntiAffinity struct {
}

type PodAffinityTerm struct {
	labelSelector []LabelSelector
}

type Perference struct {
	weight          int
	podAffinityTerm PodAffinityTerm
}

type PreferredDuringSchedulingIgnoredDuringExecution struct {
	preference []Perference
}

type RequiredDuringSchedulingIgnoredDuringExecution struct {
	labelSelector []LabelSelector
	topologyKey   string
}

type Spec struct {
	affinity   Affinity
	containers []Container
}
