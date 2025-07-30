package main

// Error classification helpers for Azure DevOps Go SDK integration.
// Each function returns true if the error matches the category.

import (
	"context"
	"encoding/json"
	"errors"
	"net"

	// Azure DevOps Go SDK error types
	"github.com/microsoft/azure-devops-go-api/azuredevops"
)

// isAuthError returns true if the error is an authentication error (e.g., 401/403, invalid PAT).
// Uses type-based checks: checks for azuredevops.WrappedError with 401/403 status codes.
func isAuthError(err error) bool {
	if err == nil {
		return false
	}
	var we azuredevops.WrappedError
	if errors.As(err, &we) && we.StatusCode != nil {
		if *we.StatusCode == 401 || *we.StatusCode == 403 {
			return true
		}
	}
	return false
}

// isNetworkError returns true if the error is a network error (timeouts, DNS, connection reset).
// Uses type-based checks only.
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
	return false
}

// isValidationError returns true if the error is a validation error (HTTP 400/422, field validation).
// Uses type-based checks for Azure DevOps SDK error types.
func isValidationError(err error) bool {
	if err == nil {
		return false
	}
	if errors.As(err, new(*azuredevops.InvalidVersionStringError)) {
		return true
	}
	if errors.As(err, new(azuredevops.InvalidApiVersion)) {
		return true
	}
	if errors.As(err, new(azuredevops.LocationIdNotRegisteredError)) {
		return true
	}
	var we azuredevops.WrappedError
	if errors.As(err, &we) && we.StatusCode != nil {
		if *we.StatusCode == 400 || *we.StatusCode == 422 {
			return true
		}
	}
	return false
}

// isRateLimitError returns true if the error is a rate limiting error (HTTP 429, rate limit messages).
// Uses type-based checks: checks for azuredevops.WrappedError with 429 status code.
func isRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	var we azuredevops.WrappedError
	if errors.As(err, &we) && we.StatusCode != nil {
		if *we.StatusCode == 429 {
			return true
		}
	}
	return false
}

// isMalformedResponseError returns true if the error is a malformed response error (JSON unmarshal, unexpected format).
// Uses type-based checks only.
func isMalformedResponseError(err error) bool {
	if err == nil {
		return false
	}
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
		return true
	}
	return false
}
