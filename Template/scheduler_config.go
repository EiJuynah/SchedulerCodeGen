package Template

type Config struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type Affinity struct {
	PodAffinity     PodAffinity     `yaml:"podAffinity"`
	PodAntiAffinity PodAntiAffinity `yaml:"podAntiAffinity"`
}

type Container struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}

func GetNewContainer(name string, image string) Container {
	var container Container
	container.Name = name
	container.Image = image

	return container
}

type LabelSelector struct {
	MatchExpressions []MatchExpression `yaml:"matchExpressions"`
}

type MatchExpression struct {
	Key      string   `yaml:"key"`
	Operator string   `yaml:"operator"`
	Values   []string `yaml:"values"`
}

type Metadata struct {
	Name  string            `yaml:"name"`
	Label map[string]string `yaml:"label"`
}

type PodAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution  RequiredDuringSchedulingIgnoredDuringExecution  `yaml:"requiredDuringSchedulingIgnoredDuringExecution"`
	PreferredDuringSchedulingIgnoredDuringExecution PreferredDuringSchedulingIgnoredDuringExecution `yaml:"preferredDuringSchedulingIgnoredDuringExecution"`
}

type PodAntiAffinity struct {
}

type PodAffinityTerm struct {
	LabelSelector []LabelSelector `yaml:"labelSelector"`
}

type Perference struct {
	Weight          int             `yaml:"weight"`
	PodAffinityTerm PodAffinityTerm `yaml:"podAffinityTerm"`
}

type PreferredDuringSchedulingIgnoredDuringExecution struct {
	Preference []Perference `yaml:"preference"`
}

type RequiredDuringSchedulingIgnoredDuringExecution struct {
	LabelSelector []LabelSelector `yaml:"labelSelector"`
	TopologyKey   string          `yaml:"topologyKey"`
}

type Spec struct {
	Affinity   Affinity    `yaml:"affinity"`
	Containers []Container `yaml:"containers"`
}

type MatchRes struct {
	Trendrule        string
	Weight           int
	Relationship     string
	Key              string
	Value            string
	MatchExpressions []MatchExpression
}
