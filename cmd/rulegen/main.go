package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ddranic/ios-proxy-rules/internal/app"
)

func main() {
	geosite := flag.String("geosite", "", "Path to v2fly domain-list directory")
	geoip := flag.String("geoip", "", "Path to geoip.dat")
	output := flag.String("output", ".", "Path to output directory")
	flag.Parse()

	generated, err := app.Run(app.Config{
		GeoSite: *geosite,
		GeoIP:   *geoip,
		Output:  *output,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %d lists in %q\n", generated, *output)
}
