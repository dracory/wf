package wf

import (
	"testing"
	"time"
)

func TestNewState(t *testing.T) {
	state := NewState()
	if state == nil {
		t.Error("NewState() returned nil")
	}
	if state.GetStatus() != StateStatusRunning {
		t.Errorf("Expected status %v, got %v", StateStatusRunning, state.GetStatus())
	}
	if state.GetData() == nil {
		t.Error("GetData() returned nil")
	}
	if len(state.GetCompletedSteps()) != 0 {
		t.Error("Expected empty completed steps")
	}
	if state.GetLastUpdated().IsZero() {
		t.Error("Expected non-zero last updated time")
	}
}

func TestStateStatus(t *testing.T) {
	state := NewState()

	// Test setting and getting status
	state.SetStatus(StateStatusPaused)
	if state.GetStatus() != StateStatusPaused {
		t.Errorf("Expected status %v, got %v", StateStatusPaused, state.GetStatus())
	}

	// Reset state for next test
	state = NewState()
	state.SetStatus(StateStatusComplete)
	if state.GetStatus() != StateStatusComplete {
		t.Errorf("Expected status %v, got %v", StateStatusComplete, state.GetStatus())
	}

	// Reset state for next test
	state = NewState()
	state.SetStatus(StateStatusFailed)
	if state.GetStatus() != StateStatusFailed {
		t.Errorf("Expected status %v, got %v", StateStatusFailed, state.GetStatus())
	}
}

func TestStateData(t *testing.T) {
	state := NewState()

	// Test setting and getting data
	testData := map[string]any{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	state.SetData(testData)
	gotData := state.GetData()
	if len(gotData) != len(testData) {
		t.Errorf("Expected data length %d, got %d", len(testData), len(gotData))
	}
	for k, v := range testData {
		if gotData[k] != v {
			t.Errorf("Expected data[%s] = %v, got %v", k, v, gotData[k])
		}
	}

	// Test workflow data methods
	workflowData := map[string]any{
		"workflow": "data",
	}
	state.SetWorkflowData(workflowData)
	gotWorkflowData := state.GetWorkflowData()
	if len(gotWorkflowData) != len(workflowData) {
		t.Errorf("Expected workflow data length %d, got %d", len(workflowData), len(gotWorkflowData))
	}
	for k, v := range workflowData {
		if gotWorkflowData[k] != v {
			t.Errorf("Expected workflow data[%s] = %v, got %v", k, v, gotWorkflowData[k])
		}
	}
}

func TestStateCompletedSteps(t *testing.T) {
	state := NewState()

	// Test adding and getting completed steps
	step1 := "step1"
	step2 := "step2"
	step3 := "step3"

	state.AddCompletedStep(step1)
	state.AddCompletedStep(step2)
	state.AddCompletedStep(step3)

	completedSteps := state.GetCompletedSteps()
	if len(completedSteps) != 3 {
		t.Errorf("Expected 3 completed steps, got %d", len(completedSteps))
	}

	// Check if all steps are present
	steps := map[string]bool{step1: true, step2: true, step3: true}
	for _, step := range completedSteps {
		if !steps[step] {
			t.Errorf("Unexpected step in completed steps: %s", step)
		}
		delete(steps, step)
	}
	if len(steps) > 0 {
		t.Errorf("Missing steps in completed steps: %v", steps)
	}
}

func TestStateCurrentStep(t *testing.T) {
	state := NewState()

	// Test setting and getting current step ID
	stepID := "current-step"
	state.SetCurrentStepID(stepID)
	if state.GetCurrentStepID() != stepID {
		t.Errorf("Expected current step ID %s, got %s", stepID, state.GetCurrentStepID())
	}
}

func TestStateLastUpdated(t *testing.T) {
	state := NewState()
	initialTime := state.GetLastUpdated()

	// Wait a bit to ensure time difference
	time.Sleep(time.Millisecond)

	// Test setting and getting last updated time
	newTime := time.Now()
	state.SetLastUpdated(newTime)
	if !state.GetLastUpdated().Equal(newTime) {
		t.Errorf("Expected last updated time %v, got %v", newTime, state.GetLastUpdated())
	}
	if state.GetLastUpdated().Equal(initialTime) {
		t.Error("Last updated time should not equal initial time")
	}
}

func TestStateJSON(t *testing.T) {
	state := NewState()

	// Set up test data
	state.SetStatus(StateStatusRunning)
	state.SetData(map[string]any{"test": "data"})
	state.SetCurrentStepID("step1")
	state.AddCompletedStep("step1")
	state.AddCompletedStep("step2")

	// Test ToJSON
	jsonData, err := state.ToJSON()
	if err != nil {
		t.Errorf("ToJSON failed: %v", err)
	}
	if len(jsonData) == 0 {
		t.Error("ToJSON returned empty data")
	}

	// Test FromJSON
	newState := NewState()
	err = newState.FromJSON(jsonData)
	if err != nil {
		t.Errorf("FromJSON failed: %v", err)
	}

	// Verify all fields were correctly serialized and deserialized
	if state.GetStatus() != newState.GetStatus() {
		t.Errorf("Expected status %v, got %v", state.GetStatus(), newState.GetStatus())
	}
	if len(state.GetData()) != len(newState.GetData()) {
		t.Errorf("Expected data length %d, got %d", len(state.GetData()), len(newState.GetData()))
	}
	if state.GetCurrentStepID() != newState.GetCurrentStepID() {
		t.Errorf("Expected current step ID %s, got %s", state.GetCurrentStepID(), newState.GetCurrentStepID())
	}
	if len(state.GetCompletedSteps()) != len(newState.GetCompletedSteps()) {
		t.Errorf("Expected completed steps length %d, got %d", len(state.GetCompletedSteps()), len(newState.GetCompletedSteps()))
	}
}

func TestStateJSONError(t *testing.T) {
	state := NewState()

	// Test FromJSON with invalid JSON
	err := state.FromJSON([]byte("invalid json"))
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	// Test FromJSON with empty data
	err = state.FromJSON([]byte{})
	if err == nil {
		t.Error("Expected error for empty data")
	}
}

func TestStateInterface(t *testing.T) {
	// This test ensures that State implements StateInterface
	var _ StateInterface = (*State)(nil)
}

func TestStateStatusTransitions(t *testing.T) {
	// Test valid state transitions
	transitions := []struct {
		from  StateStatus
		to    StateStatus
		valid bool
	}{
		{"", StateStatusRunning, true},
		{StateStatusRunning, StateStatusPaused, true},
		{StateStatusPaused, StateStatusRunning, true},
		{StateStatusRunning, StateStatusComplete, true},
		{StateStatusRunning, StateStatusFailed, true},
		{StateStatusComplete, StateStatusRunning, false},
		{StateStatusFailed, StateStatusRunning, false},
	}

	for _, tt := range transitions {
		state := NewState()
		if tt.from != "" {
			state.SetStatus(tt.from)
		}

		// Try to transition to the target state
		state.SetStatus(tt.to)

		if tt.valid {
			if state.GetStatus() != tt.to {
				t.Errorf("transition from %s to %s should be valid, got %s", tt.from, tt.to, state.GetStatus())
			}
		} else {
			if state.GetStatus() == tt.to {
				t.Errorf("transition from %s to %s should be invalid", tt.from, tt.to)
			}
		}
	}
}
