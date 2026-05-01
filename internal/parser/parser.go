package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ddranic/ios-proxy-rules/internal/core"
)

func Parse(dataDir string) (Result, error) {
	absoluteDataDir, err := filepath.Abs(dataDir)
	if err != nil {
		return Result{}, fmt.Errorf("resolve data directory %q: %w", dataDir, err)
	}

	files, err := os.ReadDir(absoluteDataDir)
	if err != nil {
		return Result{}, fmt.Errorf("read data directory %q: %w", dataDir, err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	result := Result{
		Lists: make([]core.RuleList, 0, len(files)),
	}

	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}

		listRules, err := parseList(absoluteDataDir, file.Name())
		if err != nil {
			return Result{}, fmt.Errorf("parse list %q: %w", file.Name(), err)
		}
		if len(listRules) == 0 {
			continue
		}
		result.Lists = append(result.Lists, core.RuleList{
			Name:  file.Name(),
			Rules: listRules,
		})
	}

	return result, nil
}

func parseList(dataDir, listName string) ([]core.Rule, error) {
	st := state{
		visited: make(map[string]struct{}),
		seen:    make(map[string]struct{}),
	}

	var out []core.Rule
	if err := parseInto(dataDir, listName, &st, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func parseInto(dataDir, listName string, st *state, out *[]core.Rule) error {
	path, rel, err := resolveListPath(dataDir, listName)
	if err != nil {
		return err
	}

	if _, exists := st.visited[rel]; exists {
		return nil
	}
	st.visited[rel] = struct{}{}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %q: %w", rel, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		token := parseToken(scanner.Text())
		if token == "" {
			continue
		}

		if err := handleToken(dataDir, rel, lineNumber, token, st, out); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan %q: %w", rel, err)
	}
	return nil
}

func handleToken(dataDir, rel string, line int, token string, st *state, out *[]core.Rule) error {
	if strings.HasPrefix(token, "include:") {
		return handleInclude(dataDir, rel, line, token, st, out)
	}
	return handleRule(rel, line, token, st, out)
}

func handleInclude(dataDir, rel string, line int, token string, st *state, out *[]core.Rule) error {
	target := strings.TrimPrefix(token, "include:")
	if target == "" {
		return fmt.Errorf("%s:%d: empty include target (%s)", rel, line, token)
	}
	if err := parseInto(dataDir, target, st, out); err != nil {
		return fmt.Errorf("include %q in %s:%d: %w", target, rel, line, err)
	}
	return nil
}

func handleRule(rel string, line int, token string, st *state, out *[]core.Rule) error {
	rule, ok := parseRule(token)
	if !ok {
		return fmt.Errorf("%s:%d: unsupported token (%s)", rel, line, token)
	}
	appendUniqueRule(st, out, rule)
	return nil
}

func appendUniqueRule(st *state, out *[]core.Rule, rule core.Rule) {
	key := ruleKey(rule)
	if _, exists := st.seen[key]; exists {
		return
	}
	st.seen[key] = struct{}{}
	*out = append(*out, rule)
}

func parseRule(token string) (core.Rule, bool) {
	switch {
	case strings.HasPrefix(token, "full:"):
		return newRule(core.KindDomain, strings.TrimPrefix(token, "full:"))
	case strings.HasPrefix(token, "domain:"):
		return newRule(core.KindDomainSuffix, strings.TrimPrefix(token, "domain:"))
	case strings.HasPrefix(token, "keyword:"):
		return newRule(core.KindDomainKeyword, strings.TrimPrefix(token, "keyword:"))
	case strings.HasPrefix(token, "regexp:"):
		return newRule(core.KindDomainRegex, strings.TrimPrefix(token, "regexp:"))
	default:
		if strings.ContainsRune(token, ':') {
			return core.Rule{}, false
		}
		return newRule(core.KindDomainSuffix, token)
	}
}

func newRule(kind core.RuleKind, value string) (core.Rule, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return core.Rule{}, false
	}
	return core.Rule{Kind: kind, Value: value}, true
}

func stripComment(line string) string {
	if i := strings.Index(line, "#"); i >= 0 {
		line = line[:i]
	}
	return strings.TrimSpace(line)
}

func firstToken(line string) string {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func parseToken(line string) string {
	token := firstToken(stripComment(line))
	if token == "" {
		return ""
	}

	// Attributes are expected as separate tokens (e.g. "gov.cn @cn"),
	// but we also tolerate accidental inline form like "gov.cn@cn".
	if !strings.HasPrefix(token, "regexp:") {
		if i := strings.IndexRune(token, '@'); i > 0 {
			token = token[:i]
		}
	}
	return token
}

func ruleKey(rule core.Rule) string {
	return string(rule.Kind) + "\x00" + rule.Value
}

func resolveListPath(dataDir, listName string) (string, string, error) {
	name := strings.TrimSpace(listName)
	if name == "" {
		return "", "", fmt.Errorf("empty list name")
	}
	if name == "." || name == ".." {
		return "", "", fmt.Errorf("invalid list name: %q", listName)
	}
	if strings.ContainsAny(name, `/\`) {
		return "", "", fmt.Errorf("path separators are not allowed in list name: %q", listName)
	}

	fullPath := filepath.Join(dataDir, name)
	return fullPath, name, nil
}
