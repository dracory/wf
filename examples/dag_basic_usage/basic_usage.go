package basic_usage

import (
	"context"

	"github.com/dracory/wf"
)

// NewSetValueStep creates a new step that sets a value
func NewSetValueStep() wf.StepInterface {
	return wf.NewStep(
		wf.WithName("Set Value"),
		wf.WithID("set-value-step"),
		wf.WithHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
			data["value"] = 42
			return ctx, data, nil
		}),
	)
}

// NewIncrementStep creates a new step that increments a value
func NewIncrementStep() wf.StepInterface {
	return wf.NewStep(
		wf.WithName("Increment Value"),
		wf.WithID("increment-step"),
		wf.WithHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
			value := data["value"].(int)
			value++
			data["value"] = value
			return ctx, data, nil
		}),
	)
}

// NewMultipleIncrementDag creates a DAG with multiple increment steps
func NewMultipleIncrementDag() wf.DagInterface {
	dag := wf.NewDag(
		wf.WithName("Multiple Increment DAG"),
		wf.WithID("multiple-increment-dag"),
	)
	
	// Add 4 increment steps
	for i := 0; i < 4; i++ {
		dag.RunnableAdd(NewIncrementStep())
	}
	
	return dag
}
