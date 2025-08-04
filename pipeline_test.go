package wf

import (
	"context"
	"fmt"
	"testing"
)

func Test_Pipeline_Basic(t *testing.T) {
	// Create a simple pipeline with two steps
	pipeline := NewPipeline()
	step1 := NewStep()
	step2 := NewStep()
	step1.SetName("Step1")
	step2.SetName("Step2")

	// Set handlers for steps
	step1.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		return ctx, data, nil
	})
	step2.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		return ctx, data, nil
	})

	// Add steps to pipeline
	pipeline.RunnableAdd(step1, step2)

	// Test basic pipeline structure
	if len(pipeline.RunnableList()) != 2 {
		t.Errorf("Expected 2 runnables, got %d", len(pipeline.RunnableList()))
	}

	// Test execution order
	ctx := context.Background()
	_, _, err := pipeline.Run(ctx, make(map[string]any))
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}
}

func Test_Pipeline_Remove(t *testing.T) {
	// Create a pipeline with steps
	pipeline := NewPipeline()
	step1 := NewStep()
	step2 := NewStep()
	step1.SetName("Step1")
	step2.SetName("Step2")

	// Set handlers for steps
	step1.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		return ctx, data, nil
	})
	step2.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		return ctx, data, nil
	})

	// Add steps to pipeline
	pipeline.RunnableAdd(step1, step2)

	// Remove step1
	removed := pipeline.RunnableRemove(step1)
	if !removed {
		t.Errorf("Failed to remove step1")
	}

	// Verify step1 is removed
	if len(pipeline.RunnableList()) != 1 {
		t.Errorf("Expected 1 runnable after removal, got %d", len(pipeline.RunnableList()))
	}
}

func Test_Pipeline_Execution(t *testing.T) {
	// Create a pipeline with steps that modify data
	pipeline := NewPipeline()
	step1 := NewStep()
	step2 := NewStep()
	step1.SetName("Step1")
	step2.SetName("Step2")

	// Set handlers for steps
	step1.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		data["step1"] = true
		return ctx, data, nil
	})
	step2.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		data["step2"] = true
		return ctx, data, nil
	})

	// Add steps to pipeline
	pipeline.RunnableAdd(step1, step2)

	// Test execution and data modification
	ctx := context.Background()
	data := make(map[string]any)
	_, data, err := pipeline.Run(ctx, data)
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}

	if !data["step1"].(bool) || !data["step2"].(bool) {
		t.Errorf("Expected both steps to modify data, got: %v", data)
	}
}

func Test_Pipeline_Empty(t *testing.T) {
	// Create an empty pipeline
	pipeline := NewPipeline()

	// Test execution of empty pipeline
	ctx := context.Background()
	data := make(map[string]any)
	_, data, err := pipeline.Run(ctx, data)
	if err != nil {
		t.Errorf("Run failed for empty pipeline: %v", err)
	}

	// Verify data is unchanged
	if len(data) != 0 {
		t.Errorf("Expected empty data for empty pipeline, got: %v", data)
	}
}

func Test_Pipeline_ErrorPropagation(t *testing.T) {
	// Create a pipeline with a step that returns an error
	pipeline := NewPipeline()
	step1 := NewStep()
	step1.SetName("Step1")

	// Set handler that returns an error
	step1.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		return ctx, data, fmt.Errorf("step1 failed")
	})

	// Add step to pipeline
	pipeline.RunnableAdd(step1)

	// Test error propagation
	ctx := context.Background()
	_, _, err := pipeline.Run(ctx, make(map[string]any))
	if err == nil {
		t.Errorf("Expected error from step1, got nil")
	} else if err.Error() != "step1 failed" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func Test_Pipeline_LargeSequence(t *testing.T) {
	// Create a pipeline with 100 increment steps
	pipeline := NewPipeline()
	steps := make([]RunnableInterface, 100)

	// Create and add steps
	for i := range 100 {
		step := NewStep()
		step.SetName(fmt.Sprintf("Step%d", i))
		step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
			current, ok := data["counter"]
			if !ok {
				current = 0
			}
			data["counter"] = current.(int) + 1
			return ctx, data, nil
		})
		steps[i] = step
	}

	// Add all steps to pipeline
	pipeline.RunnableAdd(steps...)

	// Test execution
	ctx := context.Background()
	data := make(map[string]any)
	_, data, err := pipeline.Run(ctx, data)
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}

	// Verify counter incremented correctly
	counter, ok := data["counter"]
	if !ok {
		t.Errorf("Expected counter in data, got: %v", data)
	}
	if counter.(int) != 100 {
		t.Errorf("Expected counter to be 100, got: %d", counter)
	}
}
