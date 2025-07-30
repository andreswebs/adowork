package main

import (
	"fmt"
	"os"
	"regexp"
)

const (
	EnvADOOrg         string = "ADO_ORG"
	EnvADOProject     string = "ADO_PROJECT"
	EnvADOPAT         string = "ADO_PAT"
	EnvADOBaseURL     string = "ADO_BASE_URL"
	defaultADOBaseURL string = "https://dev.azure.com"
)

type Config struct {
	Organization string `validate:"required" env:"ADO_ORG"`
	Project      string `validate:"required" env:"ADO_PROJECT"`
	PAT          string `validate:"required" env:"ADO_PAT"`
	BaseURL      string `env:"ADO_BASE_URL"`
}

// readConfigFromEnv reads ADO_* environment variables and returns a Config struct.
func readConfigFromEnv() Config {
	c := Config{
		Organization: os.Getenv(EnvADOOrg),
		Project:      os.Getenv(EnvADOProject),
		PAT:          os.Getenv(EnvADOPAT),
		BaseURL:      os.Getenv(EnvADOBaseURL),
	}

	c.BaseURL = c.normalizeBaseURL()

	return c
}

// normalizeBaseURL returns the base URL for Azure DevOps, ensuring it does not end with a slash.
func (c *Config) normalizeBaseURL() string {
	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = defaultADOBaseURL
	} else {
		re := regexp.MustCompile(`/+$`)
		baseURL = re.ReplaceAllString(baseURL, "")
	}
	return baseURL
}

// checkMissing checks that all required fields in Config are non-empty.
// Returns an error if any are missing.
func (c *Config) checkMissing() (missing []string, err error) {
	if c.Organization == "" {
		missing = append(missing, EnvADOOrg)
	}
	if c.Project == "" {
		missing = append(missing, EnvADOProject)
	}
	if c.PAT == "" {
		missing = append(missing, EnvADOPAT)
	}
	if c.BaseURL == "" {
		missing = append(missing, EnvADOBaseURL)
	}
	err = formatMissingEnvError(missing)
	return
}

// formatMissingEnvError returns a grouped, user-friendly error message for missing env vars.
func formatMissingEnvError(missing []string) error {
	if len(missing) == 0 {
		return nil
	}
	msg := "Missing required environment variables:\n"
	for _, env := range missing {
		msg += "  - " + env + "\n"
	}
	msg += "\nPlease set the above variables in your environment. Example (bash/zsh):\n"
	for _, env := range missing {
		msg += "  export " + env + "=value\n"
	}
	return fmt.Errorf("%s", msg)
}

// loadConfig reads env vars and validates.
func loadConfig() (cfg Config, err error) {
	cfg = readConfigFromEnv()
	_, err = cfg.checkMissing()
	return
}
