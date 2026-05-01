package generator

import (
	"strings"

	"github.com/ddranic/ios-proxy-rules/internal/core"
)

func renderShadowrocketFamily(_ string, rules []core.Rule) ([]string, error) {
	lines := make([]string, 0, len(rules))
	for _, rule := range rules {
		switch rule.Kind {
		case core.KindDomain:
			lines = append(lines, "DOMAIN,"+rule.Value)
		case core.KindDomainSuffix:
			lines = append(lines, "DOMAIN-SUFFIX,"+rule.Value)
		case core.KindDomainKeyword:
			lines = append(lines, "DOMAIN-KEYWORD,"+rule.Value)
		case core.KindDomainRegex:
			lines = append(lines, "URL-REGEX,"+domainRegexToURLRegex(rule.Value))
		}
	}
	return lines, nil
}

func renderClashYAML(_ string, rules []core.Rule) ([]string, error) {
	lines := make([]string, 0, len(rules)+1)
	lines = append(lines, "payload:")
	for _, rule := range rules {
		switch rule.Kind {
		case core.KindDomain:
			lines = append(lines, "  - DOMAIN,"+rule.Value)
		case core.KindDomainSuffix:
			lines = append(lines, "  - DOMAIN-SUFFIX,"+rule.Value)
		case core.KindDomainKeyword:
			lines = append(lines, "  - DOMAIN-KEYWORD,"+rule.Value)
		case core.KindDomainRegex:
			lines = append(lines, "  - DOMAIN-REGEX,"+rule.Value)
		}
	}
	return lines, nil
}

func domainRegexToURLRegex(value string) string {
	pattern := strings.TrimSpace(value)
	startAnchored := strings.HasPrefix(pattern, "^")
	endAnchored := strings.HasSuffix(pattern, "$")

	pattern = strings.TrimPrefix(pattern, "^")
	pattern = strings.TrimSuffix(pattern, "$")

	switch {
	case startAnchored && endAnchored:
		pattern = "^https?://" + pattern
	case startAnchored && !endAnchored:
		pattern = "^https?://" + pattern + `[^/:]*`
	case !startAnchored && endAnchored:
		pattern = "^https?://[^/:]*" + pattern
	default:
		pattern = "^https?://[^/:]*" + pattern + `[^/:]*`
	}

	return pattern + `(?::\d+)?(?:/|$)`
}
