package main

import (
	"context"
	"fmt"

	"github.com/dracory/wf"
)

// NewInitialValueStep creates a step that sets an initial value
func NewInitialValueStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Set Initial Value")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		data["value"] = 1
		return ctx, data, nil
	})
	return step
}

// NewProcessDataStep creates a step that processes data
func NewProcessDataStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Process Data")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		if value, ok := data["value"].(int); !ok {
			return ctx, data, fmt.Errorf("value not found")
		} else {
			data["value"] = value * 2
			return ctx, data, nil
		}
	})
	return step
}

// NewErrorStep creates a step that intentionally fails
func NewErrorStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Intentional Error")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		return ctx, data, fmt.Errorf("intentional error")
	})
	return step
}

// NewErrorHandlingDag creates a DAG with error handling
func NewErrorHandlingDag() wf.DagInterface {
	dag := wf.NewDag()
	dag.SetName("Error Handling Example DAG")

	// Create steps
	initialStep := NewInitialValueStep()
	processStep := NewProcessDataStep()
	errorStep := NewErrorStep()

	// Add steps to DAG
	dag.RunnableAdd(initialStep, processStep, errorStep)

	// Set up dependencies
	dag.DependencyAdd(processStep, initialStep)
	dag.DependencyAdd(errorStep, processStep)

	return dag
}

// RunErrorHandlingExample runs the error handling example
func RunErrorHandlingExample() (map[string]any, error) {
	dag := NewErrorHandlingDag()
	ctx := context.Background()
	data := make(map[string]any)
	_, data, err := dag.Run(ctx, data)
	return data, err
}
