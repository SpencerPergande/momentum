// Package config reads service settings from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// String returns the value of key, or def when it is unset or empty.
func String(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

// Required returns the value of key, or an error when it is unset or empty.
func Required(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return "", fmt.Errorf("config: %s must be set", key)
	}
	return v, nil
}

// Bool returns key parsed as a boolean, or def when it is unset or empty.
func Bool(key string, def bool) (bool, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def, nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("config: %s=%q is not a boolean: %w", key, v, err)
	}
	return b, nil
}

// Duration returns key parsed as a time.Duration, or def when it is unset or
// empty.
func Duration(key string, def time.Duration) (time.Duration, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("config: %s=%q is not a duration: %w", key, v, err)
	}
	return d, nil
}
