package app

import (
	"fmt"

	"github.com/ddranic/ios-proxy-rules/internal/generator"
	"github.com/ddranic/ios-proxy-rules/internal/parser"
)

type Config struct {
	Input  string
	Output string
}

func Run(cfg Config) (int, error) {
	parser := parser.NewParser(cfg.Input)
	lists, err := parser.Parse()
	if err != nil {
		return 0, fmt.Errorf("parse lists: %w", err)
	}

	generator := generator.NewGenerator(cfg.Output)
	err = generator.Generate(lists)
	if err != nil {
		return 0, fmt.Errorf("generate rules: %w", err)
	}

	return len(lists), nil
}
