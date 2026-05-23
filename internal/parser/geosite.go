package parser

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"google.golang.org/protobuf/encoding/protowire"

	"github.com/ddranic/ios-proxy-rules/internal/core"
)

// dlc.dat domain type constants
const (
	domainTypePlain  = 0 // keyword
	domainTypeRegex  = 1 // regex
	domainTypeDomain = 2 // suffix
	domainTypeFull   = 3 // exact
)

// ParseGeoSite reads a v2fly dlc.dat and returns rule lists sorted by name.
func ParseGeoSite(path string) ([]core.RuleList, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", path, err)
	}
	lists, err := parseGeoSiteList(data)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", path, err)
	}
	sort.Slice(lists, func(i, j int) bool {
		return lists[i].Name < lists[j].Name
	})
	return lists, nil
}

func parseGeoSiteList(data []byte) ([]core.RuleList, error) {
	var lists []core.RuleList
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return nil, fmt.Errorf("bad tag")
		}
		data = data[n:]
		if num == 1 && typ == protowire.BytesType {
			b, n := protowire.ConsumeBytes(data)
			if n < 0 {
				return nil, fmt.Errorf("bad GeoSite message")
			}
			data = data[n:]
			list, err := parseGeoSiteEntry(b)
			if err != nil {
				return nil, err
			}
			if len(list.Rules) > 0 {
				lists = append(lists, list)
			}
		} else {
			n := protowire.ConsumeFieldValue(num, typ, data)
			if n < 0 {
				return nil, fmt.Errorf("bad field")
			}
			data = data[n:]
		}
	}
	return lists, nil
}

func parseGeoSiteEntry(data []byte) (core.RuleList, error) {
	var list core.RuleList
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return list, fmt.Errorf("bad tag in GeoSite")
		}
		data = data[n:]
		switch {
		case num == 1 && typ == protowire.BytesType: // country_code
			b, n := protowire.ConsumeBytes(data)
			if n < 0 {
				return list, fmt.Errorf("bad country_code")
			}
			data = data[n:]
			list.Name = strings.ToLower(string(b))
		case num == 2 && typ == protowire.BytesType: // domain
			b, n := protowire.ConsumeBytes(data)
			if n < 0 {
				return list, fmt.Errorf("bad Domain message")
			}
			data = data[n:]
			rule, err := parseDomainEntry(b)
			if err != nil {
				return list, err
			}
			if rule.Value != "" {
				list.Rules = append(list.Rules, rule)
			}
		default:
			n := protowire.ConsumeFieldValue(num, typ, data)
			if n < 0 {
				return list, fmt.Errorf("bad field in GeoSite")
			}
			data = data[n:]
		}
	}
	return list, nil
}

func parseDomainEntry(data []byte) (core.Rule, error) {
	var domType uint64
	var value string
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return core.Rule{}, fmt.Errorf("bad tag in Domain")
		}
		data = data[n:]
		switch {
		case num == 1 && typ == protowire.VarintType: // type
			v, n := protowire.ConsumeVarint(data)
			if n < 0 {
				return core.Rule{}, fmt.Errorf("bad domain type")
			}
			data = data[n:]
			domType = v
		case num == 2 && typ == protowire.BytesType: // value
			b, n := protowire.ConsumeBytes(data)
			if n < 0 {
				return core.Rule{}, fmt.Errorf("bad domain value")
			}
			data = data[n:]
			value = string(b)
		default:
			n := protowire.ConsumeFieldValue(num, typ, data)
			if n < 0 {
				return core.Rule{}, fmt.Errorf("bad field in Domain")
			}
			data = data[n:]
		}
	}
	switch domType {
	case domainTypePlain:
		return core.NewRule(core.KindDomainKeyword, value), nil
	case domainTypeRegex:
		return core.NewRule(core.KindDomainRegex, value), nil
	case domainTypeDomain:
		return core.NewRule(core.KindDomainSuffix, value), nil
	default: // domainTypeFull
		return core.NewRule(core.KindDomain, value), nil
	}
}
