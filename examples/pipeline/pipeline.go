package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dracory/wf"
)

func main() {
	// Create a new pipeline
	p := wf.NewPipeline()
	p.SetName("Example Pipeline")

	// Create steps
	step1 := wf.NewStep()
	step1.SetName("Step 1")
	step1.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		fmt.Println("Executing Step 1")
		return ctx, data, nil
	})

	step2 := wf.NewStep()
	step2.SetName("Step 2")
	step2.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		fmt.Println("Executing Step 2")
		return ctx, data, nil
	})

	step3 := wf.NewStep()
	step3.SetName("Step 3")
	step3.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		fmt.Println("Executing Step 3")
		return ctx, data, nil
	})

	// Add steps to pipeline
	p.RunnableAdd(step1, step2, step3)

	// Run the pipeline
	ctx := context.Background()
	data := make(map[string]any)
	ctx, data, err := p.Run(ctx, data)
	if err != nil {
		log.Fatalf("Pipeline failed: %v", err)
	}

	// Example of pausing and resuming
	err = p.Pause()
	if err != nil {
		log.Fatalf("Failed to pause pipeline: %v", err)
	}

	// Resume the pipeline
	ctx, data, err = p.Resume(ctx, data)
	if err != nil {
		log.Fatalf("Pipeline failed to resume: %v", err)
	}

	// Check pipeline state
	state := p.GetState()
	if state.GetStatus() == wf.StateStatus(wf.StateStatusComplete) {
		fmt.Println("Pipeline completed successfully")
	}
}
