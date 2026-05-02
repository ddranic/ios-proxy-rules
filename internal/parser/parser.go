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

type Parser struct {
	input string
	state state
}

type state struct {
	visited map[string]struct{}
	seen    map[core.Rule]struct{}
}

func NewParser(input string) *Parser {
	return &Parser{input: input}
}

func (parser *Parser) Parse() ([]core.RuleList, error) {
	files, err := os.ReadDir(parser.input)
	if err != nil {
		return nil, fmt.Errorf("read input directory %q: %w", parser.input, err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	lists := make([]core.RuleList, 0, len(files))

	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}

		parser.state = state{
			visited: make(map[string]struct{}),
			seen:    make(map[core.Rule]struct{}),
		}

		rules, err := parser.parseList(file.Name())
		if err != nil {
			return nil, fmt.Errorf("parse list %q: %w", file.Name(), err)
		}

		lists = append(lists, core.RuleList{
			Name:  file.Name(),
			Rules: rules,
		})
	}

	return lists, nil
}

func (parser *Parser) parseList(name string) ([]core.Rule, error) {
	if _, ok := parser.state.visited[name]; ok {
		return nil, nil
	}
	parser.state.visited[name] = struct{}{}

	path := filepath.Join(parser.input, name)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var out []core.Rule

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		token := parseToken(line)
		if token == "" {
			continue
		}

		if include, ok := strings.CutPrefix(token, PrefixInclude); ok {
			nested, err := parser.parseList(include)
			if err != nil {
				return nil, err
			}
			out = append(out, nested...)
			continue
		}

		rule, err := parseRule(token)
		if err != nil {
			return nil, err
		}

		if _, exists := parser.state.seen[rule]; !exists {
			parser.state.seen[rule] = struct{}{}
			out = append(out, rule)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func parseRule(token string) (core.Rule, error) {
	switch {
	case strings.HasPrefix(token, PrefixFull):
		return core.NewRule(core.KindDomain, strings.TrimPrefix(token, PrefixFull)), nil
	case strings.HasPrefix(token, PrefixDomain):
		return core.NewRule(core.KindDomainSuffix, strings.TrimPrefix(token, PrefixDomain)), nil
	case strings.HasPrefix(token, PrefixKeyword):
		return core.NewRule(core.KindDomainKeyword, strings.TrimPrefix(token, PrefixKeyword)), nil
	case strings.HasPrefix(token, PrefixRegexp):
		return core.NewRule(core.KindDomainRegex, strings.TrimPrefix(token, PrefixRegexp)), nil
	default:
		if strings.ContainsRune(token, ':') {
			return core.Rule{}, fmt.Errorf("unknown rule prefix or invalid format: %q", token)
		}
		return core.NewRule(core.KindDomainSuffix, token), nil
	}
}

func parseToken(line string) string {
	if i := strings.Index(line, "#"); i >= 0 {
		line = line[:i]
	}
	line = strings.TrimSpace(line)

	fields := strings.Fields(line)
	if len(fields) == 0 {
		return ""
	}
	token := fields[0]

	if strings.HasPrefix(token, PrefixRegexp) {
		return token
	}

	if i := strings.IndexRune(token, '@'); i > 0 {
		token = token[:i]
	}
	return token
}
