package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ddranic/ios-proxy-rules/internal/app"
)

func main() {
	input := flag.String("input", "data", "Path to source domain-list directory")
	output := flag.String("output", "rules", "Path to generated rules directory")
	flag.Parse()

	generated, err := app.Run(app.Config{
		Input:  *input,
		Output: *output,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %d lists in %q\n", generated, *output)
}
