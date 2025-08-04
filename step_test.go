package wf

import (
	"context"
	"testing"
)

func Test_Step_Basic(t *testing.T) {
	// Create a step
	step := NewStep()
	step.SetName("TestStep")

	// Test basic properties
	if step.GetName() != "TestStep" {
		t.Errorf("Expected name TestStep, got %s", step.GetName())
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
