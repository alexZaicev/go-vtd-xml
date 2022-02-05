package pilot

type IterationType int

const (
	Undefined IterationType = iota
	Simple
	SimpleNs
	Descending
	DescendingNs
	Following
	FollowingNs
	Preceding
	PrecedingNs
	Attr
	AttrNs
	Namespace
	SimpleNode
	DescendantNode
	FollowingNode
	PrecedingNode
)
