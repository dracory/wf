package wf

import (
	"context"
	"errors"
	"testing"
)

func Test_Dag_Basic(t *testing.T) {
	// Create a simple DAG with two steps
	dag := NewDag()
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

	// Add steps to DAG
	dag.RunnableAdd(step1, step2)

	// Add dependency
	dag.DependencyAdd(step2, step1)

	// Test basic DAG structure
	if len(dag.RunnableList()) != 2 {
		t.Errorf("Expected 2 runnables, got %d", len(dag.RunnableList()))
	}

	// Test dependency
	deps := dag.DependencyList(context.Background(), step2, make(map[string]any))
	if len(deps) != 1 {
		t.Errorf("Expected 1 dependency, got %d", len(deps))
	}
}

func Test_Dag_Remove(t *testing.T) {
	// Create a DAG with steps
	dag := NewDag()
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

	// Add steps to DAG
	dag.RunnableAdd(step1, step2)

	// Add dependency
	dag.DependencyAdd(step2, step1)

	// Remove step1
	removed := dag.RunnableRemove(step1)
	if !removed {
		t.Errorf("Failed to remove step1")
	}

	// Verify step1 is removed
	if len(dag.RunnableList()) != 1 {
		t.Errorf("Expected 1 runnable after removal, got %d", len(dag.RunnableList()))
	}

	// Verify dependencies are cleaned up
	deps := dag.DependencyList(context.Background(), step2, make(map[string]any))
	if len(deps) != 0 {
		t.Errorf("Expected 0 dependencies after removal, got %d", len(deps))
	}
}

func Test_Dag_TopologicalSort(t *testing.T) {
	// Create a DAG with a cycle
	dag := NewDag()
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

	// Add steps to DAG
	dag.RunnableAdd(step1, step2)

	// Create a cycle
	dag.DependencyAdd(step1, step2)
	dag.DependencyAdd(step2, step1)

	// Test Run with cycle
	ctx := context.Background()
	_, _, err := dag.Run(ctx, make(map[string]any))
	if err == nil {
		t.Errorf("Expected error for cycle, got nil")
	}
}

func Test_Dag_VisitNode(t *testing.T) {
	// Create a simple DAG with nodes
	dag := NewDag()
	step1 := NewStep()
	step2 := NewStep()
	step3 := NewStep()
	step1.SetName("Step1")
	step2.SetName("Step2")
	step3.SetName("Step3")

	// Set handlers for steps
	step1.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		return ctx, data, nil
	})
	step2.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		return ctx, data, nil
	})
	step3.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		return ctx, data, nil
	})

	// Add steps to DAG
	dag.RunnableAdd(step1, step2, step3)

	// Create a graph with dependencies
	graph := map[RunnableInterface][]RunnableInterface{
		step1: {step2},
		step2: {step3},
		step3: {},
	}

	// Test regular dependency chain
	visited := make(map[RunnableInterface]bool)
	tempMark := make(map[RunnableInterface]bool)
	result := []RunnableInterface{}

	if err := visitNode(step1, graph, visited, tempMark, &result); err != nil {
		t.Errorf("visitNode failed for regular dependency chain: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 nodes in result, got %d", len(result))
	}

	// Test cycle detection
	graphWithCycle := map[RunnableInterface][]RunnableInterface{
		step1: {step2},
		step2: {step3},
		step3: {step1},
	}

	visited = make(map[RunnableInterface]bool)
	tempMark = make(map[RunnableInterface]bool)
	result = []RunnableInterface{}

	if err := visitNode(step1, graphWithCycle, visited, tempMark, &result); err == nil {
		t.Error("Expected cycle detection error, got nil")
	} else if err.Error() != "cycle detected" {
		t.Errorf("Expected cycle detected error, got: %v", err)
	}

	// Test conditional dependencies
	graphWithConditional := map[RunnableInterface][]RunnableInterface{
		step1: {step2},
		step2: {},
		step3: {},
	}

	visited = make(map[RunnableInterface]bool)
	tempMark = make(map[RunnableInterface]bool)
	result = []RunnableInterface{}

	if err := visitNode(step1, graphWithConditional, visited, tempMark, &result); err != nil {
		t.Errorf("visitNode failed for conditional dependencies: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 nodes in result, got %d", len(result))
	}
}

