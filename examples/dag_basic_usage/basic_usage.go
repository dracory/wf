package basic_usage

import (
	"context"

	"github.com/dracory/wf"
)

// NewSetValueStep creates a new step that sets a value
func NewSetValueStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Set Value")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		data["value"] = 42
		return ctx, data, nil
	})
	return step
}

// NewIncrementStep creates a new step that increments a value
func NewIncrementStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Increment Value")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		value := data["value"].(int)
		value++
		data["value"] = value
		return ctx, data, nil
	})
	return step
}

// NewMultipleIncrementDag creates a DAG with multiple increment steps
func NewMultipleIncrementDag() wf.DagInterface {
	dag := wf.NewDag()
	dag.SetName("Multiple Increment DAG")
	
	// Add 4 increment steps
	for i := 0; i < 4; i++ {
		dag.RunnableAdd(NewIncrementStep())
	}
	
	return dag
}
