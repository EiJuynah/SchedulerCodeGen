package template

//
//import (
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	api "k8s.io/kubernetes/pkg/apis/core/v1"
//)
//
//type Affinity struct {
//	NodeAffinity    *NodeAffinity    `json:"nodeAffinity,omitempty"`
//	PodAffinity     *PodAffinity     `json:"podAffinity,omitempty"`
//	PodAntiAffinity *PodAntiAffinity `json:"podAntiAffinity,omitempty"`
//}
//
//type NodeAffinity struct {
//	// If the affinity requirements specified by this field are not met at
//	// scheduling time, the pod will not be scheduled onto the node.
//	// If the affinity requirements specified by this field cease to be met
//	// at some point during pod execution (e.g. due to a node label update),
//	// the system will try to eventually evict the pod from its node.
//	RequiredDuringSchedulingRequiredDuringExecution *NodeSelector `json:"requiredDuringSchedulingRequiredDuringExecution,omitempty"`
//	// If the affinity requirements specified by this field are not met at
//	// scheduling time, the pod will not be scheduled onto the node.
//	// If the affinity requirements specified by this field cease to be met
//	// at some point during pod execution (e.g. due to a node label update),
//	// the system may or may not try to eventually evict the pod from its node.
//	RequiredDuringSchedulingIgnoredDuringExecution *NodeSelector `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
//	// The scheduler will prefer to schedule pods to nodes that satisfy
//	// the affinity expressions specified by this field, but it may choose
//	// a node that violates one or more of the expressions. The node that is
//	// most preferred is the one with the greatest sum of weights, i.e.
//	// for each node that meets all of the scheduling requirements (resource
//	// request, RequiredDuringScheduling affinity expressions, etc.),
//	// compute a sum by iterating through the elements of this field and adding
//	// "weight" to the sum if the node matches the corresponding MatchExpressions; the
//	// node(s) with the highest sum are the most preferred.
//	PreferredDuringSchedulingIgnoredDuringExecution []PreferredSchedulingTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`
//}
//
//// An empty preferred scheduling term matches all objects with implicit weight 0
//// (i.e. it's a no-op). A null preferred scheduling term matches no objects.
//type PreferredSchedulingTerm struct {
//	// weight is in the range 1-100
//	Weight int `json:"weight"`
//	// matchExpressions is a list of node selector requirements. The requirements are ANDed.
//	MatchExpressions []NodeSelectorRequirement `json:"matchExpressions,omitempty"`
//}
//
//type PodAffinity struct {
//	// If the affinity requirements specified by this field are not met at
//	// scheduling time, the pod will not be scheduled onto the node.
//	// If the affinity requirements specified by this field cease to be met
//	// at some point during pod execution (e.g. due to a pod label update), the
//	// system will try to eventually evict the pod from its node.
//	// When there are multiple elements, the lists of nodes corresponding to each
//	// PodAffinityTerm are intersected, i.e. all terms must be satisfied.
//	RequiredDuringSchedulingRequiredDuringExecution []PodAffinityTerm `json:"requiredDuringSchedulingRequiredDuringExecution,omitempty"`
//	// If the affinity requirements specified by this field are not met at
//	// scheduling time, the pod will not be scheduled onto the node.
//	// If the affinity requirements specified by this field cease to be met
//	// at some point during pod execution (e.g. due to a pod label update), the
//	// system may or may not try to eventually evict the pod from its node.
//	// When there are multiple elements, the lists of nodes corresponding to each
//	// PodAffinityTerm are intersected, i.e. all terms must be satisfied.
//	RequiredDuringSchedulingIgnoredDuringExecution []PodAffinityTerm `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
//	// The scheduler will prefer to schedule pods to nodes that satisfy
//	// the affinity expressions specified by this field, but it may choose
//	// a node that violates one or more of the expressions. The node that is
//	// most preferred is the one with the greatest sum of weights, i.e.
//	// for each node that meets all of the scheduling requirements (resource
//	// request, RequiredDuringScheduling affinity expressions, etc.),
//	// compute a sum by iterating through the elements of this field and adding
//	// "weight" to the sum if the node matches the corresponding MatchExpressions; the
//	// node(s) with the highest sum are the most preferred.
//	PreferredDuringSchedulingIgnoredDuringExecution []WeightedPodAffinityTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`
//}
//
//type PodAntiAffinity struct {
//	// If the anti-affinity requirements specified by this field are not met at
//	// scheduling time, the pod will not be scheduled onto the node.
//	// If the anti-affinity requirements specified by this field cease to be met
//	// at some point during pod execution (e.g. due to a pod label update), the
//	// system will try to eventually evict the pod from its node.
//	// When there are multiple elements, the lists of nodes corresponding to each
//	// PodAffinityTerm are intersected, i.e. all terms must be satisfied.
//	RequiredDuringSchedulingRequiredDuringExecution []PodAffinityTerm `json:"requiredDuringSchedulingRequiredDuringExecution,omitempty"`
//	// If the anti-affinity requirements specified by this field are not met at
//	// scheduling time, the pod will not be scheduled onto the node.
//	// If the anti-affinity requirements specified by this field cease to be met
//	// at some point during pod execution (e.g. due to a pod label update), the
//	// system may or may not try to eventually evict the pod from its node.
//	// When there are multiple elements, the lists of nodes corresponding to each
//	// PodAffinityTerm are intersected, i.e. all terms must be satisfied.
//	RequiredDuringSchedulingIgnoredDuringExecution []PodAffinityTerm `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
//	// The scheduler will prefer to schedule pods to nodes that satisfy
//	// the anti-affinity expressions specified by this field, but it may choose
//	// a node that violates one or more of the expressions. The node that is
//	// most preferred is the one with the greatest sum of weights, i.e.
//	// for each node that meets all of the scheduling requirements (resource
//	// request, RequiredDuringScheduling anti-affinity expressions, etc.),
//	// compute a sum by iterating through the elements of this field and adding
//	// "weight" to the sum if the node matches the corresponding MatchExpressions; the
//	// node(s) with the highest sum are the most preferred.
//	PreferredDuringSchedulingIgnoredDuringExecution []WeightedPodAffinityTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`
//}
//
//type WeightedPodAffinityTerm struct {
//	// weight is in the range 1-100
//	Weight          int             `json:"weight"`
//	PodAffinityTerm PodAffinityTerm `json:"podAffinityTerm"`
//}
//
//type PodAffinityTerm struct {
//	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty"`
//	// namespaces specifies which namespaces the LabelSelector applies to (matches against);
//	// nil list means "this pod's namespace," empty list means "all namespaces"
//	// The json tag here is not "omitempty" since we need to distinguish nil and empty.
//	// See https://golang.org/pkg/encoding/json/#Marshal for more details.
//	Namespaces []api.Namespace `json:"namespaces,omitempty"`
//	// empty topology key is interpreted by the scheduler as "all topologies"
//	TopologyKey string `json:"topologyKey,omitempty"`
//}
