package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
)

// ADOClient is the Azure DevOps API client structure using the official library.
type ADOClient struct {
	Organization string
	Project      string
	PAT          string
	BaseURL      string
	Connection   *azuredevops.Connection
	WITClient    workitemtracking.Client
}

// NewADOClient creates a new ADOClient using the official azure-devops-go-api library.
func NewADOClient(c *Config) (*ADOClient, error) {
	_, err := c.checkMissing()
	if err != nil {
		return nil, err
	}

	organizationURL := fmt.Sprintf("%s/%s", c.BaseURL, c.Organization)
	connection := azuredevops.NewPatConnection(organizationURL, c.PAT)
	ctx := context.Background()
	witClient, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		return nil, FormatADOError(err, "Creating work item tracking client")
	}

	return &ADOClient{
		Organization: c.Organization,
		Project:      c.Project,
		PAT:          c.PAT,
		BaseURL:      c.BaseURL,
		Connection:   connection,
		WITClient:    witClient,
	}, nil
}

// BuildWorkItemPatchDocument constructs the JSON patch document for creating or updating a work item.
func (c *ADOClient) BuildWorkItemPatchDocument(title, description string, parentID *int, assignedTo string) ([]webapi.JsonPatchOperation, error) {
	var patchDoc []webapi.JsonPatchOperation

	// Add title field
	patchDoc = append(patchDoc, webapi.JsonPatchOperation{
		Op:    &webapi.OperationValues.Add,
		Path:  stringPtr("/fields/System.Title"),
		Value: title,
	})

	// Add description field if provided
	if description != "" {
		patchDoc = append(patchDoc, webapi.JsonPatchOperation{
			Op:    &webapi.OperationValues.Add,
			Path:  stringPtr("/fields/System.Description"),
			Value: description,
		})
	}

	// Add assigned to field if provided
	if assignedTo != "" {
		patchDoc = append(patchDoc, webapi.JsonPatchOperation{
			Op:    &webapi.OperationValues.Add,
			Path:  stringPtr("/fields/System.AssignedTo"),
			Value: assignedTo,
		})
	}

	// Add parent relationship if specified
	if parentID != nil {
		patchDoc = append(patchDoc, webapi.JsonPatchOperation{
			Op:   &webapi.OperationValues.Add,
			Path: stringPtr("/relations/-"),
			Value: map[string]interface{}{
				"rel": "System.LinkTypes.Hierarchy-Reverse",
				"url": fmt.Sprintf("%s/%s/%s/_apis/wit/workItems/%d",
					c.BaseURL, c.Organization, c.Project, *parentID),
			},
		})
	}

	return patchDoc, nil
}

// CreateWorkItem creates a new work item using the official Azure DevOps Go API library.
func (c *ADOClient) CreateWorkItem(ctx context.Context, workItemType string, patchDoc []webapi.JsonPatchOperation) (*workitemtracking.WorkItem, error) {
	// Create the work item using the typed client
	args := workitemtracking.CreateWorkItemArgs{
		Document: &patchDoc,
		Project:  &c.Project,
		Type:     &workItemType,
	}

	workItem, err := c.WITClient.CreateWorkItem(ctx, args)
	if err != nil {
		return nil, FormatADOError(err, "Creating work item")
	}

	return workItem, nil
}

// GetWorkItemURL returns the URL for accessing a work item in the Azure DevOps web interface.
func (c *ADOClient) GetWorkItemURL(workItemID int) string {
	return fmt.Sprintf("%s/%s/%s/_workitems/edit/%d",
		c.BaseURL, c.Organization, c.Project, workItemID)
}

// stringPtr is a helper function to return a pointer to a string.
func stringPtr(s string) *string {
	return &s
}

// Error handling utilities for Azure DevOps API errors

// IsArgumentError checks if an error is an argument validation error from the Azure DevOps library
func IsArgumentError(err error) bool {
	var argNilErr *azuredevops.ArgumentNilError
	var argNilOrEmptyErr *azuredevops.ArgumentNilOrEmptyError
	return errors.As(err, &argNilErr) || errors.As(err, &argNilOrEmptyErr)
}

// IsAPIError checks if an error is an API error from the Azure DevOps library
func IsAPIError(err error) bool {
	var wrappedErr *azuredevops.WrappedError
	return errors.As(err, &wrappedErr)
}

// GetAPIErrorDetails extracts detailed error information from Azure DevOps API errors
func GetAPIErrorDetails(err error) (statusCode int, message string, details map[string]interface{}) {
	var wrappedErr *azuredevops.WrappedError
	if errors.As(err, &wrappedErr) {
		if wrappedErr.StatusCode != nil {
			statusCode = *wrappedErr.StatusCode
		}
		if wrappedErr.Message != nil {
			message = *wrappedErr.Message
		}
		if wrappedErr.CustomProperties != nil {
			details = *wrappedErr.CustomProperties
		}
		return statusCode, message, details
	}
	return 0, "", nil
}

// FormatADOError provides a user-friendly error message with context from Azure DevOps API errors
func FormatADOError(err error, operation string) error {
	if IsArgumentError(err) {
		return fmt.Errorf("%s failed due to invalid arguments: %w", operation, err)
	}

	if IsAPIError(err) {
		statusCode, message, _ := GetAPIErrorDetails(err)
		if statusCode != 0 && message != "" {
			return fmt.Errorf("%s failed (HTTP %d): %s", operation, statusCode, message)
		}
	}

	return fmt.Errorf("%s failed: %w", operation, err)
}
