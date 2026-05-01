package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ddranic/ios-proxy-rules/internal/core"
)

type Generator struct {
	outputDir string
	platforms []Platform
}

func New(outputDir string, platforms ...Platform) (*Generator, error) {
	if len(platforms) == 0 {
		platforms = AllPlatforms
	}

	selectedPlatforms := make([]Platform, len(platforms))
	copy(selectedPlatforms, platforms)

	for _, platform := range selectedPlatforms {
		if _, _, ok := configForPlatform(platform); !ok {
			return nil, fmt.Errorf("unsupported platform %q", platform)
		}
	}

	return &Generator{
		outputDir: outputDir,
		platforms: selectedPlatforms,
	}, nil
}

func (g *Generator) Generate(lists []core.RuleList) error {
	outputDir, err := filepath.Abs(g.outputDir)
	if err != nil {
		return fmt.Errorf("resolve output directory %q: %w", g.outputDir, err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolve current working directory: %w", err)
	}
	if err := validateOutputDir(outputDir, cwd); err != nil {
		return err
	}

	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("remove output directory %q: %w", outputDir, err)
	}
	return g.generateToDir(outputDir, lists)
}

func (g *Generator) generateToDir(outputDir string, lists []core.RuleList) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory %q: %w", outputDir, err)
	}

	for _, platform := range g.platforms {
		targetDir := filepath.Join(outputDir, string(platform))
		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			return fmt.Errorf("create target dir %q: %w", targetDir, err)
		}
	}

	for _, list := range lists {
		for _, platform := range g.platforms {
			ext, render, _ := configForPlatform(platform)
			lines, err := render(list.Name, list.Rules)
			if err != nil {
				return fmt.Errorf("render %s/%s: %w", platform, list.Name, err)
			}
			if len(lines) == 0 {
				continue
			}

			outputPath := filepath.Join(outputDir, string(platform), list.Name+ext)
			data := linesToBytes(lines)
			if err := os.WriteFile(outputPath, data, 0o644); err != nil {
				return fmt.Errorf("write %q: %w", outputPath, err)
			}
		}
	}

	return nil
}

func linesToBytes(lines []string) []byte {
	return []byte(strings.Join(lines, "\n"))
}

func validateOutputDir(outputDir, cwd string) error {
	clean := filepath.Clean(outputDir)
	if clean == string(filepath.Separator) {
		return fmt.Errorf("refusing to use filesystem root as output directory %q", outputDir)
	}

	volume := filepath.VolumeName(clean)
	if volume != "" && clean == volume+string(filepath.Separator) {
		return fmt.Errorf("refusing to use filesystem root as output directory %q", outputDir)
	}
	if filepath.Clean(cwd) == clean {
		return fmt.Errorf("refusing to use current working directory as output directory %q", outputDir)
	}

	return nil
}
