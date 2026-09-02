package detector

import (
	"path/filepath"
	"strings"
)

// DefaultAllowlist contains standard safe template and example file patterns.
var DefaultAllowlist = []string{
	".env.example",
	".env.sample",
	".env.template",
	".env.*.example",
	".env.*.sample",
	".env.*.template",
	"*.env.example",
	"*.env.sample",
	"*.env.template",
}

// Detector evaluates file paths against environment patterns and allowlist rules.
type Detector struct {
	customPatterns []string
	allowlist      []string
}

// New creates a Detector initialized with the default allowlist.
func New() *Detector {
	return &Detector{
		customPatterns: make([]string, 0),
		allowlist:      DefaultAllowlist,
	}
}

// NewWithAllowlist creates a Detector with a custom allowlist appended to defaults.
func NewWithAllowlist(allowlist []string) *Detector {
	return NewWithPatterns(nil, allowlist)
}

// NewWithPatterns creates a Detector with custom environment patterns and allowlist appended to defaults.
func NewWithPatterns(customPatterns []string, customAllowlist []string) *Detector {
	combinedAllowlist := make([]string, 0, len(DefaultAllowlist)+len(customAllowlist))
	combinedAllowlist = append(combinedAllowlist, DefaultAllowlist...)
	combinedAllowlist = append(combinedAllowlist, customAllowlist...)

	return &Detector{
		customPatterns: customPatterns,
		allowlist:      combinedAllowlist,
	}
}


// extractBaseName extracts the filename from a given path across Unix and Windows formats.
func extractBaseName(path string) string {
	cleanPath := strings.ReplaceAll(path, "\\", "/")
	cleanPath = strings.TrimRight(cleanPath, "/")
	if cleanPath == "" || cleanPath == "." {
		return ""
	}
	idx := strings.LastIndex(cleanPath, "/")
	if idx >= 0 {
		return cleanPath[idx+1:]
	}
	return cleanPath
}

// IsEnvFile reports whether path corresponds to an environment file name (.env, .env.*, *.env, etc.).
// Operates exclusively on file path strings without inspecting disk or file contents.
func (d *Detector) IsEnvFile(path string) bool {
	base := extractBaseName(path)
	if base == "" {
		return false
	}

	lower := strings.ToLower(base)

	// Check custom patterns first if configured
	for _, pattern := range d.customPatterns {
		patternLower := strings.ToLower(pattern)
		if matchAllowlistPattern(patternLower, lower) {
			return true
		}
	}

	// Direct match: .env
	if lower == ".env" {
		return true
	}

	// Prefix match: .env.* (e.g., .env.local, .env.production, .env.example)
	if strings.HasPrefix(lower, ".env.") {
		return true
	}

	// Suffix match: *.env (e.g., backend.env, prod.env)
	if strings.HasSuffix(lower, ".env") {
		return true
	}

	// Segment match: *.env.* (e.g., app.env.local, backend.env.example)
	if strings.Contains(lower, ".env.") {
		return true
	}

	return false
}


// IsAllowed reports whether path is a recognized safe template/example file according to the allowlist.
func (d *Detector) IsAllowed(path string) bool {
	base := extractBaseName(path)
	if base == "" {
		return false
	}

	lower := strings.ToLower(base)

	// Check configured allowlist patterns
	for _, pattern := range d.allowlist {
		patternLower := strings.ToLower(pattern)
		if matchAllowlistPattern(patternLower, lower) {
			return true
		}
	}

	// General safe suffixes for env files (*.example, *.sample, *.template)
	if d.IsEnvFile(path) {
		if strings.HasSuffix(lower, ".example") ||
			strings.HasSuffix(lower, ".sample") ||
			strings.HasSuffix(lower, ".template") {
			return true
		}
	}

	return false
}

// matchAllowlistPattern matches a pattern against a filename.
func matchAllowlistPattern(pattern, name string) bool {
	if pattern == name {
		return true
	}
	matched, err := filepath.Match(pattern, name)
	if err == nil && matched {
		return true
	}
	return false
}

// Detect analyzes a path and returns whether it is an environment file and whether it is allowed.
func (d *Detector) Detect(path string) (isEnv bool, isAllowed bool) {
	isEnv = d.IsEnvFile(path)
	if !isEnv {
		return false, false
	}
	isAllowed = d.IsAllowed(path)
	return isEnv, isAllowed
}

// DefaultDetector is a package-level instance using standard configurations.
var DefaultDetector = New()

// IsEnvFile reports whether the given path represents an environment file using default patterns.
func IsEnvFile(path string) bool {
	return DefaultDetector.IsEnvFile(path)
}

// IsAllowed reports whether the given path represents a safe exception using default allowlist.
func IsAllowed(path string) bool {
	return DefaultDetector.IsAllowed(path)
}

// Detect analyzes a path using the default detector instance.
func Detect(path string) (isEnv bool, isAllowed bool) {
	return DefaultDetector.Detect(path)
}
