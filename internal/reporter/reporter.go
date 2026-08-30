package reporter

import (
	"fmt"
	"io"

	"github.com/joaooncode/envguard/internal/scanner"
)

// Format represents the supported output report format.
type Format string

const (
	// FormatTerminal outputs human-readable styled terminal text.
	FormatTerminal Format = "terminal"
	// FormatJSON outputs structured JSON data.
	FormatJSON Format = "json"
)

// Reporter defines the interface for rendering scan results.
type Reporter interface {
	Render(result *scanner.Result, w io.Writer) error
}

// Options configures reporter behavior.
type Options struct {
	NoColor bool
	Pretty  bool
	Version string
}

// Option configures an Options instance.
type Option func(*Options)

// WithNoColor disables ANSI color escape sequences in terminal output.
func WithNoColor(noColor bool) Option {
	return func(o *Options) {
		o.NoColor = noColor
	}
}

// WithPretty enables pretty printing with indentation.
func WithPretty(pretty bool) Option {
	return func(o *Options) {
		o.Pretty = pretty
	}
}

// WithVersion specifies the application version for report headers/metadata.
func WithVersion(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

// New creates a new Reporter implementation based on the requested format.
func New(format Format, opts ...Option) (Reporter, error) {
	options := Options{
		Pretty: true,
	}
	for _, opt := range opts {
		opt(&options)
	}

	switch format {
	case FormatTerminal:
		return NewTerminalReporter(options), nil
	case FormatJSON:
		return NewJSONReporter(options), nil
	default:
		return nil, fmt.Errorf("unsupported reporter format: %q", format)
	}
}
