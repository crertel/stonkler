package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunConfigShowHidesAPIKey(t *testing.T) {
	var stdout, stderr bytes.Buffer
	getenv := func(key string) string {
		if key == "FMP_API_KEY" {
			return "secret-value"
		}
		return ""
	}

	code := runConfig([]string{"show"}, &stdout, &stderr, getenv)

	if code != 0 {
		t.Fatalf("runConfig() code = %d, want 0", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), "fmp.api_key_configured=true") {
		t.Fatalf("stdout = %q, want configured key state", stdout.String())
	}
	if strings.Contains(stdout.String(), "secret-value") {
		t.Fatalf("stdout leaked API key: %q", stdout.String())
	}
}

func TestRunConfigDoctorMissingKey(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runConfig([]string{"doctor"}, &stdout, &stderr, func(string) string { return "" })

	if code != 1 {
		t.Fatalf("runConfig() code = %d, want 1", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), "FMP_API_KEY is not configured") {
		t.Fatalf("stdout = %q, want missing key diagnostic", stdout.String())
	}
}

func TestRunConfigProviders(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := runConfig([]string{"providers"}, &stdout, &stderr, func(string) string { return "secret-value" })

	if code != 0 {
		t.Fatalf("runConfig() code = %d, want 0", code)
	}
	if got := stdout.String(); got != "fmp\tready\n" {
		t.Fatalf("stdout = %q, want provider readiness", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}
