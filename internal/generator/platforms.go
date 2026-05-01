package generator

import "github.com/ddranic/ios-proxy-rules/internal/core"

func configForPlatform(platform Platform) (string, func(string, []core.Rule) ([]string, error), bool) {
	switch platform {
	case PlatformShadowrocket, PlatformLoon, PlatformSurge, PlatformStash:
		return extList, renderShadowrocketFamily, true
	case PlatformClash:
		return extYAML, renderClashYAML, true
	default:
		return "", nil, false
	}
}
