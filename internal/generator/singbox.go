package generator

import (
	"encoding/json"

	"github.com/ddranic/ios-proxy-rules/internal/core"
	"github.com/ddranic/ios-proxy-rules/internal/parser"
)

type singboxRule struct {
	Domain        []string `json:"domain,omitempty"`
	DomainSuffix  []string `json:"domain_suffix,omitempty"`
	DomainKeyword []string `json:"domain_keyword,omitempty"`
	DomainRegex   []string `json:"domain_regex,omitempty"`
	IPCidr        []string `json:"ip_cidr,omitempty"`
}

type singboxRuleSet struct {
	Version int           `json:"version"`
	Rules   []singboxRule `json:"rules"`
}

func renderSingboxGeoSite(rules []core.Rule) ([]string, error) {
	var domain, suffix, keyword, regex []string
	for _, r := range rules {
		switch r.Kind {
		case core.KindDomain:
			domain = append(domain, r.Value)
		case core.KindDomainSuffix:
			suffix = append(suffix, r.Value)
		case core.KindDomainKeyword:
			keyword = append(keyword, r.Value)
		case core.KindDomainRegex:
			regex = append(regex, r.Value)
		}
	}

	var ruleObjs []singboxRule
	if len(domain) > 0 {
		ruleObjs = append(ruleObjs, singboxRule{Domain: domain})
	}
	if len(suffix) > 0 {
		ruleObjs = append(ruleObjs, singboxRule{DomainSuffix: suffix})
	}
	if len(keyword) > 0 {
		ruleObjs = append(ruleObjs, singboxRule{DomainKeyword: keyword})
	}
	if len(regex) > 0 {
		ruleObjs = append(ruleObjs, singboxRule{DomainRegex: regex})
	}

	return marshalRuleSet(singboxRuleSet{Version: 2, Rules: ruleObjs})
}

func renderSingboxGeoIP(entry parser.GeoIPEntry) ([]string, error) {
	rule := singboxRule{IPCidr: entry.CIDRs}
	return marshalRuleSet(singboxRuleSet{Version: 2, Rules: []singboxRule{rule}})
}

func marshalRuleSet(rs singboxRuleSet) ([]string, error) {
	b, err := json.MarshalIndent(rs, "", "  ")
	if err != nil {
		return nil, err
	}
	return []string{string(b)}, nil
}
