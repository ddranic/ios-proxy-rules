package generator

import (
	"fmt"
	"strings"

	"github.com/ddranic/ios-proxy-rules/internal/parser"
)

func renderGeoIP(platform Platform, entry parser.GeoIPEntry) ([]string, error) {
	switch platform {
	case PlatformShadowrocket, PlatformLoon, PlatformSurge:
		return renderGeoIPList(entry.CIDRs), nil
	case PlatformClash:
		return renderGeoIPClash(entry.CIDRs), nil
	case PlatformSingbox:
		return renderSingboxGeoIP(entry)
	default:
		return nil, fmt.Errorf("unsupported platform %q", platform)
	}
}

func renderGeoIPList(cidrs []string) []string {
	lines := make([]string, 0, len(cidrs))
	for _, cidr := range cidrs {
		if isIPv6(cidr) {
			lines = append(lines, "IP-CIDR6,"+cidr)
		} else {
			lines = append(lines, "IP-CIDR,"+cidr)
		}
	}
	return lines
}

func renderGeoIPClash(cidrs []string) []string {
	lines := make([]string, 0, len(cidrs)+1)
	lines = append(lines, "payload:")
	for _, cidr := range cidrs {
		if isIPv6(cidr) {
			lines = append(lines, "  - IP-CIDR6,"+cidr)
		} else {
			lines = append(lines, "  - IP-CIDR,"+cidr)
		}
	}
	return lines
}

func isIPv6(cidr string) bool {
	return strings.ContainsRune(cidr, ':')
}
