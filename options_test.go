package wf

import "testing"

func Test_WithRunnables_Pipeline(t *testing.T) {
	// Create test steps
	step1 := NewStep(WithName("step1"), WithID("step1"))
	step2 := NewStep(WithName("step2"), WithID("step2"))

	// Create pipeline with WithRunnables
	pipeline := NewPipeline(
		WithName("Test Pipeline"),
		WithRunnables(step1, step2),
	)

	// Verify pipeline was created with correct nodes
	nodes := pipeline.RunnableList()
	if len(nodes) != 2 {
		t.Fatalf("Pipeline should have 2 nodes, got %d", len(nodes))
	}
	if id := nodes[0].GetID(); id != "step1" {
		t.Errorf("First node should be step1, got %s", id)
	}
	if id := nodes[1].GetID(); id != "step2" {
		t.Errorf("Second node should be step2, got %s", id)
	}
}

func Test_WithRunnables_Dag(t *testing.T) {
	// Create test steps
	step1 := NewStep(WithName("step1"), WithID("step1"))
	step2 := NewStep(WithName("step2"), WithID("step2"))

	// Create DAG with WithRunnables
	dag := NewDag(
		WithName("Test DAG"),
		WithRunnables(step1, step2),
	)

	// Verify DAG was created with correct nodes
	nodes := dag.RunnableList()
	if len(nodes) != 2 {
		t.Fatalf("DAG should have 2 nodes, got %d", len(nodes))
	}
	// Order in DAG is not guaranteed, so we'll check by ID
	ids := make(map[string]bool)
	for _, node := range nodes {
		ids[node.GetID()] = true
	}
	if !ids["step1"] {
		t.Error("DAG should contain step1")
	}
	if !ids["step2"] {
		t.Error("DAG should contain step2")
	}
}

func Test_WithRunnables_Empty(t *testing.T) {
	// Test with empty runnables
	pipeline := NewPipeline(WithRunnables())
	if nodes := pipeline.RunnableList(); len(nodes) != 0 {
		t.Errorf("Pipeline should have no nodes, got %d", len(nodes))
	}

	dag := NewDag(WithRunnables())
	if nodes := dag.RunnableList(); len(nodes) != 0 {
		t.Errorf("DAG should have no nodes, got %d", len(nodes))
	}
}

func Test_WithRunnables_NilNodes(t *testing.T) {
	// Test with nil nodes
	pipeline := NewPipeline(WithRunnables(nil, nil))
	// Should not panic and should not add nil nodes
	if nodes := pipeline.RunnableList(); len(nodes) != 0 {
		t.Errorf("Pipeline should not add nil nodes, got %d nodes", len(nodes))
	}

	dag := NewDag(WithRunnables(nil, nil))
	if nodes := dag.RunnableList(); len(nodes) != 0 {
		t.Errorf("DAG should not add nil nodes, got %d nodes", len(nodes))
	}
}

func Test_WithRunnables_CombinedWithOtherOptions(t *testing.T) {
	// Test WithRunnables combined with other options
	step1 := NewStep(WithName("step1"), WithID("step1"))
	pipeline := NewPipeline(
		WithName("Test Pipeline"),
		WithID("test-pipeline-1"),
		WithRunnables(step1),
	)

	if name := pipeline.GetName(); name != "Test Pipeline" {
		t.Errorf("Expected pipeline name 'Test Pipeline', got '%s'", name)
	}
	if id := pipeline.GetID(); id != "test-pipeline-1" {
		t.Errorf("Expected pipeline ID 'test-pipeline-1', got '%s'", id)
	}
	nodes := pipeline.RunnableList()
	if len(nodes) != 1 {
		t.Fatalf("Pipeline should have 1 node, got %d", len(nodes))
	}
	if id := nodes[0].GetID(); id != "step1" {
		t.Errorf("Expected node ID 'step1', got '%s'", id)
	}
}
