package wf

import (
	"context"
	"fmt"
	"testing"
)

func Test_VisitNode(t *testing.T) {
	// Create test steps
	step1 := NewStep()
	step1.SetName("Step1")
	step1.SetID("1")
	step2 := NewStep()
	step2.SetName("Step2")
	step2.SetID("2")
	step3 := NewStep()
	step3.SetName("Step3")
	step3.SetID("3")

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

	// Test regular dependency chain
	graph := map[RunnableInterface][]RunnableInterface{
		step1: {step2},
		step2: {step3},
		step3: {},
	}

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

func Test_TopologicalSort(t *testing.T) {
	// Create test steps with unique IDs
	step1 := NewStep()
	step1.SetName("Step1")
	step1.SetID("1")
	step2 := NewStep()
	step2.SetName("Step2")
	step2.SetID("2")
	step3 := NewStep()
	step3.SetName("Step3")
	step3.SetID("3")
	step4 := NewStep()
	step4.SetName("Step4")
	step4.SetID("4")

	// Create graph with dependencies
	graph := map[RunnableInterface][]RunnableInterface{
		step1: {},
		step2: {step1},
		step3: {step2},
		step4: {step1},
	}

	// Test case 1: Regular dependency chain
	result, err := topologicalSort(graph)
	if err != nil {
		t.Errorf("topologicalSort failed: %v", err)
	}

	// Verify result order
	if len(result) != 4 {
		t.Errorf("Expected 4 nodes in result, got %d", len(result))
	}

	// Verify step1 is first since it has no dependencies and lowest ID
	if result[0] != step1 {
		t.Errorf("Expected step1 to be first, got %s", result[0].GetName())
	}

	// Verify step2 comes after step1 since it depends on step1 and has next lowest ID
	if result[1] != step2 {
		t.Errorf("Expected step2 to be second, got %s", result[1].GetName())
	}

	// Verify step3 comes after step2 since it depends on step2 and has next lowest ID
	if result[2] != step3 {
		t.Errorf("Expected step3 to be third, got %s", result[2].GetName())
	}

	// Verify step4 comes after step1 since it depends on step1 and has highest ID
	if result[3] != step4 {
		t.Errorf("Expected step4 to be fourth, got %s", result[3].GetName())
	}

	// Test case 2: Circular dependency
	graphWithCycle := map[RunnableInterface][]RunnableInterface{
		step1: {step2},
		step2: {step3},
		step3: {step1},
	}

	_, err = topologicalSort(graphWithCycle)
	if err == nil {
		t.Error("Expected cycle detection error, got nil")
	} else if err.Error() != "cycle detected" {
		t.Errorf("Expected cycle detected error, got: %v", err)
	}

	// Test case 3: Multiple independent chains
	graphWithMultipleChains := map[RunnableInterface][]RunnableInterface{
		step1: {},
		step2: {step1},
		step3: {},
		step4: {step3},
	}

	result, err = topologicalSort(graphWithMultipleChains)
	if err != nil {
		t.Errorf("topologicalSort failed: %v", err)
	}

	// Verify result order - independent chains should be sorted by ID
	if len(result) != 4 {
		t.Errorf("Expected 4 nodes in result, got %d", len(result))
	}

	// Verify step1 is first since it has no dependencies and lowest ID
	if result[0] != step1 {
		t.Errorf("Expected step1 to be first, got %s", result[0].GetName())
	}

	// Verify step2 comes after step1 since it depends on step1 and has next lowest ID
	if result[1] != step2 {
		t.Errorf("Expected step2 to be second, got %s", result[1].GetName())
	}

	// Verify step3 is next since it has no dependencies and next lowest ID
	if result[2] != step3 {
		t.Errorf("Expected step3 to be third, got %s", result[2].GetName())
	}

	// Verify step4 comes after step3 since it depends on step3 and has highest ID
	if result[3] != step4 {
		t.Errorf("Expected step4 to be fourth, got %s", result[3].GetName())
	}
}

// Sorting by ID
func Test_TopologicalSort_DuplicateNames(t *testing.T) {
	// Create test steps with duplicate names
	step1 := NewStep()
	step1.SetName("CommonName1")
	step1.SetID("1")

	step2 := NewStep()
	step2.SetName("UniqueName")
	step2.SetID("2")

	step3 := NewStep()
	step3.SetName("CommonName2")
	step3.SetID("3")

	// Create graph with dependencies
	graph := map[RunnableInterface][]RunnableInterface{
		step1: {},      // step1 has no dependencies
		step2: {step1}, // step2 depends on step1
		step3: {},      // step3 has no dependencies
	}

	// Test case: Duplicate names
	result, err := topologicalSort(graph)
	if err != nil {
		t.Errorf("topologicalSort failed: %v", err)
	}

	// Verify result order
	if len(result) != 3 {
		t.Errorf("Expected 3 nodes in result, got %d", len(result))
	}

	// Verify nodes are sorted by ID first (1, 2, 3)
	if result[0].GetID() != "1" || result[1].GetID() != "2" || result[2].GetID() != "3" {
		t.Errorf("Expected steps to be sorted by ID (1, 2, 3), got %s, %s, %s",
			result[0].GetID(), result[1].GetID(), result[2].GetID())
	}

	// Verify step2 comes after step1 due to dependency
	if result[1] != step2 {
		t.Errorf("Expected step2 to be second due to dependency on step1, got %s", result[1].GetID())
	}

	// Verify step1 and step3 are in the correct group
	if result[0] != step1 || result[2] != step3 {
		t.Errorf("Expected step1 and step3 to be in positions 0 and 2, got %s and %s", result[0].GetID(), result[2].GetID())
	}
}

func Test_TopologicalSort_MultipleValidOrderings(t *testing.T) {
	// Create test steps
	stepA := NewStep()
	stepA.SetName("StepA")
	stepA.SetID("1")
	stepB := NewStep()
	stepB.SetName("StepB")
	stepB.SetID("2")
	stepC := NewStep()
	stepC.SetName("StepC")
	stepC.SetID("3")
	stepD := NewStep()
	stepD.SetName("StepD")
	stepD.SetID("4")

	// Create graph with multiple valid orderings
	graph := map[RunnableInterface][]RunnableInterface{
		stepA: {stepC},
		stepB: {stepC},
		stepC: {stepD},
		stepD: {},
	}

	result, err := topologicalSort(graph)
	if err != nil {
		t.Errorf("topologicalSort failed: %v", err)
	}

	// Verify result length
	if len(result) != 4 {
		t.Errorf("Expected 4 nodes in result, got %d", len(result))
	}

	// Verify dependencies are respected
	for i, step := range result {
		for _, dep := range graph[step] {
			depIndex := -1
			for j, depStep := range result {
				if depStep.GetID() == dep.GetID() {
					depIndex = j
					break
				}
			}
			if depIndex == -1 {
				t.Errorf("Dependency %s not found in result", dep.GetID())
			}
			if i <= depIndex {
				t.Errorf("Step %s should come after its dependency %s", step.GetID(), dep.GetID())
			}
		}
	}
}

func Test_TopologicalSort_LargeComplexGraph(t *testing.T) {
	// Create a larger set of steps
	steps := make(map[string]RunnableInterface)
	for i := 0; i < 10; i++ {
		step := NewStep()
		step.SetName(fmt.Sprintf("Step%d", i))
		step.SetID(fmt.Sprintf("%d", i))
		steps[step.GetName()] = step
	}

	// Create a complex dependency graph
	graph := map[RunnableInterface][]RunnableInterface{
		steps["Step0"]: {steps["Step1"], steps["Step2"]},
		steps["Step1"]: {steps["Step3"]},
		steps["Step2"]: {steps["Step4"], steps["Step5"]},
		steps["Step3"]: {steps["Step6"]},
		steps["Step4"]: {steps["Step6"]},
		steps["Step5"]: {steps["Step7"]},
		steps["Step6"]: {steps["Step8"]},
		steps["Step7"]: {steps["Step8"]},
		steps["Step8"]: {steps["Step9"]},
		steps["Step9"]: {},
	}

	result, err := topologicalSort(graph)
	if err != nil {
		t.Errorf("topologicalSort failed: %v", err)
	}

	// Verify result length
	if len(result) != 10 {
		t.Errorf("Expected 10 nodes in result, got %d", len(result))
	}

	// Verify dependencies are respected
	for i, step := range result {
		for _, dep := range graph[step] {
			depIndex := -1
			for j, depStep := range result {
				if depStep.GetID() == dep.GetID() {
					depIndex = j
					break
				}
			}
			if depIndex == -1 {
				t.Errorf("Dependency %s not found in result", dep.GetID())
			}
			if i <= depIndex {
				t.Errorf("Step %s should come after its dependency %s", step.GetID(), dep.GetID())
			}
		}
	}
}

// Test_TopologicalSort_DuplicateNamesWithDependencies verifies that the topological sort:
//  1. Sorts steps by ID in ascending order (1, 2, 3, 4)
//  2. Maintains dependency order despite ID sorting
//  3. Ensures that steps with dependencies appear after their dependencies
//     even though they might have lower IDs
func Test_TopologicalSort_DuplicateNamesWithDependencies(t *testing.T) {
	// Create test steps with duplicate names and dependencies
	step1 := NewStep()
	step1.SetName("CommonName1")
	step1.SetID("1")

	step2 := NewStep()
	step2.SetName("UniqueName")
	step2.SetID("2")

	step3 := NewStep()
	step3.SetName("CommonName2")
	step3.SetID("3")

	step4 := NewStep()
	step4.SetName("OtherName")
	step4.SetID("4")

	// Create graph with dependencies
	graph := map[RunnableInterface][]RunnableInterface{
		step4: {step1, step3}, // step4 depends on step1 and step3
		step1: {step2},        // step1 depends on step2
		step2: {},             // step2 has no dependencies
		step3: {},             // step3 has no dependencies
	}

	// Test case: Duplicate names
	result, err := topologicalSort(graph)
	if err != nil {
		t.Errorf("topologicalSort failed: %v", err)
	}

	// Verify result length
	if len(result) != 4 {
		t.Errorf("Expected 4 nodes in result, got %d", len(result))
	}

	// Verify nodes are sorted by ID first (1, 2, 3, 4)
	if result[0].GetID() != "2" || result[1].GetID() != "1" || result[2].GetID() != "3" || result[3].GetID() != "4" {
		t.Errorf("Expected steps to be sorted by dependency order (2, 1, 3, 4), got %s, %s, %s, %s",
			result[0].GetID(), result[1].GetID(), result[2].GetID(), result[3].GetID())
	}

	// Verify dependencies are respected despite ID sorting
	// Step2 depends on step1, so step1 must appear before step2
	// Step1 and step3 depend on step4, so step4 must appear before both
	for i, step := range result {
		for _, dep := range graph[step] {
			// Find the position of the dependency by ID
			depIndex := -1
			for j := 0; j < len(result); j++ {
				if result[j].GetID() == dep.GetID() {
					depIndex = j
					break
				}
			}
			if depIndex == -1 {
				t.Errorf("Dependency %s not found in result", dep.GetID())
			}
			// Ensure the dependency appears before the step
			if i < depIndex {
				t.Errorf("Step %s should come after its dependency %s", step.GetID(), dep.GetID())
			}
		}
	}
}

func Test_BuildDependencyGraph_BasicChain(t *testing.T) {
	// Create test steps
	step1 := NewStep()
	step1.SetName("Step1")
	step1.SetID("1")
	step2 := NewStep()
	step2.SetName("Step2")
	step2.SetID("2")
	step3 := NewStep()
	step3.SetName("Step3")
	step3.SetID("3")

	// Create runnables map
	runnables := map[string]RunnableInterface{
		"1": step1,
		"2": step2,
		"3": step3,
	}

	// Set up dependencies
	dependencies := map[string][]string{
		"2": {"1"}, // Step2 depends on Step1
		"3": {"2"}, // Step3 depends on Step2
	}

	// Build dependency graph
	graph := buildDependencyGraph(runnables, dependencies)

	// Verify graph structure
	if len(graph[step1]) != 0 {
		t.Errorf("Step1 should have no dependencies")
	}
	if len(graph[step2]) != 1 || graph[step2][0] != step1 {
		t.Errorf("Step2 should depend only on Step1")
	}
	if len(graph[step3]) != 1 || graph[step3][0] != step2 {
		t.Errorf("Step3 should depend only on Step2")
	}
}

func Test_BuildDependencyGraph_CircularDependencies(t *testing.T) {
	// Create test steps
	step1 := NewStep()
	step1.SetName("Step1")
	step1.SetID("1")
	step2 := NewStep()
	step2.SetName("Step2")
	step2.SetID("2")

	// Create runnables map
	runnables := map[string]RunnableInterface{
		"1": step1,
		"2": step2,
	}

	// Set up dependencies
	dependencies := map[string][]string{
		"1": {"2"}, // Step1 depends on Step2
		"2": {"1"}, // Step2 depends on Step1
	}

	// Build dependency graph
	graph := buildDependencyGraph(runnables, dependencies)

	// Verify graph structure
	if len(graph) != 2 {
		t.Errorf("Expected 2 nodes in graph, got %d", len(graph))
	}

	// Verify circular dependencies
	if len(graph[step1]) != 1 || graph[step1][0] != step2 {
		t.Errorf("Expected step1 to depend on step2")
	}

	if len(graph[step2]) != 1 || graph[step2][0] != step1 {
		t.Errorf("Expected step2 to depend on step1")
	}
}
