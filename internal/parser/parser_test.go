package parser

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ddranic/ios-proxy-rules/internal/core"
)

func TestParseWithIncludes(t *testing.T) {
	dataDir := t.TempDir()
	writeTestFile(t, filepath.Join(dataDir, "base"), strings.Join([]string{
		"full:example.com",
		"domain:example.org",
		"",
	}, "\n"))
	writeTestFile(t, filepath.Join(dataDir, "mainlist"), strings.Join([]string{
		"include:base",
		"keyword:cdn",
		"",
	}, "\n"))

	result, err := Parse(dataDir)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.Lists) != 2 {
		t.Fatalf("expected 2 lists, got %d", len(result.Lists))
	}

	mainRules := findListRules(result.Lists, "mainlist")
	if len(mainRules) != 3 {
		t.Fatalf("expected 3 rules in mainlist, got %d", len(mainRules))
	}
}

func TestParseRejectsUnsupportedToken(t *testing.T) {
	dataDir := t.TempDir()
	writeTestFile(t, filepath.Join(dataDir, "mainlist"), "unknown:value\n")

	_, err := Parse(dataDir)
	if err == nil {
		t.Fatalf("expected unsupported token error")
	}
	if !strings.Contains(err.Error(), "unsupported token") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseRejectsEmptyIncludeTarget(t *testing.T) {
	dataDir := t.TempDir()
	writeTestFile(t, filepath.Join(dataDir, "mainlist"), "include:\n")

	_, err := Parse(dataDir)
	if err == nil {
		t.Fatalf("expected empty include target error")
	}
	if !strings.Contains(err.Error(), "empty include target") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseReturnsListsInDeterministicOrder(t *testing.T) {
	dataDir := t.TempDir()
	writeTestFile(t, filepath.Join(dataDir, "z-last"), "full:last.example\n")
	writeTestFile(t, filepath.Join(dataDir, "a-first"), "full:first.example\n")

	result, err := Parse(dataDir)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(result.Lists) != 2 {
		t.Fatalf("expected 2 lists, got %d", len(result.Lists))
	}
	if result.Lists[0].Name != "a-first" || result.Lists[1].Name != "z-last" {
		t.Fatalf("unexpected list order: %q, %q", result.Lists[0].Name, result.Lists[1].Name)
	}
}

func TestParseRejectsPathTraversalInInclude(t *testing.T) {
	dataDir := t.TempDir()
	writeTestFile(t, filepath.Join(dataDir, "mainlist"), "include:../secrets\n")

	_, err := Parse(dataDir)
	if err == nil {
		t.Fatalf("expected path traversal error")
	}
	if !strings.Contains(err.Error(), "path separators are not allowed") {
		t.Fatalf("unexpected traversal error: %v", err)
	}
}

func TestParseRule(t *testing.T) {
	tests := []struct {
		name   string
		token  string
		want   core.Rule
		wantOK bool
	}{
		{
			name:   "full",
			token:  "full:example.com",
			want:   core.Rule{Kind: core.KindDomain, Value: "example.com"},
			wantOK: true,
		},
		{
			name:   "domain",
			token:  "domain:example.com",
			want:   core.Rule{Kind: core.KindDomainSuffix, Value: "example.com"},
			wantOK: true,
		},
		{
			name:   "keyword",
			token:  "keyword:cdn",
			want:   core.Rule{Kind: core.KindDomainKeyword, Value: "cdn"},
			wantOK: true,
		},
		{
			name:   "regexp",
			token:  "regexp:^ads\\..*$",
			want:   core.Rule{Kind: core.KindDomainRegex, Value: "^ads\\..*$"},
			wantOK: true,
		},
		{
			name:   "default suffix",
			token:  "example.org",
			want:   core.Rule{Kind: core.KindDomainSuffix, Value: "example.org"},
			wantOK: true,
		},
		{
			name:   "unknown prefixed token",
			token:  "unknown:value",
			want:   core.Rule{},
			wantOK: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := parseRule(test.token)
			if ok != test.wantOK {
				t.Fatalf("ok mismatch: got %v want %v", ok, test.wantOK)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("rule mismatch: got %+v want %+v", got, test.want)
			}
		})
	}
}

func TestParseIgnoresAttributesAndComments(t *testing.T) {
	dataDir := t.TempDir()
	writeTestFile(t, filepath.Join(dataDir, "mainlist"), strings.Join([]string{
		"gov.cn @cn",
		"example.org@cn",
		"domain:foo.bar @!cn @ads",
		"full:api.test # comment",
		"regexp:^foo@bar$ @cn",
		"",
	}, "\n"))

	result, err := Parse(dataDir)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	rules := findListRules(result.Lists, "mainlist")
	if len(rules) != 5 {
		t.Fatalf("expected 5 rules, got %d", len(rules))
	}

	want := []core.Rule{
		{Kind: core.KindDomainSuffix, Value: "gov.cn"},
		{Kind: core.KindDomainSuffix, Value: "example.org"},
		{Kind: core.KindDomainSuffix, Value: "foo.bar"},
		{Kind: core.KindDomain, Value: "api.test"},
		{Kind: core.KindDomainRegex, Value: "^foo@bar$"},
	}
	if !reflect.DeepEqual(rules, want) {
		t.Fatalf("rules mismatch:\n got: %#v\nwant: %#v", rules, want)
	}
}

func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func findListRules(lists []core.RuleList, name string) []core.Rule {
	for _, list := range lists {
		if list.Name == name {
			return list.Rules
		}
	}
	return nil
}
