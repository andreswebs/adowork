package main

import (
	"fmt"
	"os"
)

// handleError is the unified error handler for the CLI application.
// It classifies the error, prints a user-friendly message and suggestion to stderr,
// optionally logs technical details if DEBUG is set, and exits with code 1.
func handleError(err error) {
	if err == nil {
		return
	}

	// Classify error using helpers
	switch {
	case isAuthError(err):
		fmt.Fprintln(os.Stderr, "Error: Authentication failed. Unable to access Azure DevOps with the provided credentials.")
		fmt.Fprintln(os.Stderr, "Suggestion: Check your Personal Access Token (PAT) for validity, permissions, and expiration. Ensure it is set in the ADO_PAT environment variable.")
	case isNetworkError(err):
		fmt.Fprintln(os.Stderr, "Error: Network error. Unable to connect to Azure DevOps services.")
		fmt.Fprintln(os.Stderr, "Suggestion: Check your internet connection and verify Azure DevOps is reachable. Retry after a few moments.")
	case isValidationError(err):
		fmt.Fprintln(os.Stderr, "Error: Validation error. One or more input parameters are invalid.")
		fmt.Fprintln(os.Stderr, "Suggestion: Review your command-line arguments and environment variables for missing or incorrect values.")
	case isRateLimitError(err):
		fmt.Fprintln(os.Stderr, "Error: Rate limit exceeded. Too many requests sent to Azure DevOps.")
		fmt.Fprintln(os.Stderr, "Suggestion: Wait a few minutes before retrying. Consider reducing request frequency or checking your organization's API quota.")
	case isMalformedResponseError(err):
		fmt.Fprintln(os.Stderr, "Error: Malformed response. Received unexpected or invalid data from Azure DevOps.")
		fmt.Fprintln(os.Stderr, "Suggestion: Retry the operation. If the problem persists, check for Azure DevOps service issues or API changes.")
	default:
		fmt.Fprintln(os.Stderr, "Error: An unexpected error occurred.")
		fmt.Fprintln(os.Stderr, "Suggestion: Retry the operation or contact support if the issue continues.")
		// Print the actual error message for debugging and testability
		fmt.Fprintf(os.Stderr, "Details: %v\n", err)
	}

	// If DEBUG is set, print technical details
	if os.Getenv("DEBUG") != "" {
		fmt.Fprintln(os.Stderr, "--- Technical details ---")
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}

	os.Exit(1)
}
