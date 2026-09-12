package secretscanner

import (
	"bufio"
	"bytes"
	"testing"
)

func TestScanDetectsAWSAccessKey(t *testing.T) {
	content := []byte("AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n")

	s := New()
	matches, err := s.Scan(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d: %+v", len(matches), matches)
	}

	m := matches[0]
	if m.Line != 1 {
		t.Errorf("Line = %d, want 1", m.Line)
	}
	if m.Key != "AWS_ACCESS_KEY_ID" {
		t.Errorf("Key = %q, want %q", m.Key, "AWS_ACCESS_KEY_ID")
	}
	if m.Method != MethodPattern {
		t.Errorf("Method = %q, want %q", m.Method, MethodPattern)
	}
	if m.Provider != "aws" {
		t.Errorf("Provider = %q, want %q", m.Provider, "aws")
	}
}

func TestScanDetectsStripeSecretKey(t *testing.T) {
	content := []byte("STRIPE_SECRET_KEY=sk_live_4eC39HqLyjWDarjtT1zdp7dc\n")

	s := New()
	matches, err := s.Scan(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d: %+v", len(matches), matches)
	}

	m := matches[0]
	if m.Provider != "stripe" {
		t.Errorf("Provider = %q, want %q", m.Provider, "stripe")
	}
	if m.Key != "STRIPE_SECRET_KEY" {
		t.Errorf("Key = %q, want %q", m.Key, "STRIPE_SECRET_KEY")
	}
}

func TestScanDetectsGitHubToken(t *testing.T) {
	content := []byte("GITHUB_TOKEN=ghp_16C7e42F292c6912E7710c838347Ae178B4a\n")

	s := New()
	matches, err := s.Scan(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d: %+v", len(matches), matches)
	}

	m := matches[0]
	if m.Provider != "github" {
		t.Errorf("Provider = %q, want %q", m.Provider, "github")
	}
	if m.Key != "GITHUB_TOKEN" {
		t.Errorf("Key = %q, want %q", m.Key, "GITHUB_TOKEN")
	}
}

func TestScanDetectsPEMPrivateKeyHeader(t *testing.T) {
	content := []byte("PRIVATE_KEY=\"-----BEGIN RSA PRIVATE KEY-----\n")

	s := New()
	matches, err := s.Scan(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d: %+v", len(matches), matches)
	}

	m := matches[0]
	if m.Provider != "pem" {
		t.Errorf("Provider = %q, want %q", m.Provider, "pem")
	}
}

func TestScanDetectsBearerToken(t *testing.T) {
	content := []byte("AUTH_HEADER=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.abcdef\n")

	s := New()
	matches, err := s.Scan(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d: %+v", len(matches), matches)
	}

	m := matches[0]
	if m.Provider != "bearer-token" {
		t.Errorf("Provider = %q, want %q", m.Provider, "bearer-token")
	}
}

func TestScanTracksLineNumbersAndIgnoresPlainLines(t *testing.T) {
	content := []byte("NODE_ENV=production\nAWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\nPORT=3000\nSTRIPE_SECRET_KEY=sk_live_4eC39HqLyjWDarjtT1zdp7dc\n")

	s := New()
	matches, err := s.Scan(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d: %+v", len(matches), matches)
	}
	if matches[0].Line != 2 {
		t.Errorf("matches[0].Line = %d, want 2", matches[0].Line)
	}
	if matches[1].Line != 4 {
		t.Errorf("matches[1].Line = %d, want 4", matches[1].Line)
	}
}

func TestScanEntropyHeuristicDisabledByDefault(t *testing.T) {
	content := []byte("SECRET_TOKEN=K7mP9xQ2vL8nR4tY6wZ1aB3cD5eF0gHj\n")

	s := New()
	matches, err := s.Scan(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(matches) != 0 {
		t.Fatalf("expected 0 matches with entropy scan disabled, got %d: %+v", len(matches), matches)
	}
}

func TestScanEntropyHeuristicOptIn(t *testing.T) {
	content := []byte("NODE_ENV=production\nSECRET_TOKEN=K7mP9xQ2vL8nR4tY6wZ1aB3cD5eF0gHj\n")

	s := NewWithOptions(Options{EntropyScan: true})
	matches, err := s.Scan(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d: %+v", len(matches), matches)
	}

	m := matches[0]
	if m.Method != MethodEntropy {
		t.Errorf("Method = %q, want %q", m.Method, MethodEntropy)
	}
	if m.Key != "SECRET_TOKEN" {
		t.Errorf("Key = %q, want %q", m.Key, "SECRET_TOKEN")
	}
	if m.Line != 2 {
		t.Errorf("Line = %d, want 2", m.Line)
	}
}

func TestScanRespectsProviderAllowlist(t *testing.T) {
	content := []byte("AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\nSTRIPE_SECRET_KEY=sk_live_4eC39HqLyjWDarjtT1zdp7dc\n")

	s := NewWithOptions(Options{Providers: []string{"stripe"}})
	matches, err := s.Scan(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d: %+v", len(matches), matches)
	}
	if matches[0].Provider != "stripe" {
		t.Errorf("Provider = %q, want %q", matches[0].Provider, "stripe")
	}
}

func TestScanHandlesLinesLargerThan64KB(t *testing.T) {
	// A line larger than bufio.Scanner's default 64KB token limit (e.g. a PEM
	// certificate or JSON credentials blob assigned to one env var) must not
	// truncate the scan and hide secrets on later lines.
	hugeValue := bytes.Repeat([]byte("a"), bufio.MaxScanTokenSize+1024)
	var content bytes.Buffer
	content.WriteString("BIG_BLOB=")
	content.Write(hugeValue)
	content.WriteString("\nAWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n")

	s := New()
	matches, err := s.Scan(content.Bytes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d: %+v", len(matches), matches)
	}
	if matches[0].Line != 2 {
		t.Errorf("Line = %d, want 2", matches[0].Line)
	}
	if matches[0].Provider != "aws" {
		t.Errorf("Provider = %q, want %q", matches[0].Provider, "aws")
	}
}

func TestScanRespectsIgnoreList(t *testing.T) {
	content := []byte("AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\nSTRIPE_SECRET_KEY=sk_live_4eC39HqLyjWDarjtT1zdp7dc\n")

	s := NewWithOptions(Options{Ignore: []string{"AWS_ACCESS_KEY_ID"}})
	matches, err := s.Scan(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d: %+v", len(matches), matches)
	}
	if matches[0].Key != "STRIPE_SECRET_KEY" {
		t.Errorf("Key = %q, want %q", matches[0].Key, "STRIPE_SECRET_KEY")
	}
}
