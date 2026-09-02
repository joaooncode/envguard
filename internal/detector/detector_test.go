package detector

import (
	"testing"
)

func TestIsEnvFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		// Standard .env files
		{"exact .env", ".env", true},
		{"exact .ENV (uppercase)", ".ENV", true},
		{"exact .Env (mixed case)", ".Env", true},

		// Prefix .env.*
		{".env.local", ".env.local", true},
		{".env.production", ".env.production", true},
		{".env.staging", ".env.staging", true},
		{".env.development", ".env.development", true},
		{".env.test", ".env.test", true},
		{".env.dev.local", ".env.dev.local", true},
		{".ENV.LOCAL uppercase", ".ENV.LOCAL", true},

		// Suffix *.env
		{"backend.env", "backend.env", true},
		{"prod.env", "prod.env", true},
		{"docker.env", "docker.env", true},
		{"app.service.env", "app.service.env", true},
		{"PROD.ENV uppercase", "PROD.ENV", true},

		// Segment *.env.*
		{"app.env.local", "app.env.local", true},
		{"docker.env.production", "docker.env.production", true},

		// Allowlist candidates (they are env files, but safe templates)
		{".env.example", ".env.example", true},
		{".env.sample", ".env.sample", true},
		{".env.template", ".env.template", true},
		{".env.local.example", ".env.local.example", true},
		{"backend.env.template", "backend.env.template", true},

		// Nested paths (Unix and Windows)
		{"nested unix path", "configs/sub/.env", true},
		{"nested windows path", `C:\Users\dev\repo\.env.production`, true},
		{"relative current dir", "./.env", true},
		{"nested unix with prefix", "src/backend/.env.staging", true},
		{"nested windows with suffix", `services\auth\api.env`, true},

		// Negative cases (non-env files)
		{"empty string", "", false},
		{"dot only", ".", false},
		{"main.go", "main.go", false},
		{"environment.js", "environment.js", false},
		{"envelope.py", "envelope.py", false},
		{"env.go", "env.go", false},
		{"dotenv.rs", "dotenv.rs", false},
		{"env_vars.json", "env_vars.json", false},
		{"package.json", "package.json", false},
		{"README.md", "README.md", false},
		{"config.yaml", "config.yaml", false},
		{"app.environment", "app.environment", false},
		{"environment_test.go", "environment_test.go", false},
		{"some-env-file.txt", "some-env-file.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsEnvFile(tt.path)
			if got != tt.expected {
				t.Errorf("IsEnvFile(%q) = %v; want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestIsAllowed(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		// Standard default allowlist
		{".env.example", ".env.example", true},
		{".env.sample", ".env.sample", true},
		{".env.template", ".env.template", true},
		{".ENV.EXAMPLE (uppercase)", ".ENV.EXAMPLE", true},
		{".env.SAMPLE (mixed case)", ".env.SAMPLE", true},

		// Variants with suffixes
		{".env.local.example", ".env.local.example", true},
		{".env.production.sample", ".env.production.sample", true},
		{"backend.env.template", "backend.env.template", true},
		{"api.env.example", "api.env.example", true},

		// Paths with directories
		{"nested unix example", "config/docker/.env.example", true},
		{"nested windows template", `C:\Projects\app\.env.local.template`, true},

		// Env files that are NOT allowed (real secret files)
		{".env", ".env", false},
		{".env.local", ".env.local", false},
		{".env.production", ".env.production", false},
		{"backend.env", "backend.env", false},
		{"docker.env", "docker.env", false},

		// Non-env files (should not be in allowlist)
		{"main.go", "main.go", false},
		{"readme.md", "readme.md", false},
		{"example.go", "example.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAllowed(tt.path)
			if got != tt.expected {
				t.Errorf("IsAllowed(%q) = %v; want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestDetect(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		expectedIsEnv   bool
		expectedAllowed bool
	}{
		{"sensitive .env", ".env", true, false},
		{"sensitive .env.production", ".env.production", true, false},
		{"sensitive backend.env", "services/backend.env", true, false},
		{"safe .env.example", ".env.example", true, true},
		{"safe .env.local.template", "config/.env.local.template", true, true},
		{"regular code file", "pkg/detector/detector.go", false, false},
		{"regular config", "configs/settings.json", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isEnv, isAllowed := Detect(tt.path)
			if isEnv != tt.expectedIsEnv || isAllowed != tt.expectedAllowed {
				t.Errorf("Detect(%q) = (%v, %v); want (%v, %v)",
					tt.path, isEnv, isAllowed, tt.expectedIsEnv, tt.expectedAllowed)
			}
		})
	}
}

func TestCustomPatternsAndAllowlist(t *testing.T) {
	d := NewWithPatterns(
		[]string{"*.env.vault", ".env.secret"},
		[]string{".env.dist", "*.env.custom-example"},
	)

	// Built-in defaults should still work
	if !d.IsEnvFile(".env") {
		t.Errorf("expected default .env to be recognized")
	}
	if !d.IsAllowed(".env.example") {
		t.Errorf("expected default .env.example to be allowed")
	}

	// Custom patterns
	if !d.IsEnvFile("backend.env.vault") {
		t.Errorf("expected custom pattern backend.env.vault to be recognized")
	}
	if !d.IsEnvFile(".env.secret") {
		t.Errorf("expected custom pattern .env.secret to be recognized")
	}

	// Custom allowlist
	if !d.IsAllowed(".env.dist") {
		t.Errorf("expected custom allowlist .env.dist to be allowed")
	}
	if !d.IsAllowed("app.env.custom-example") {
		t.Errorf("expected custom allowlist app.env.custom-example to be allowed")
	}
}
