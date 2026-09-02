package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"

	"github.com/joaooncode/envguard/internal/initializer"
)

type initConfig struct {
	path         string
	force        bool
	template     bool
	templateFrom string
}

func runInitCommand(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var cfg initConfig
	fs.StringVar(&cfg.path, "path", ".", "Target directory path to initialize")
	fs.StringVar(&cfg.path, "p", ".", "Target directory path to initialize (shorthand)")
	fs.BoolVar(&cfg.force, "force", false, "Overwrite existing configuration or template files")
	fs.BoolVar(&cfg.force, "f", false, "Overwrite existing files (shorthand)")
	fs.BoolVar(&cfg.template, "template", false, "Generate a safe .env.example template file")
	fs.BoolVar(&cfg.template, "t", false, "Generate a safe .env.example template file (shorthand)")
	fs.StringVar(&cfg.templateFrom, "template-from", "", "Source .env file to sanitize and create .env.example from")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return ExitCodeSuccess
		}
		return ExitCodeUsageError
	}

	if cfg.path == "" {
		cfg.path = "."
	}

	// 1. Generate configuration file (.envguard.yaml)
	if err := initializer.GenerateConfig(cfg.path, cfg.force); err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return ExitCodeInternalError
	}
	configFilePath := filepath.Join(cfg.path, ".envguard.yaml")
	fmt.Fprintf(stdout, "Created configuration file: %s\n", configFilePath)

	// 2. Generate template if requested
	if cfg.template || cfg.templateFrom != "" {
		if err := initializer.GenerateTemplate(cfg.path, cfg.templateFrom, cfg.force); err != nil {
			fmt.Fprintf(stderr, "Error: %v\n", err)
			return ExitCodeInternalError
		}
		templateFilePath := filepath.Join(cfg.path, ".env.example")
		fmt.Fprintf(stdout, "Created template file: %s\n", templateFilePath)
	}

	return ExitCodeSuccess
}
