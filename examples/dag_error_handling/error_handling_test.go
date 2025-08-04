package main

import (
	"context"
	"testing"
)

func TestErrorHandling(t *testing.T) {
	// Create and run the DAG
	dag := NewErrorHandlingDag()

	_, data, err := dag.Run(context.Background(), map[string]any{})
	if err == nil {
		t.Error("Expected error, got nil")
		return
	}

	// Verify the error message
	if err.Error() != "intentional error" {
		t.Errorf("Expected error 'intentional error', got '%v'", err)
	}

	// Verify the value was still processed
	value, ok := data["value"].(int)
	if !ok {
		t.Error("Expected value to be an integer")
		return
	}
	if value != 2 {
		t.Errorf("Expected value 2, got %d", value)
	}
}
