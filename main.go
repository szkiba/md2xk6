// Package main contains the main entry point for the md2xk6 CLI.
package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	filename := "README.md"

	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	contents, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	modules, err := extract(contents)
	if err != nil {
		log.Fatal(err)
	}

	for _, module := range modules {
		fmt.Fprintf(os.Stdout, " --with %s", module) //nolint:errcheck
	}
}
