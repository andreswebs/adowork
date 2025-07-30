package main

import (
	"os"
	"slices"
	"testing"
)

func TestReadConfigFromEnv(t *testing.T) {
	os.Setenv(EnvADOOrg, "test_org")
	os.Setenv(EnvADOProject, "test_project")
	os.Setenv(EnvADOPAT, "test_token")

	defer os.Unsetenv(EnvADOOrg)
	defer os.Unsetenv(EnvADOProject)
	defer os.Unsetenv(EnvADOPAT)

	cfg := readConfigFromEnv()

	if cfg.Organization != "test_org" {
		t.Errorf("Org: got %q, want %q", cfg.Organization, "test_org")
	}
	if cfg.Project != "test_project" {
		t.Errorf("Project: got %q, want %q", cfg.Project, "test_project")
	}
	if cfg.PAT != "test_token" {
		t.Errorf("PAT: got %q, want %q", cfg.PAT, "test_token")
	}
}

func TestConfigValidate_AllPresent(t *testing.T) {
	cfg := Config{
		Organization: "test_org",
		Project:      "test_proj",
		PAT:          "test_pat",
		BaseURL:      "test_url",
	}
	_, err := cfg.checkMissing()

	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestConfigValidate_MissingFields(t *testing.T) {
	cfg := Config{}
	missing, _ := cfg.checkMissing()
	want := []string{EnvADOOrg, EnvADOProject, EnvADOPAT, EnvADOBaseURL}
	if len(missing) != len(want) {
		t.Errorf("Expected %d missing fields, got %d", len(want), len(missing))
	}
	for _, field := range want {
		found := slices.Contains(missing, field)
		if !found {
			t.Errorf("Missing field %q not found in result", field)
		}
	}
}

func TestConfigValidate_SomeMissing(t *testing.T) {
	cfg := Config{Organization: "org"}
	missing, _ := cfg.checkMissing()
	if len(missing) != 3 {
		t.Errorf("Expected 3 missing fields, got %d", len(missing))
	}
	if !slices.Contains(missing, EnvADOProject) {
		t.Errorf("Expected missing field %q", EnvADOProject)
	}
	if !slices.Contains(missing, EnvADOPAT) {
		t.Errorf("Expected missing field %q", EnvADOPAT)
	}
	if !slices.Contains(missing, EnvADOBaseURL) {
		t.Errorf("Expected missing field %q", EnvADOBaseURL)
	}
}
