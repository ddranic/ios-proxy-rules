package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ddranic/ios-proxy-rules/internal/core"
)

const (
	DirMode  = 0755
	FileMode = 0644
)

type Generator struct {
	output string
}

func NewGenerator(output string) *Generator {
	return &Generator{
		output: output,
	}
}

func (generator *Generator) Generate(lists []core.RuleList) error {
	if err := os.RemoveAll(generator.output); err != nil {
		return fmt.Errorf("clean output directory %q: %w", generator.output, err)
	}

	if err := os.MkdirAll(generator.output, DirMode); err != nil {
		return fmt.Errorf("create output directory %q: %w", generator.output, err)
	}

	for _, platform := range Platforms {
		target := filepath.Join(generator.output, string(platform))
		if err := os.MkdirAll(target, DirMode); err != nil {
			return fmt.Errorf("create target dir %q: %w", target, err)
		}
	}

	for _, list := range lists {
		err := generator.generateList(list)
		if err != nil {
			return err
		}
	}

	return nil
}

func (generator *Generator) generateList(list core.RuleList) error {
	for _, platform := range Platforms {
		extension := PlatformsExtensions[platform]

		lines, err := render(platform, list.Name, list.Rules)
		if err != nil {
			return fmt.Errorf("render %s/%s: %w", platform, list.Name, err)
		}

		path := filepath.Join(generator.output, string(platform), list.Name+extension)

		data := linesToBytes(lines)
		if err := os.WriteFile(path, data, FileMode); err != nil {
			return fmt.Errorf("write %q: %w", path, err)
		}
	}

	return nil
}

func linesToBytes(lines []string) []byte {
	return []byte(strings.Join(lines, "\n"))
}
