package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ddranic/ios-proxy-rules/internal/app"
)

func main() {
	inputDir := flag.String("input", "data", "Path to source domain-list directory")
	outputDir := flag.String("output", "rules", "Path to generated rules directory")
	flag.Parse()

	report, err := app.Run(app.Config{
		InputDir:  *inputDir,
		OutputDir: *outputDir,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %d lists in %q\n", report.GeneratedLists, *outputDir)
}
