package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ddranic/ios-proxy-rules/internal/core"
)

func TestGenerateWritesPlatformOutputs(t *testing.T) {
	root := t.TempDir()
	outputDir := filepath.Join(root, "rules")

	g, err := New(outputDir)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	lists := []core.RuleList{
		{
			Name: "sample",
			Rules: []core.Rule{
				{Kind: core.KindDomain, Value: "example.com"},
				{Kind: core.KindDomainSuffix, Value: "example.org"},
				{Kind: core.KindDomainKeyword, Value: "cdn"},
				{Kind: core.KindDomainRegex, Value: "^ads\\..*\\.example\\.org$"},
			},
		},
	}

	if err := g.Generate(lists); err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	assertFileEquals(t, filepath.Join(outputDir, "shadowrocket", "sample.list"), strings.Join([]string{
		"DOMAIN,example.com",
		"DOMAIN-SUFFIX,example.org",
		"DOMAIN-KEYWORD,cdn",
		"URL-REGEX,^https?://ads\\..*\\.example\\.org(?::\\d+)?(?:/|$)",
	}, "\n"))

	assertFileEquals(t, filepath.Join(outputDir, "clash", "sample.yaml"), strings.Join([]string{
		"payload:",
		"  - DOMAIN,example.com",
		"  - DOMAIN-SUFFIX,example.org",
		"  - DOMAIN-KEYWORD,cdn",
		"  - DOMAIN-REGEX,^ads\\..*\\.example\\.org$",
	}, "\n"))
}

func TestGenerateRejectsUnsafeOutputDir(t *testing.T) {
	g, err := New(".")
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	err = g.Generate([]core.RuleList{
		{
			Name: "sample",
			Rules: []core.Rule{
				{Kind: core.KindDomainSuffix, Value: "example.org"},
			},
		},
	})
	if err == nil {
		t.Fatalf("expected unsafe output dir error")
	}
	if !strings.Contains(err.Error(), "refusing to use current working directory") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGenerateRejectsUnknownPlatform(t *testing.T) {
	outputDir := filepath.Join(t.TempDir(), "rules")
	_, err := New(outputDir, Platform("unknown"))
	if err == nil {
		t.Fatalf("expected unknown platform error")
	}
	if !strings.Contains(err.Error(), "unsupported platform") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewKeepsProvidedPlatforms(t *testing.T) {
	outputDir := filepath.Join(t.TempDir(), "rules")
	g, err := New(outputDir, PlatformClash, PlatformClash, PlatformSurge)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if len(g.platforms) != 3 {
		t.Fatalf("expected 3 platforms, got %d", len(g.platforms))
	}
	if g.platforms[0] != PlatformClash || g.platforms[1] != PlatformClash || g.platforms[2] != PlatformSurge {
		t.Fatalf("unexpected platform order: %#v", g.platforms)
	}
}

func assertFileEquals(t *testing.T, path string, expected string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(data) != expected {
		t.Fatalf("file mismatch for %s:\n--- got ---\n%s\n--- want ---\n%s", path, data, expected)
	}
}
