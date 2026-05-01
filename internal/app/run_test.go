package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunGeneratesRules(t *testing.T) {
	root := t.TempDir()
	inputDir := filepath.Join(root, "data")
	outputDir := filepath.Join(root, "rules")

	if err := os.MkdirAll(inputDir, 0o755); err != nil {
		t.Fatalf("mkdir input dir: %v", err)
	}
	content := strings.Join([]string{
		"full:example.com",
		"domain:example.org",
		"",
	}, "\n")
	if err := os.WriteFile(filepath.Join(inputDir, "mainlist"), []byte(content), 0o644); err != nil {
		t.Fatalf("write list file: %v", err)
	}

	report, err := Run(Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if report.GeneratedLists != 1 {
		t.Fatalf("expected 1 generated list, got %d", report.GeneratedLists)
	}

	shadowrocketPath := filepath.Join(outputDir, "shadowrocket", "mainlist.list")
	if _, err := os.Stat(shadowrocketPath); err != nil {
		t.Fatalf("expected generated file %q: %v", shadowrocketPath, err)
	}
}

func TestRunRejectsUnsupportedParserToken(t *testing.T) {
	root := t.TempDir()
	inputDir := filepath.Join(root, "data")
	outputDir := filepath.Join(root, "rules")

	if err := os.MkdirAll(inputDir, 0o755); err != nil {
		t.Fatalf("mkdir input dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "mainlist"), []byte("unknown:value\n"), 0o644); err != nil {
		t.Fatalf("write list file: %v", err)
	}

	_, err := Run(Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
	})
	if err == nil {
		t.Fatalf("expected run error")
	}
	if !strings.Contains(err.Error(), "unsupported token") {
		t.Fatalf("unexpected error: %v", err)
	}
}
