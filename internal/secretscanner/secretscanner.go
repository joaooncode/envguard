package secretscanner

import (
	"bufio"
	"bytes"
	"math"
	"path/filepath"
	"regexp"
)

// Method identifies which detection technique produced a Match.
type Method string

const (
	// MethodPattern indicates a match found via a known-provider signature.
	MethodPattern Method = "pattern"
	// MethodEntropy indicates a match found via the entropy heuristic.
	MethodEntropy Method = "entropy"
)

// Match represents a single detected secret instance within scanned content.
type Match struct {
	Line     int
	Key      string
	Method   Method
	Provider string
}

type providerPattern struct {
	provider string
	regex    *regexp.Regexp
}

var providerPatterns = []providerPattern{
	{provider: "aws", regex: regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
	{provider: "stripe", regex: regexp.MustCompile(`sk_(live|test)_[0-9a-zA-Z]{16,}`)},
	{provider: "github", regex: regexp.MustCompile(`gh[pousr]_[0-9A-Za-z]{36}`)},
	{provider: "pem", regex: regexp.MustCompile(`-----BEGIN ((RSA|EC|OPENSSH|DSA) )?PRIVATE KEY-----`)},
	{provider: "bearer-token", regex: regexp.MustCompile(`Bearer [A-Za-z0-9\-_.]{20,}`)},
}

// keyValueLine extracts the key name from a "KEY=VALUE" style line, if present.
var keyValueLine = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=`)

// Options configures which detection techniques a Scanner runs.
type Options struct {
	// EntropyScan enables the entropy-based heuristic (opt-in; higher false-positive rate).
	EntropyScan bool
	// Providers restricts pattern matching to the named providers. Empty/nil means all shipped providers.
	Providers []string
	// Ignore suppresses matches whose Key equals or glob-matches one of these entries, regardless of provider.
	Ignore []string
}

// ValidProviders returns the names of all shipped secret providers.
func ValidProviders() []string {
	names := make([]string, len(providerPatterns))
	for i, pp := range providerPatterns {
		names[i] = pp.provider
	}
	return names
}

// IsValidProvider reports whether provider is one of the shipped provider names.
func IsValidProvider(provider string) bool {
	for _, pp := range providerPatterns {
		if pp.provider == provider {
			return true
		}
	}
	return false
}

func (o Options) providerEnabled(provider string) bool {
	if len(o.Providers) == 0 {
		return true
	}
	for _, p := range o.Providers {
		if p == provider {
			return true
		}
	}
	return false
}

func (o Options) keyIgnored(key string) bool {
	if key == "" {
		return false
	}
	for _, pattern := range o.Ignore {
		if pattern == key {
			return true
		}
		if matched, err := filepath.Match(pattern, key); err == nil && matched {
			return true
		}
	}
	return false
}

const (
	entropyMinValueLength = 16
	entropyThreshold      = 3.5
)

// Scanner inspects file content for embedded secrets using known-provider
// signature patterns and, optionally, an entropy-based heuristic.
type Scanner struct {
	opts Options
}

// New creates a Scanner with default (pattern-only) detection enabled.
func New() *Scanner {
	return &Scanner{}
}

// NewWithOptions creates a Scanner with the given detection Options.
func NewWithOptions(opts Options) *Scanner {
	return &Scanner{opts: opts}
}

// Scan inspects content line by line and returns every detected Match.
//
// The scanner's buffer is sized to cover content in full: bufio.Scanner's
// default 64KB per-line limit would otherwise make Scan() stop silently
// (bufio.ErrTooLong) on any single line larger than that — e.g. a PEM
// certificate or JSON credentials blob assigned to one env var — dropping
// every match on subsequent lines without reporting an error. Since no line
// can be longer than content itself, sizing the buffer to len(content)
// guarantees no legitimate line ever hits the limit. Any other scanning
// error is returned rather than swallowed.
func (s *Scanner) Scan(content []byte) ([]Match, error) {
	matches := make([]Match, 0)

	scanner := bufio.NewScanner(bytes.NewReader(content))
	maxTokenSize := len(content)
	if maxTokenSize < bufio.MaxScanTokenSize {
		maxTokenSize = bufio.MaxScanTokenSize
	}
	scanner.Buffer(make([]byte, 0, bufio.MaxScanTokenSize), maxTokenSize)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		matchedPattern := false
		for _, pp := range providerPatterns {
			if !s.opts.providerEnabled(pp.provider) {
				continue
			}
			if pp.regex.MatchString(line) {
				key := extractKey(line)
				matchedPattern = true
				if s.opts.keyIgnored(key) {
					continue
				}
				matches = append(matches, Match{
					Line:     lineNum,
					Key:      key,
					Method:   MethodPattern,
					Provider: pp.provider,
				})
			}
		}

		if !matchedPattern && s.opts.EntropyScan {
			if key, value, ok := extractKeyValue(line); ok && looksLikeSecret(value) && !s.opts.keyIgnored(key) {
				matches = append(matches, Match{
					Line:   lineNum,
					Key:    key,
					Method: MethodEntropy,
				})
			}
		}
	}

	return matches, scanner.Err()
}

// extractKeyValue splits a "KEY=VALUE" line into its key and value, if present.
func extractKeyValue(line string) (key string, value string, ok bool) {
	idx := bytes.IndexByte([]byte(line), '=')
	if idx < 0 {
		return "", "", false
	}
	key = extractKey(line)
	if key == "" {
		return "", "", false
	}
	value = line[idx+1:]
	return key, value, true
}

// looksLikeSecret reports whether value is long enough and random-looking enough
// (by Shannon entropy) to be flagged as a likely secret.
func looksLikeSecret(value string) bool {
	if len(value) < entropyMinValueLength {
		return false
	}
	return shannonEntropy(value) >= entropyThreshold
}

// shannonEntropy computes the Shannon entropy (bits per character) of s.
func shannonEntropy(s string) float64 {
	if s == "" {
		return 0
	}

	counts := make(map[rune]int)
	for _, r := range s {
		counts[r]++
	}

	length := float64(len(s))
	var entropy float64
	for _, c := range counts {
		p := float64(c) / length
		entropy -= p * math.Log2(p)
	}
	return entropy
}

// extractKey returns the key name from a "KEY=VALUE" line, or "" if not present.
func extractKey(line string) string {
	m := keyValueLine.FindStringSubmatch(line)
	if m == nil {
		return ""
	}
	return m[1]
}
