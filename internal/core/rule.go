package core

type RuleKind string

const (
	KindDomain        RuleKind = "domain"
	KindDomainSuffix  RuleKind = "domain-suffix"
	KindDomainKeyword RuleKind = "domain-keyword"
	KindDomainRegex   RuleKind = "domain-regex"
)

type Rule struct {
	Kind  RuleKind
	Value string
}

type RuleList struct {
	Name  string
	Rules []Rule
}
