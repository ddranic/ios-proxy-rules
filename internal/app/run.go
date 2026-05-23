package app

import (
	"fmt"

	"github.com/ddranic/ios-proxy-rules/internal/core"
	"github.com/ddranic/ios-proxy-rules/internal/generator"
	"github.com/ddranic/ios-proxy-rules/internal/parser"
)

type Config struct {
	GeoSite string
	GeoIP   string
	Output  string
}

func Run(cfg Config) (int, error) {
	var lists []core.RuleList
	if cfg.GeoSite != "" {
		var err error
		lists, err = parser.ParseGeoSite(cfg.GeoSite)
		if err != nil {
			return 0, fmt.Errorf("parse geosite: %w", err)
		}
	}

	var entries []parser.GeoIPEntry
	if cfg.GeoIP != "" {
		var err error
		entries, err = parser.ParseGeoIP(cfg.GeoIP)
		if err != nil {
			return 0, fmt.Errorf("parse geoip: %w", err)
		}
	}

	gen := generator.NewGenerator(cfg.Output)
	if err := gen.Generate(lists, entries); err != nil {
		return 0, fmt.Errorf("generate rules: %w", err)
	}

	return len(lists), nil
}
