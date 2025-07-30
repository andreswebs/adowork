package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/urfave/cli/v3"
)

// mockADOClient is a mock implementation of the ADOClient for testing purposes.
type mockADOClient struct {
	CreateWorkItemFunc func(ctx context.Context, workItemType string, patchDoc []webapi.JsonPatchOperation) (*workitemtracking.WorkItem, error)
}

// BuildWorkItemPatchDocument is a mock implementation.
func (m *mockADOClient) BuildWorkItemPatchDocument(title, description string, parentID *int, assignedTo string) ([]webapi.JsonPatchOperation, error) {
	return []webapi.JsonPatchOperation{}, nil
}

// CreateWorkItem is a mock implementation.
func (m *mockADOClient) CreateWorkItem(ctx context.Context, workItemType string, patchDoc []webapi.JsonPatchOperation) (*workitemtracking.WorkItem, error) {
	if m.CreateWorkItemFunc != nil {
		return m.CreateWorkItemFunc(ctx, workItemType, patchDoc)
	}
	return nil, errors.New("CreateWorkItemFunc not implemented")
}

// GetWorkItemURL is a mock implementation.
func (m *mockADOClient) GetWorkItemURL(workItemID int) string {
	return fmt.Sprintf("https://dev.azure.com/mock-org/mock-project/_workitems/edit/%d", workItemID)
}

func TestAction_Success(t *testing.T) {
	origHandler := GetErrorHandler()
	SetErrorHandler(func(err error) {
		panic(err)
	})
	t.Cleanup(func() { SetErrorHandler(origHandler) })

	mockClient := &mockADOClient{
		CreateWorkItemFunc: func(ctx context.Context, workItemType string, patchDoc []webapi.JsonPatchOperation) (*workitemtracking.WorkItem, error) {
			id := 123
			return &workitemtracking.WorkItem{Id: &id}, nil
		},
	}

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "type"},
			&cli.StringFlag{Name: "title"},
			&cli.StringFlag{Name: "description"},
			&cli.StringFlag{Name: "assigned-to"},
			&cli.IntFlag{Name: "parent"},
			&cli.BoolFlag{Name: "dry-run"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return actionWithClient(ctx, cmd, mockClient)
		},
	}

	err := cmd.Run(context.Background(), []string{"", "--type", "Task", "--title", "Test Task"})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestAction_CreateWorkItemError(t *testing.T) {
	origHandler := GetErrorHandler()
	SetErrorHandler(func(err error) {
		panic(err)
	})
	t.Cleanup(func() { SetErrorHandler(origHandler) })

	apiError := errors.New("API call failed")
	mockClient := &mockADOClient{
		CreateWorkItemFunc: func(ctx context.Context, workItemType string, patchDoc []webapi.JsonPatchOperation) (*workitemtracking.WorkItem, error) {
			return nil, apiError
		},
	}

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "type"},
			&cli.StringFlag{Name: "title"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return actionWithClient(ctx, cmd, mockClient)
		},
	}

	var recovered any
	func() {
		defer func() { recovered = recover() }()
		_ = cmd.Run(context.Background(), []string{"", "--type", "Bug", "--title", "Test Bug"})
	}()
	if recovered == nil {
		t.Errorf("Expected an error, but got none")
		return
	}
	if err, ok := recovered.(error); ok {
		if !strings.Contains(err.Error(), "API call failed") {
			t.Errorf("Expected error message to contain 'API call failed', but got '%s'", err.Error())
		}
	} else {
		t.Errorf("Expected error, got %T", recovered)
	}
}

func TestAction_InvalidWorkItemType(t *testing.T) {
	origHandler := GetErrorHandler()
	SetErrorHandler(func(err error) {
		panic(err)
	})
	t.Cleanup(func() { SetErrorHandler(origHandler) })

	mockClient := &mockADOClient{} // No functions needed as it should fail before client interaction.

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "type"},
			&cli.StringFlag{Name: "title"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return actionWithClient(ctx, cmd, mockClient)
		},
	}

	var recovered any
	func() {
		defer func() { recovered = recover() }()
		_ = cmd.Run(context.Background(), []string{"", "--type", "InvalidType", "--title", "Test"})
	}()
	if recovered == nil {
		t.Errorf("Expected an error for invalid work item type, but got none")
		return
	}
	if err, ok := recovered.(error); ok {
		expectedMsg := "Invalid work item type: 'InvalidType'. Please use a common type like 'Task', 'Bug', or 'User Story'."
		if err.Error() != expectedMsg {
			t.Errorf("Expected error message '%s', but got '%s'", expectedMsg, err.Error())
		}
	} else {
		t.Errorf("Expected error, got %T", recovered)
	}
}

// TestFormatError ensures FormatError produces the expected error message.
func TestFormatError(t *testing.T) {
	wrapped := FormatADOError(errors.New("API call failed"), "creating work item")
	actual := wrapped.Error()
	if !strings.Contains(actual, "API call failed") {
		t.Errorf("Expected error message to contain 'API call failed', got '%s'", actual)
	}
	if !strings.Contains(actual, "creating work item failed") {
		t.Errorf("Expected error message to contain 'creating work item failed', got '%s'", actual)
	}
}

// TestCLIErrorWithExec runs the CLI as a subprocess to test error output and exit code.
func TestCLIErrorWithExec(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("os/exec CLI test skipped on Windows")
	}
	bin := "./adowork"
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build CLI binary: %v\n%s", err, string(out))
	}
	cmd := exec.Command(bin, "--type", "Bug", "--title", "Test Bug")
	cmd.Dir = "."
	cmd.Env = append(os.Environ(),
		"ADO_ORG=dummy",
		"ADO_PROJECT=dummy",
		"ADO_PAT=dummy",
	)
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("Expected non-zero exit code, got nil error")
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 1 {
			t.Errorf("Expected exit code 1, got %d", exitErr.ExitCode())
		}
	} else {
		t.Errorf("Expected exec.ExitError, got %T", err)
	}
	outStr := string(output)
	if !strings.Contains(outStr, "creating work item failed") &&
		!strings.Contains(outStr, "An unexpected error occurred") &&
		!strings.Contains(outStr, "Validation error") &&
		!strings.Contains(outStr, "Authentication failed. Unable to access Azure DevOps with the provided credentials.") {
		t.Errorf("Expected output to contain a user-facing error message, got: %s", outStr)
	}
	_ = exec.Command("rm", bin).Run()
}
