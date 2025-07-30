package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
)

func TestIsAuthError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		// All string-based cases should now return false
		{"401 string", errors.New("401 Unauthorized"), false},
		{"403 string", errors.New("403 Forbidden"), false},
		{"unauthorized msg", errors.New("user unauthorized"), false},
		{"forbidden msg", errors.New("access forbidden"), false},
		{"unrelated", errors.New("some other error"), false},
		{"wrapped unauthorized", fmt.Errorf("wrap: %w", errors.New("401 Unauthorized")), false},
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
		// All string-based cases should now return false
		{"connection refused", errors.New("connection refused"), false},
		{"timeout string", errors.New("timeout occurred"), false},
		{"network unreachable", errors.New("network unreachable"), false},
		{"unrelated", errors.New("some other error"), false},
		{"wrapped timeout", fmt.Errorf("wrap: %w", errors.New("timeout occurred")), false},
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
		// All string-based cases should now return false
		{"400 bad request", errors.New("400 Bad Request"), false},
		{"422 validation", errors.New("422 validation error"), false},
		{"validation msg", errors.New("validation failed"), false},
		{"bad request msg", errors.New("bad request"), false},
		{"unrelated", errors.New("some other error"), false},
		{"wrapped validation", fmt.Errorf("wrap: %w", errors.New("validation failed")), false},
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
		// All string-based cases should now return false
		{"429 too many requests", errors.New("429 Too Many Requests"), false},
		{"rate limit msg", errors.New("rate limit exceeded"), false},
		{"too many requests msg", errors.New("too many requests"), false},
		{"unrelated", errors.New("some other error"), false},
		{"wrapped rate limit", fmt.Errorf("wrap: %w", errors.New("rate limit exceeded")), false},
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
	syntaxErr := &json.SyntaxError{}
	typeErr := &json.UnmarshalTypeError{}
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"syntax error", syntaxErr, true},
		{"unmarshal type error", typeErr, true},
		// All string-based cases should now return false
		{"invalid character", errors.New("invalid character '}' looking for beginning of object key string"), false},
		{"unexpected end of json", errors.New("unexpected end of JSON input"), false},
		{"unmarshal msg", errors.New("json: cannot unmarshal string into Go value"), false},
		{"unrelated", errors.New("some other error"), false},
		{"wrapped unmarshal", fmt.Errorf("wrap: %w", errors.New("json: cannot unmarshal string into Go value")), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMalformedResponseError(tt.err); got != tt.want {
				t.Errorf("isMalformedResponseError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func makeWrappedError(status int) error {
	code := status
	msg := "mock error"
	return azuredevops.WrappedError{StatusCode: &code, Message: &msg}
}

func TestIsAuthError_TypeBased(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"401 WrappedError", makeWrappedError(401), true},
		{"403 WrappedError", makeWrappedError(403), true},
		{"404 WrappedError", makeWrappedError(404), false},
		{"nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAuthError(tt.err); got != tt.want {
				t.Errorf("isAuthError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestIsValidationError_TypeBased(t *testing.T) {
	verr := &azuredevops.InvalidVersionStringError{}
	apiErr := azuredevops.InvalidApiVersion{}
	locErr := azuredevops.LocationIdNotRegisteredError{}
	badReq := makeWrappedError(400)
	val422 := makeWrappedError(422)
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"InvalidVersionStringError", verr, true},
		{"InvalidApiVersion", apiErr, true},
		{"LocationIdNotRegisteredError", locErr, true},
		{"400 WrappedError", badReq, true},
		{"422 WrappedError", val422, true},
		{"404 WrappedError", makeWrappedError(404), false},
		{"nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidationError(tt.err); got != tt.want {
				t.Errorf("isValidationError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestIsRateLimitError_TypeBased(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"429 WrappedError", makeWrappedError(429), true},
		{"401 WrappedError", makeWrappedError(401), false},
		{"nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRateLimitError(tt.err); got != tt.want {
				t.Errorf("isRateLimitError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
