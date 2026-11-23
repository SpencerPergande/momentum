package config

import (
	"testing"
	"time"
)

func TestString(t *testing.T) {
	t.Setenv("MOMENTUM_TEST_STRING", "")
	if got := String("MOMENTUM_TEST_STRING", "def"); got != "def" {
		t.Errorf("empty value: got %q, want %q", got, "def")
	}
	t.Setenv("MOMENTUM_TEST_STRING", "set")
	if got := String("MOMENTUM_TEST_STRING", "def"); got != "set" {
		t.Errorf("set value: got %q, want %q", got, "set")
	}
}

func TestRequired(t *testing.T) {
	t.Setenv("MOMENTUM_TEST_REQUIRED", "")
	if _, err := Required("MOMENTUM_TEST_REQUIRED"); err == nil {
		t.Error("empty value: want error, got nil")
	}
	t.Setenv("MOMENTUM_TEST_REQUIRED", "x")
	if got, err := Required("MOMENTUM_TEST_REQUIRED"); err != nil || got != "x" {
		t.Errorf("set value: got %q, %v", got, err)
	}
}

func TestBool(t *testing.T) {
	t.Setenv("MOMENTUM_TEST_BOOL", "")
	if got, err := Bool("MOMENTUM_TEST_BOOL", true); err != nil || !got {
		t.Errorf("default: got %v, %v", got, err)
	}
	t.Setenv("MOMENTUM_TEST_BOOL", "false")
	if got, err := Bool("MOMENTUM_TEST_BOOL", true); err != nil || got {
		t.Errorf("false: got %v, %v", got, err)
	}
	t.Setenv("MOMENTUM_TEST_BOOL", "maybe")
	if _, err := Bool("MOMENTUM_TEST_BOOL", true); err == nil {
		t.Error("invalid: want error, got nil")
	}
}

func TestDuration(t *testing.T) {
	t.Setenv("MOMENTUM_TEST_DURATION", "")
	if got, err := Duration("MOMENTUM_TEST_DURATION", time.Second); err != nil || got != time.Second {
		t.Errorf("default: got %v, %v", got, err)
	}
	t.Setenv("MOMENTUM_TEST_DURATION", "15m")
	if got, err := Duration("MOMENTUM_TEST_DURATION", time.Second); err != nil || got != 15*time.Minute {
		t.Errorf("15m: got %v, %v", got, err)
	}
	t.Setenv("MOMENTUM_TEST_DURATION", "soon")
	if _, err := Duration("MOMENTUM_TEST_DURATION", time.Second); err == nil {
		t.Error("invalid: want error, got nil")
	}
}
