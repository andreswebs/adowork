package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
)

func TestIsAuthError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"401 string", errors.New("401 Unauthorized"), true},
		{"403 string", errors.New("403 Forbidden"), true},
		{"unauthorized msg", errors.New("user unauthorized"), true},
		{"forbidden msg", errors.New("access forbidden"), true},
		{"unrelated", errors.New("some other error"), false},
		{"wrapped unauthorized", fmt.Errorf("wrap: %w", errors.New("401 Unauthorized")), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAuthError(tt.err); got != tt.want {
				t.Errorf("isAuthError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestIsNetworkError(t *testing.T) {
	timeoutErr := &net.DNSError{IsTimeout: true}
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"timeout net.Error", timeoutErr, true},
		{"deadline exceeded", context.DeadlineExceeded, true},
		{"connection refused", errors.New("connection refused"), true},
		{"timeout string", errors.New("timeout occurred"), true},
		{"network unreachable", errors.New("network unreachable"), true},
		{"unrelated", errors.New("some other error"), false},
		{"wrapped timeout", fmt.Errorf("wrap: %w", errors.New("timeout occurred")), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNetworkError(tt.err); got != tt.want {
				t.Errorf("isNetworkError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"400 bad request", errors.New("400 Bad Request"), true},
		{"422 validation", errors.New("422 validation error"), true},
		{"validation msg", errors.New("validation failed"), true},
		{"bad request msg", errors.New("bad request"), true},
		{"unrelated", errors.New("some other error"), false},
		{"wrapped validation", fmt.Errorf("wrap: %w", errors.New("validation failed")), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidationError(tt.err); got != tt.want {
				t.Errorf("isValidationError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestIsRateLimitError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"429 too many requests", errors.New("429 Too Many Requests"), true},
		{"rate limit msg", errors.New("rate limit exceeded"), true},
		{"too many requests msg", errors.New("too many requests"), true},
		{"unrelated", errors.New("some other error"), false},
		{"wrapped rate limit", fmt.Errorf("wrap: %w", errors.New("rate limit exceeded")), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRateLimitError(tt.err); got != tt.want {
				t.Errorf("isRateLimitError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestIsMalformedResponseError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"invalid character", errors.New("invalid character '}' looking for beginning of object key string"), true},
		{"unexpected end of json", errors.New("unexpected end of JSON input"), true},
		{"unmarshal msg", errors.New("json: cannot unmarshal string into Go value"), true},
		{"unrelated", errors.New("some other error"), false},
		{"wrapped unmarshal", fmt.Errorf("wrap: %w", errors.New("json: cannot unmarshal string into Go value")), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMalformedResponseError(tt.err); got != tt.want {
				t.Errorf("isMalformedResponseError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
