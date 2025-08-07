package wf

import (
	"context"
	"testing"
)

func Test_Step_Basic(t *testing.T) {
	// Create a step with options
	customID := "custom-step-id"
	step := NewStep()

	// Test basic properties
	step.SetName("TestStep")
	step.SetID(customID)

	// Test basic properties
	if step.GetName() != "TestStep" {
		t.Errorf("Expected name TestStep, got %s", step.GetName())
	}

	// Test ID was set
	if step.GetID() != customID {
		t.Errorf("Expected ID %s, got %s", customID, step.GetID())
	}

	// Test handler execution
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		data["test"] = true
		return ctx, data, nil
	})

	// Execute step
	ctx := context.Background()
	data := make(map[string]any)
	_, data, err := step.Run(ctx, data)
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}

	if !data["test"].(bool) {
		t.Errorf("Expected test data to be true, got %v", data["test"])
	}
}

func Test_Step_WithOptions(t *testing.T) {
	// Create a step with options
	customID := "custom-step-id"
	step := NewStep(
		WithName("TestStep"),
		WithID(customID),
	)

	// Test basic properties
	if step.GetName() != "TestStep" {
		t.Errorf("Expected name TestStep, got %s", step.GetName())
	}

	// Test ID was set
	if step.GetID() != customID {
		t.Errorf("Expected ID %s, got %s", customID, step.GetID())
	}

	// Test handler execution
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		data["test"] = true
		return ctx, data, nil
	})

	// Execute step
	ctx := context.Background()
	data := make(map[string]any)
	_, data, err := step.Run(ctx, data)
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}

	if !data["test"].(bool) {
		t.Errorf("Expected test data to be true, got %v", data["test"])
	}
}
