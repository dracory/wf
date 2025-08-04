package main

import (
	"context"

	"github.com/dracory/wf"
)

// NewDataStep creates a step that sets initial data
func NewDataStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Set Initial Data")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		data["numbers"] = []int{1, 2, 3}
		return ctx, data, nil
	})
	return step
}

// NewProcessStep creates a step that processes data
func NewProcessStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Process Data")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		numbers := data["numbers"].([]int)
		for i := range numbers {
			numbers[i] *= 2
		}
		data["numbers"] = numbers
		return ctx, data, nil
	})
	return step
}

// NewSumStep creates a step that calculates the sum of processed data
func NewSumStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Calculate Sum")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		numbers := data["numbers"].([]int)
		sum := 0
		for _, num := range numbers {
			sum += num
		}
		data["sum"] = sum
		return ctx, data, nil
	})
	return step
}

// NewDependenciesDag creates a DAG with multiple dependent steps
func NewDependenciesDag() wf.DagInterface {
	dag := wf.NewDag()
	dag.SetName("Dependencies Example DAG")

	// Create steps
	dataStep := NewDataStep()
	processStep := NewProcessStep()
	sumStep := NewSumStep()

	// Add steps to DAG
	dag.RunnableAdd(dataStep, processStep, sumStep)

	// Set up dependencies
	// processStep depends on dataStep
	// sumStep depends on processStep
	dag.DependencyAdd(processStep, dataStep)
	dag.DependencyAdd(sumStep, processStep)

	return dag
}

// RunDependenciesExample runs the dependencies example
func RunDependenciesExample() (map[string]any, error) {
	dag := NewDependenciesDag()
	ctx := context.Background()
	data := make(map[string]any)
	_, data, err := dag.Run(ctx, data)
	return data, err
}
