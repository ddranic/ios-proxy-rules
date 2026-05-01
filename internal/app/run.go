package app

import (
	"fmt"

	"github.com/ddranic/ios-proxy-rules/internal/generator"
	"github.com/ddranic/ios-proxy-rules/internal/parser"
)

type Config struct {
	InputDir  string
	OutputDir string
}

type Report struct {
	GeneratedLists int
}

func Run(cfg Config) (Report, error) {
	result, err := parser.Parse(cfg.InputDir)
	if err != nil {
		return Report{}, fmt.Errorf("parse domain lists: %w", err)
	}
	if len(result.Lists) == 0 {
		return Report{}, fmt.Errorf("no source lists found in %q", cfg.InputDir)
	}

	g, err := generator.New(cfg.OutputDir)
	if err != nil {
		return Report{}, fmt.Errorf("create generator: %w", err)
	}
	if err := g.Generate(result.Lists); err != nil {
		return Report{}, fmt.Errorf("generate rules: %w", err)
	}

	return Report{
		GeneratedLists: len(result.Lists),
	}, nil
}