func Test_Dag_Run_SuccessfulExecution(t *testing.T) {
	step1 := NewStep()
	step1.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		data["value"] = 1
		return ctx, data, nil
	})

	step2 := NewStep()
	step2.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		if _, ok := data["value"].(int); !ok {
			return ctx, data, errors.New("value not found")
		}
		value := data["value"].(int)
		data["value"] = value * 2
		return ctx, data, nil
	})

	dag := NewDag()
	dag.RunnableAdd(step1, step2)
	dag.DependencyAdd(step2, step1)

	ctx := context.Background()
	data := make(map[string]any)
	data["value"] = 0 // Initialize the value
	_, data, err := dag.Run(ctx, data)
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}

	if data["value"].(int) != 2 {
		t.Errorf("Expected value to be 2, got: %v", data["value"])
	}
}

func Test_Dag_Run_ErrorPropagation(t *testing.T) {
	step1 := NewStep()
	step1.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		return ctx, data, errors.New("step1 failed")
	})

	step2 := NewStep()
	step2.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		if _, ok := data["value"].(int); !ok {
			return ctx, data, errors.New("value not found")
		}
		value := data["value"].(int)
		data["value"] = value * 2
		return ctx, data, nil
	})

	dag := NewDag()
	dag.RunnableAdd(step1, step2)
	dag.DependencyAdd(step2, step1)

	ctx := context.Background()
	data := make(map[string]any)
	_, _, err := dag.Run(ctx, data)
	if err == nil {
		t.Error("Expected error from step1, got nil")
	} else if err.Error() != "step1 failed" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func Test_Dag_Run_CircularDependency(t *testing.T) {
	dag := NewDag()
	step1 := NewStep()
	step2 := NewStep()
	dag.RunnableAdd(step1, step2)
	dag.DependencyAdd(step1, step2)
	dag.DependencyAdd(step2, step1)

	ctx := context.Background()
	data := make(map[string]any)
	_, _, err := dag.Run(ctx, data)
	if err == nil {
		t.Errorf("Expected error for circular dependency, got nil")
	}
}

func Test_Dag_Run_StepDependencies(t *testing.T) {
	dag := NewDag()
	step1 := NewStep()
	step2 := NewStep()

	// Set handlers for steps
	step1.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		data["value"] = 1
		return ctx, data, nil
	})

	step2.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		if value, ok := data["value"].(int); ok {
			data["value"] = value * 2
		}
		return ctx, data, nil
	})

	dag.RunnableAdd(step1, step2)
	dag.DependencyAdd(step2, step1)

	ctx := context.Background()
	data := make(map[string]any)
	data["value"] = 0 // Initialize the value
	_, data, err := dag.Run(ctx, data)
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}

	if data["value"].(int) != 2 {
		t.Errorf("Expected value to be 2, got: %v", data["value"])
	}
}

func Test_Dag_Run_ParallelStepExecution(t *testing.T) {
	dag := NewDag()
	step1 := NewStep()
	step2 := NewStep()

	// Set handlers for steps
	step1.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		data["value"] = 1
		return ctx, data, nil
	})

	step2.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		if value, ok := data["value"].(int); ok {
			data["value"] = value * 2
		}
		return ctx, data, nil
	})

	dag.RunnableAdd(step1, step2)
	dag.DependencyAdd(step2, step1)

	ctx := context.Background()
	data := make(map[string]any)
	data["value"] = 0 // Initialize the value
	_, data, err := dag.Run(ctx, data)
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}

	if data["value"].(int) != 2 {
		t.Errorf("Expected value to be 2, got: %v", data["value"])
	}
}

func Test_Dag_Run_DependencyResolution(t *testing.T) {
	dag := NewDag()
	step1 := NewStep()
	step2 := NewStep()

	// Set handlers for steps
	step1.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		data["value"] = 1
		return ctx, data, nil
	})

	step2.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		if value, ok := data["value"].(int); ok {
			data["value"] = value * 2
		}
		return ctx, data, nil
	})

	dag.RunnableAdd(step1, step2)
	dag.DependencyAdd(step2, step1)

	ctx := context.Background()
	data := make(map[string]any)
	data["value"] = 0 // Initialize the value
	_, data, err := dag.Run(ctx, data)
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}

	if data["value"].(int) != 2 {
		t.Errorf("Expected value to be 2, got: %v", data["value"])
	}
}
