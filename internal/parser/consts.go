package parser

type Prefix string

const (
	PrefixInclude = "include:"
	PrefixFull    = "full:"
	PrefixDomain  = "domain:"
	PrefixKeyword = "keyword:"
	PrefixRegexp  = "regexp:"
)
