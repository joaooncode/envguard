package main

import (
	"os"

	"github.com/joaooncode/envguard/internal/cli"
)

func main() {
	code := cli.Run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(code)
}

