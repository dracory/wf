package dagstate

import (
	"context"
	"fmt"
	"time"

	"github.com/dracory/wf"
)

// This example demonstrates how to use a Directed Acyclic Graph (DAG) with state management capabilities.
// The workflow simulates a data processing pipeline with the following steps:
// 1. Process Data: Initial data processing step
// 2. Validate Data: Validates the processed data
// 3. Transform Data: Transforms the validated data
// 4. Analyze Data: Performs analysis on the transformed data
// 5. Generate Report: Creates a final report from the analyzed data
//
// The example shows how to:
// - Create and configure a DAG with multiple steps
// - Add dependencies between steps to control execution order
// - Run the workflow and handle its execution
// - Pause the workflow at any point
// - Save the workflow state to JSON
// - Create a new DAG instance and restore its state
// - Resume the workflow from where it was paused
//
// This is particularly useful for long-running workflows that need to be:
// - Paused and resumed
// - Saved and restored
// - Executed across different sessions
// - Distributed across multiple machines

// NewProcessDataStep creates a step that processes data
func NewProcessDataStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Process Data")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		// Simulate data processing
		fmt.Println("Processing data...")
		time.Sleep(1 * time.Second)
		data["processed"] = true
		return ctx, data, nil
	})
	return step
}

// NewValidateDataStep creates a step that validates data
func NewValidateDataStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Validate Data")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		// Simulate data validation
		fmt.Println("Validating data...")
		time.Sleep(1 * time.Second)
		if !data["processed"].(bool) {
			return ctx, data, fmt.Errorf("data not processed")
		}
		data["validated"] = true
		return ctx, data, nil
	})
	return step
}

// NewTransformDataStep creates a step that transforms data
func NewTransformDataStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Transform Data")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		// Simulate data transformation
		fmt.Println("Transforming data...")
		time.Sleep(1 * time.Second)
		if !data["validated"].(bool) {
			return ctx, data, fmt.Errorf("data not validated")
		}
		data["transformed"] = true
		return ctx, data, nil
	})
	return step
}

// NewAnalyzeDataStep creates a step that analyzes data
func NewAnalyzeDataStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Analyze Data")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		// Simulate data analysis
		fmt.Println("Analyzing data...")
		time.Sleep(1 * time.Second)
		if !data["transformed"].(bool) {
			return ctx, data, fmt.Errorf("data not transformed")
		}
		data["analyzed"] = true
		return ctx, data, nil
	})
	return step
}

// NewGenerateReportStep creates a step that generates a report
func NewGenerateReportStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Generate Report")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		// Simulate report generation
		fmt.Println("Generating report...")
		time.Sleep(1 * time.Second)
		if !data["analyzed"].(bool) {
			return ctx, data, fmt.Errorf("data not analyzed")
		}
		data["report"] = "Final Report"
		return ctx, data, nil
	})
	return step
}

// NewDataDag creates a DAG for data processing
func NewDataDag() wf.DagInterface {
	dag := wf.NewDag()
	dag.SetName("Data Processing DAG")

	// Create steps
	processStep := NewProcessDataStep()
	validateStep := NewValidateDataStep()
	transformStep := NewTransformDataStep()
	analyzeStep := NewAnalyzeDataStep()
	reportStep := NewGenerateReportStep()

	// Add steps to DAG
	dag.RunnableAdd(processStep, validateStep, transformStep, analyzeStep, reportStep)

	// Add dependencies:
	// - Validate depends on Process
	// - Transform depends on Validate
	// - Analyze depends on Transform
	// - Report depends on Analyze
	dag.DependencyAdd(validateStep, processStep)
	dag.DependencyAdd(transformStep, validateStep)
	dag.DependencyAdd(analyzeStep, transformStep)
	dag.DependencyAdd(reportStep, analyzeStep)

	return dag
}

// RunDagStateExample demonstrates the DAG with state management
func RunDagStateExample() error {
	// Create DAG
	dag := NewDataDag()

	// Initialize data
	ctx := context.Background()
	data := map[string]any{
		"input": "test data",
	}

	// Start DAG
	ctx, data, err := dag.Run(ctx, data)
	if err != nil {
		return fmt.Errorf("DAG failed: %v", err)
	}

	// Pause the DAG after processing data
	if dag.IsRunning() {
		err = dag.Pause()
		if err != nil {
			return fmt.Errorf("failed to pause DAG: %v", err)
		}
		fmt.Println("DAG paused after processing data")
	}

	// Save DAG state
	state := dag.GetState()
	stateJSON, err := state.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to save state: %v", err)
	}
	fmt.Printf("Saved DAG state: %s\n", string(stateJSON))

	// Create a new DAG instance
	newDag := NewDataDag()

	// Load saved state
	newState := wf.NewState()
	if err := newState.FromJSON(stateJSON); err != nil {
		return fmt.Errorf("failed to load state: %v", err)
	}
	newDag.SetState(newState)

	// Resume DAG
	_, _, err = newDag.Resume(ctx, data)
	if err != nil {
		return fmt.Errorf("DAG resume failed: %v", err)
	}

	return nil
}
