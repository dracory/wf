package wf

import (
	"encoding/json"
	"slices"
	"time"
)

// StateInterface defines the interface for workflow state management
type StateInterface interface {
	GetStatus() StateStatus
	SetStatus(status StateStatus)

	GetData() map[string]any
	SetData(data map[string]any)

	ToJSON() ([]byte, error)
	FromJSON(data []byte) error

	GetCurrentStepID() string
	SetCurrentStepID(id string)

	GetCompletedSteps() []string
	AddCompletedStep(id string)

	GetWorkflowData() map[string]any
	SetWorkflowData(data map[string]any)

	GetLastUpdated() time.Time
	SetLastUpdated(t time.Time)
}

// State represents the current state of a workflow
type State struct {
	Status         StateStatus
	Data           map[string]any
	CurrentStepID  string
	CompletedSteps []string
	LastUpdated    time.Time
}

// NewState creates a new workflow state
func NewState() StateInterface {
	return &State{
		Status:         StateStatusRunning,
		Data:           make(map[string]any),
		CompletedSteps: make([]string, 0),
		LastUpdated:    time.Now(),
	}
}

// GetStatus returns the current status of the workflow
func (s *State) GetStatus() StateStatus {
	return s.Status
}

// SetStatus sets the current status of the workflow
func (s *State) SetStatus(status StateStatus) {
	// Define valid state transitions
	validTransitions := map[StateStatus][]StateStatus{
		"":                  {StateStatusRunning},
		StateStatusRunning:  {StateStatusPaused, StateStatusComplete, StateStatusFailed},
		StateStatusPaused:   {StateStatusRunning},
		StateStatusComplete: {}, // No valid transitions from complete
		StateStatusFailed:   {}, // No valid transitions from failed
	}

	// Check if the transition is valid
	currentStatus := s.Status
	validNextStates, exists := validTransitions[currentStatus]

	if !exists {
		// If current status is not in the map, allow any transition (for backward compatibility)
		s.Status = status
		s.LastUpdated = time.Now()
		return
	}

	// If there are no valid next states, don't allow any transition
	if len(validNextStates) == 0 {
		return
	}

	// Check if the requested transition is valid
	if slices.Contains(validNextStates, status) {
		s.Status = status
		s.LastUpdated = time.Now()
		return
	}

	// If we get here, the transition is not valid, so we don't change the state
}

// GetData returns the current data of the workflow
func (s *State) GetData() map[string]any {
	return s.Data
}

// SetData sets the current data of the workflow
func (s *State) SetData(data map[string]any) {
	s.Data = data
	s.LastUpdated = time.Now()
}

// ToJSON converts the state to JSON
func (s *State) ToJSON() ([]byte, error) {
	s.LastUpdated = time.Now()
	return json.Marshal(s)
}

// FromJSON loads the state from JSON
func (s *State) FromJSON(data []byte) error {
	return json.Unmarshal(data, s)
}

// GetCurrentStepID returns the ID of the current step
func (s *State) GetCurrentStepID() string {
	return s.CurrentStepID
}

// SetCurrentStepID sets the ID of the current step
func (s *State) SetCurrentStepID(id string) {
	s.CurrentStepID = id
	s.LastUpdated = time.Now()
}

// GetCompletedSteps returns the list of completed step IDs
func (s *State) GetCompletedSteps() []string {
	return s.CompletedSteps
}

// AddCompletedStep adds a step ID to the completed steps list
func (s *State) AddCompletedStep(id string) {
	s.CompletedSteps = append(s.CompletedSteps, id)
	s.LastUpdated = time.Now()
}

// GetWorkflowData returns the workflow data
func (s *State) GetWorkflowData() map[string]any {
	return s.Data
}

// SetWorkflowData sets the workflow data
func (s *State) SetWorkflowData(data map[string]any) {
	s.Data = data
	s.LastUpdated = time.Now()
}

// GetLastUpdated returns the timestamp of the last update
func (s *State) GetLastUpdated() time.Time {
	return s.LastUpdated
}

// SetLastUpdated sets the timestamp of the last update
func (s *State) SetLastUpdated(t time.Time) {
	s.LastUpdated = t
}
