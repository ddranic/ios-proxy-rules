package parser

import "github.com/ddranic/ios-proxy-rules/internal/core"

type Result struct {
	Lists []core.RuleList
}

type state struct {
	visited map[string]struct{}
	seen    map[string]struct{}
}
