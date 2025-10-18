package main

import (
	"flag"
	"log"
	"os"

	"api.etin.dev/internal/version"
	"api.etin.dev/pkg/openapi"
)

func main() {
	var output string
	flag.StringVar(&output, "output", "cmd/api/openapi.json", "path to write the OpenAPI specification")
	flag.Parse()

	spec, err := openapi.Build(version.Number)
	if err != nil {
		log.Fatalf("generate openapi spec: %v", err)
	}

	if err := os.WriteFile(output, spec, 0o644); err != nil {
		log.Fatalf("write openapi spec: %v", err)
	}
}
