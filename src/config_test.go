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
	}
	missing := cfg.validate()
	if len(missing) != 0 {
		t.Errorf("Expected no missing fields, got %v", missing)
	}
}

func TestConfigValidate_MissingFields(t *testing.T) {
	cfg := Config{}
	missing := cfg.validate()
	want := []string{EnvADOOrg, EnvADOProject, EnvADOPAT}
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
	missing := cfg.validate()
	if len(missing) != 2 {
		t.Errorf("Expected 2 missing fields, got %d", len(missing))
	}
	if missing[0] != EnvADOProject && missing[1] != EnvADOProject {
		t.Errorf("Expected missing field %q", EnvADOProject)
	}
	if missing[0] != EnvADOPAT && missing[1] != EnvADOPAT {
		t.Errorf("Expected missing field %q", EnvADOPAT)
	}
}
