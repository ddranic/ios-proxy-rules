package parser

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strings"

	"google.golang.org/protobuf/encoding/protowire"
)

// GeoIPEntry holds all CIDRs for one country code.
type GeoIPEntry struct {
	CountryCode string
	CIDRs       []string
}

// ParseGeoIP reads a v2fly geoip.dat and returns entries sorted by country code.
func ParseGeoIP(path string) ([]GeoIPEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", path, err)
	}
	entries, err := parseGeoIPList(data)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", path, err)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CountryCode < entries[j].CountryCode
	})
	return entries, nil
}

func parseGeoIPList(data []byte) ([]GeoIPEntry, error) {
	var entries []GeoIPEntry
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return nil, fmt.Errorf("bad tag")
		}
		data = data[n:]
		if num == 1 && typ == protowire.BytesType {
			b, n := protowire.ConsumeBytes(data)
			if n < 0 {
				return nil, fmt.Errorf("bad GeoIP message")
			}
			data = data[n:]
			e, err := parseGeoIPEntry(b)
			if err != nil {
				return nil, err
			}
			if len(e.CIDRs) > 0 {
				entries = append(entries, e)
			}
		} else {
			n := protowire.ConsumeFieldValue(num, typ, data)
			if n < 0 {
				return nil, fmt.Errorf("bad field")
			}
			data = data[n:]
		}
	}
	return entries, nil
}

func parseGeoIPEntry(data []byte) (GeoIPEntry, error) {
	var e GeoIPEntry
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return e, fmt.Errorf("bad tag in GeoIP")
		}
		data = data[n:]
		switch {
		case num == 1 && typ == protowire.BytesType: // country_code
			b, n := protowire.ConsumeBytes(data)
			if n < 0 {
				return e, fmt.Errorf("bad country_code")
			}
			data = data[n:]
			e.CountryCode = strings.ToLower(string(b))
		case num == 2 && typ == protowire.BytesType: // cidr
			b, n := protowire.ConsumeBytes(data)
			if n < 0 {
				return e, fmt.Errorf("bad CIDR message")
			}
			data = data[n:]
			cidr, err := parseCIDREntry(b)
			if err != nil {
				return e, err
			}
			e.CIDRs = append(e.CIDRs, cidr)
		default:
			n := protowire.ConsumeFieldValue(num, typ, data)
			if n < 0 {
				return e, fmt.Errorf("bad field in GeoIP")
			}
			data = data[n:]
		}
	}
	return e, nil
}

func parseCIDREntry(data []byte) (string, error) {
	var ip []byte
	var prefix uint64
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return "", fmt.Errorf("bad tag in CIDR")
		}
		data = data[n:]
		switch {
		case num == 1 && typ == protowire.BytesType: // ip bytes
			b, n := protowire.ConsumeBytes(data)
			if n < 0 {
				return "", fmt.Errorf("bad ip bytes")
			}
			data = data[n:]
			ip = make([]byte, len(b))
			copy(ip, b)
		case num == 2 && typ == protowire.VarintType: // prefix uint32
			v, n := protowire.ConsumeVarint(data)
			if n < 0 {
				return "", fmt.Errorf("bad prefix")
			}
			data = data[n:]
			prefix = v
		default:
			n := protowire.ConsumeFieldValue(num, typ, data)
			if n < 0 {
				return "", fmt.Errorf("bad field in CIDR")
			}
			data = data[n:]
		}
	}
	if len(ip) != 4 && len(ip) != 16 {
		return "", fmt.Errorf("unexpected IP length %d", len(ip))
	}
	return fmt.Sprintf("%s/%d", net.IP(ip), prefix), nil
}
