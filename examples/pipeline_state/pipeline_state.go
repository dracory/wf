package pipelinestate

import (
	"context"
	"fmt"
	"time"

	"github.com/dracory/wf"
)

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

// NewDataPipeline creates a pipeline for data processing
func NewDataPipeline() wf.PipelineInterface {
	pipeline := wf.NewPipeline()
	pipeline.SetName("Data Processing Pipeline")

	// Create steps
	processStep := NewProcessDataStep()
	validateStep := NewValidateDataStep()
	transformStep := NewTransformDataStep()

	// Add steps to pipeline
	pipeline.RunnableAdd(processStep, validateStep, transformStep)

	return pipeline
}

// RunPipelineStateExample demonstrates the pipeline with state management
func RunPipelineStateExample() error {
	// Create pipeline
	pipeline := NewDataPipeline()

	// Initialize data
	ctx := context.Background()
	data := map[string]any{
		"input": "test data",
	}

	// Start pipeline
	ctx, data, err := pipeline.Run(ctx, data)
	if err != nil {
		return fmt.Errorf("pipeline failed: %v", err)
	}

	// Pause the pipeline after processing data
	if pipeline.IsRunning() {
		err = pipeline.Pause()
		if err != nil {
			return fmt.Errorf("failed to pause pipeline: %v", err)
		}
		fmt.Println("Pipeline paused after processing data")
	}

	// Save pipeline state
	state := pipeline.GetState()
	stateJSON, err := state.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to save state: %v", err)
	}
	fmt.Printf("Saved pipeline state: %s\n", string(stateJSON))

	// Create a new pipeline instance
	newPipeline := NewDataPipeline()

	// Load saved state
	newState := wf.NewState()
	if err := newState.FromJSON(stateJSON); err != nil {
		return fmt.Errorf("failed to load state: %v", err)
	}
	newPipeline.SetState(newState)

	// Resume pipeline
	ctx, data, err = newPipeline.Resume(ctx, data)
	if err != nil {
		return fmt.Errorf("pipeline resume failed: %v", err)
	}

	return nil
}
