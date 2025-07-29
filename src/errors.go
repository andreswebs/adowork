package main

// Error classification helpers for Azure DevOps Go SDK integration.
// Each function returns true if the error matches the category.

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"strings"
	// Azure DevOps Go SDK error types (pseudo-import, replace with actual if available)
	// "github.com/microsoft/azure-devops-go-api/azuredevops"
)

// isAuthError returns true if the error is an authentication error (e.g., 401/403, invalid PAT).
// Criteria: HTTP 401/403, error message contains 'unauthorized', 'forbidden', 'invalid personal access token', 'access denied',
// or Azure DevOps SDK-specific auth error types (e.g., AuthenticationError, IdentityError).
func isAuthError(err error) bool {
	if err == nil {
		return false
	}
	// SDK-specific: check for known auth error types (pseudo-code, replace with actual types if available)
	// if errors.As(err, &azuredevops.AuthenticationError{}) || errors.As(err, &azuredevops.IdentityError{}) {
	// 	return true
	// }
	// Check error message for common auth patterns
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "forbidden") ||
		strings.Contains(errStr, "401") || strings.Contains(errStr, "403") ||
		strings.Contains(errStr, "invalid personal access token") || strings.Contains(errStr, "access denied") {
		return true
	}
	return false
}

// isNetworkError returns true if the error is a network error (timeouts, DNS, connection reset).
// Criteria: net.Error, context.DeadlineExceeded, error message contains 'connection refused', 'timeout', 'network unreachable',
// or SDK-specific network error wrappers.
func isNetworkError(err error) bool {
	if err == nil {
		return false
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "network unreachable") || strings.Contains(errStr, "connection reset") {
		return true
	}
	return false
}

// isValidationError returns true if the error is a validation error (HTTP 400/422, field validation).
// Criteria: HTTP 400/422, error message contains 'validation', 'bad request', 'field required', 'invalid value',
// or SDK-specific validation error types (e.g., ValidationError).
func isValidationError(err error) bool {
	if err == nil {
		return false
	}
	// SDK-specific: check for known validation error types (pseudo-code)
	// if errors.As(err, &azuredevops.ValidationError{}) {
	// 	return true
	// }
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "validation") || strings.Contains(errStr, "bad request") ||
		strings.Contains(errStr, "400") || strings.Contains(errStr, "422") ||
		strings.Contains(errStr, "field required") || strings.Contains(errStr, "invalid value") {
		return true
	}
	return false
}

// isRateLimitError returns true if the error is a rate limiting error (HTTP 429, rate limit messages).
// Criteria: HTTP 429, error message contains 'rate limit', 'too many requests', 'quota exceeded',
// or SDK-specific rate limit error types (e.g., RateLimitError).
func isRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	// SDK-specific: check for known rate limit error types (pseudo-code)
	// if errors.As(err, &azuredevops.RateLimitError{}) {
	// 	return true
	// }
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "rate limit") || strings.Contains(errStr, "too many requests") ||
		strings.Contains(errStr, "429") || strings.Contains(errStr, "quota exceeded") {
		return true
	}
	return false
}

// isMalformedResponseError returns true if the error is a malformed response error (JSON unmarshal, unexpected format).
// Criteria: Go json.SyntaxError, json.UnmarshalTypeError, error message contains 'invalid character', 'unexpected end of json',
// 'unmarshal', 'malformed', or SDK-specific response parsing errors.
func isMalformedResponseError(err error) bool {
	if err == nil {
		return false
	}
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
		return true
	}
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "invalid character") || strings.Contains(errStr, "unexpected end of json") ||
		strings.Contains(errStr, "unmarshal") || strings.Contains(errStr, "malformed") {
		return true
	}
	return false
}

// Add more detailed logic and documentation as you implement each helper.
