package rapid

import (
	"testing"
	"time"
)

func TestLoadCmdlineDefaultsUsesBaseValues(t *testing.T) {
	got := loadCmdlineDefaults(func(string) (string, bool) {
		return "", false
	})
	want := defaultCmdline()

	if got != want {
		t.Fatalf("defaults mismatch: got %+v, want %+v", got, want)
	}
}

func TestLoadCmdlineDefaultsFromEnv(t *testing.T) {
	env := map[string]string{
		"RAPID_CHECKS":     "0xc8",
		"RAPID_STEPS":      "40_000",
		"RAPID_FAILFILE":   "/tmp/failfile",
		"RAPID_NOFAILFILE": "true",
		"RAPID_SEED":       "0x1234",
		"RAPID_LOG":        "true",
		"RAPID_V":          "true",
		"RAPID_DEBUG":      "true",
		"RAPID_DEBUGVIS":   "true",
		"RAPID_SHRINKTIME": "45s",
	}

	got := loadCmdlineDefaults(func(key string) (string, bool) {
		value, ok := env[key]
		return value, ok
	})

	if got.checks != 200 {
		t.Fatalf("checks: got %d, want %d", got.checks, 200)
	}
	if got.steps != 40000 {
		t.Fatalf("steps: got %d, want %d", got.steps, 40000)
	}
	if got.failfile != "/tmp/failfile" {
		t.Fatalf("failfile: got %q, want %q", got.failfile, "/tmp/failfile")
	}
	if !got.nofailfile {
		t.Fatalf("nofailfile: expected true")
	}
	if got.seed != 0x1234 {
		t.Fatalf("seed: got %d, want %d", got.seed, 0x1234)
	}
	if !got.log || !got.verbose || !got.debug || !got.debugvis {
		t.Fatalf("expected all bool flags true, got %+v", got)
	}
	if got.shrinkTime != 45*time.Second {
		t.Fatalf("shrinkTime: got %v, want %v", got.shrinkTime, 45*time.Second)
	}
}

func TestLoadCmdlineDefaultsInvalidEnvPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for invalid env value")
		}
	}()

	loadCmdlineDefaults(func(string) (string, bool) {
		return "not-an-int", true
	})
}
