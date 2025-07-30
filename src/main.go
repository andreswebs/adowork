package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/urfave/cli/v3"
)

// errorHandler is a function variable for error handling, allowing test injection.
var (
	errorHandlerMu sync.RWMutex
	errorHandler   = handleError
)

// SetErrorHandler safely sets the global errorHandler.
func SetErrorHandler(f func(error)) {
	errorHandlerMu.Lock()
	defer errorHandlerMu.Unlock()
	errorHandler = f
}

// GetErrorHandler safely retrieves the global errorHandler.
func GetErrorHandler() func(error) {
	errorHandlerMu.RLock()
	defer errorHandlerMu.RUnlock()
	return errorHandler
}

// ADOClientInterface defines the methods we use from ADOClient, allowing for mocking.
type ADOClientInterface interface {
	BuildWorkItemPatchDocument(title, description string, parentID *int, assignedTo string) ([]webapi.JsonPatchOperation, error)
	CreateWorkItem(ctx context.Context, workItemType string, patchDoc []webapi.JsonPatchOperation) (*workitemtracking.WorkItem, error)
	GetWorkItemURL(workItemID int) string
}

// containsString returns true if a string is in a slice.
func containsString(slice []string, s string) bool {
	return slices.Contains(slice, s)
}

// isValidWorkItemType checks if the given type is a valid, common Azure DevOps work item type.
func isValidWorkItemType(witType string) bool {
	validTypes := []string{"task", "bug", "user story", "feature", "epic", "issue"}
	normalizedType := strings.ToLower(witType)
	return containsString(validTypes, normalizedType)
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		GetErrorHandler()(err)
	}
	cmd := &cli.Command{
		Name:    "adowork",
		Usage:   "A command-line tool for creating Azure DevOps work items",
		Version: "0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "type", Aliases: []string{"t"}, Required: true},
			&cli.StringFlag{Name: "title", Aliases: []string{"T"}, Required: true},
			&cli.StringFlag{Name: "description", Aliases: []string{"d"}},
			&cli.StringFlag{Name: "assigned-to", Aliases: []string{"a"}},
			&cli.IntFlag{Name: "parent", Aliases: []string{"p"}},
			&cli.BoolFlag{Name: "dry-run", Aliases: []string{"n"}},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return actionDispatch(ctx, cmd, &cfg)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		// If no arguments, display help text
		if len(os.Args) != 1 {
			GetErrorHandler()(err)
		}
	}
}

func actionDispatch(ctx context.Context, cmd *cli.Command, cfg *Config) error {
	client, err := NewADOClient(cfg)
	if err != nil {
		GetErrorHandler()(FormatADOError(err, "creating ADO client"))
	}

	return actionWithClient(ctx, cmd, client)
}

func actionWithClient(ctx context.Context, cmd *cli.Command, client ADOClientInterface) error {
	typeVal := cmd.String("type")
	if !isValidWorkItemType(typeVal) {
		GetErrorHandler()(fmt.Errorf("Invalid work item type: '%s'. Please use a common type like 'Task', 'Bug', or 'User Story'.", typeVal))
	}
	titleVal := cmd.String("title")
	descVal := cmd.String("description")
	assignedToVal := cmd.String("assigned-to")
	parentVal := cmd.Int("parent")
	dryRunVal := cmd.Bool("dry-run")

	var parentID *int
	if parentVal != 0 {
		parentID = &parentVal
	}

	patchDoc, err := client.BuildWorkItemPatchDocument(titleVal, descVal, parentID, assignedToVal)
	if err != nil {
		GetErrorHandler()(FormatADOError(err, "building work item patch document"))
	}

	if dryRunVal {
		fmt.Println("--- Dry Run: Work Item Payload ---")
		jsonBytes, err := json.MarshalIndent(patchDoc, "", "  ")
		if err != nil {
			GetErrorHandler()(fmt.Errorf("Error marshaling dry-run output: %v", err))
		}
		fmt.Println(string(jsonBytes))
		fmt.Println("------------------------------------")
		return nil
	}

	workItem, err := client.CreateWorkItem(ctx, typeVal, patchDoc)
	if err != nil {
		GetErrorHandler()(FormatADOError(err, "creating work item"))
	}

	if workItem == nil || workItem.Id == nil {
		GetErrorHandler()(fmt.Errorf("Failed to create work item: received no ID from API"))
	}

	fmt.Print(client.GetWorkItemURL(*workItem.Id))

	return nil
}
