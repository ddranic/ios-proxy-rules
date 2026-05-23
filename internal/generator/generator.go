package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ddranic/ios-proxy-rules/internal/core"
	"github.com/ddranic/ios-proxy-rules/internal/parser"
)

const (
	DirMode  = 0755
	FileMode = 0644
)

type Generator struct {
	output string
}

func NewGenerator(output string) *Generator {
	return &Generator{output: output}
}

func (g *Generator) Generate(lists []core.RuleList, entries []parser.GeoIPEntry) error {
	if err := os.MkdirAll(g.output, DirMode); err != nil {
		return fmt.Errorf("create output dir %q: %w", g.output, err)
	}

	tmp := filepath.Join(g.output, ".tmp")
	if err := os.RemoveAll(tmp); err != nil {
		return fmt.Errorf("clean temp dir: %w", err)
	}

	if err := g.generateTo(tmp, lists, entries); err != nil {
		os.RemoveAll(tmp)
		return err
	}

	// atomically swap each generated subdir
	for _, sub := range []string{"geosite", "geoip"} {
		src := filepath.Join(tmp, sub)
		if _, err := os.Stat(src); os.IsNotExist(err) {
			continue
		}
		dst := filepath.Join(g.output, sub)
		os.RemoveAll(dst)
		if err := os.Rename(src, dst); err != nil {
			os.RemoveAll(tmp)
			return fmt.Errorf("swap %s: %w", sub, err)
		}
	}

	return os.RemoveAll(tmp)
}

func (g *Generator) generateTo(dir string, lists []core.RuleList, entries []parser.GeoIPEntry) error {
	if len(lists) > 0 {
		for _, platform := range Platforms {
			if err := os.MkdirAll(filepath.Join(dir, "geosite", string(platform)), DirMode); err != nil {
				return fmt.Errorf("create geosite dir: %w", err)
			}
		}
		for _, list := range lists {
			if err := g.writeList(dir, list); err != nil {
				return err
			}
		}
	}

	if len(entries) > 0 {
		for _, platform := range Platforms {
			if err := os.MkdirAll(filepath.Join(dir, "geoip", string(platform)), DirMode); err != nil {
				return fmt.Errorf("create geoip dir: %w", err)
			}
		}
		for _, entry := range entries {
			if err := g.writeGeoIPEntry(dir, entry); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Generator) writeList(dir string, list core.RuleList) error {
	for _, platform := range Platforms {
		lines, err := render(platform, list.Name, list.Rules)
		if err != nil {
			return fmt.Errorf("render %s/%s: %w", platform, list.Name, err)
		}
		path := filepath.Join(dir, "geosite", string(platform), list.Name+PlatformsExtensions[platform])
		if err := os.WriteFile(path, linesToBytes(lines), FileMode); err != nil {
			return fmt.Errorf("write %q: %w", path, err)
		}
	}
	return nil
}

func (g *Generator) writeGeoIPEntry(dir string, entry parser.GeoIPEntry) error {
	for _, platform := range Platforms {
		lines, err := renderGeoIP(platform, entry)
		if err != nil {
			return fmt.Errorf("render geoip %s/%s: %w", platform, entry.CountryCode, err)
		}
		path := filepath.Join(dir, "geoip", string(platform), entry.CountryCode+PlatformsExtensions[platform])
		if err := os.WriteFile(path, linesToBytes(lines), FileMode); err != nil {
			return fmt.Errorf("write %q: %w", path, err)
		}
	}
	return nil
}

func linesToBytes(lines []string) []byte {
	return []byte(strings.Join(lines, "\n"))
}
