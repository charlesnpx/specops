package main

import (
	"fmt"
	"os"

	"github.com/specops/specops/internal/cli"
	"github.com/specops/specops/internal/output"
)

var version = "0.1.3-dev"

func main() {
	root := cli.NewRoot(os.Stdout, os.Stderr, version)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(output.ExitCode(err))
	}
}
