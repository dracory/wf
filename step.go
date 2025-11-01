package wf

import (
	"context"
	"errors"

	"github.com/dracory/uid"
)

type stepImplementation struct {
	id      string
	name    string
	data    map[string]any
	handler StepHandler
	state   StateInterface
}

// NewStep creates a new step with the given options
func NewStep(opts ...interface{}) StepInterface {
	step := &stepImplementation{
		id:    uid.HumanUid(),
		name:  "",
		data:  make(map[string]any),
		state: NewState(),
	}

	// Apply all options
	for _, opt := range opts {
		switch o := opt.(type) {
		case func(Nameable):
			o(step) // Handles WithName
		case func(Identifiable):
			o(step) // Handles WithID
		case func(StepInterface):
			o(step) // Handles WithHandler and other Step-specific options
		}
	}

	return step
}

func (s *stepImplementation) GetID() string {
	return s.id
}

func (s *stepImplementation) SetID(id string) {
	s.id = id
}

func (s *stepImplementation) GetName() string {
	return s.name
}

func (s *stepImplementation) SetName(name string) {
	s.name = name
}

// GetHandler returns the step's execution function
func (s *stepImplementation) GetHandler() StepHandler {
	return s.handler
}

// SetHandler sets the step's execution function
func (s *stepImplementation) SetHandler(fn StepHandler) {
	s.handler = fn
}

// Run executes the step's function with the given context
func (s *stepImplementation) Run(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
	// If we have a saved state, use it
	if s.state.GetStatus() == StateStatus(StateStatusPaused) {
		return s.resumeFromState(ctx, data)
	}

	// Initialize new state
	s.state = NewState()
	s.state.SetStatus(StateStatus(StateStatusRunning))
	s.state.SetWorkflowData(data)
	s.state.SetCurrentStepID(s.id)

	// Execute step
	ctx, data, err := s.handler(ctx, data)
	if err != nil {
		s.state.SetStatus(StateStatus(StateStatusFailed))
		return ctx, data, err
	}

	// Mark step as completed
	s.state.AddCompletedStep(s.id)
	s.state.SetWorkflowData(data)
	s.state.SetStatus(StateStatus(StateStatusComplete))

	return ctx, data, nil
}

// Pause pauses the workflow execution
func (s *stepImplementation) Pause() error {
	if s.state.GetStatus() != StateStatus(StateStatusRunning) {
		return errors.New("workflow is not running")
	}
	s.state.SetStatus(StateStatus(StateStatusPaused))
	return nil
}

// Resume resumes the workflow execution from the last saved state
func (s *stepImplementation) Resume(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
	if s.state.GetStatus() != StateStatus(StateStatusPaused) {
		return ctx, data, errors.New("workflow is not paused")
	}
	return s.resumeFromState(ctx, data)
}

// resumeFromState resumes the workflow from the saved state
func (s *stepImplementation) resumeFromState(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
	// Update data with saved state
	savedData := s.state.GetWorkflowData()
	for k, v := range savedData {
		data[k] = v
	}

	// Execute step
	s.state.SetStatus(StateStatus(StateStatusRunning))
	ctx, data, err := s.handler(ctx, data)
	if err != nil {
		s.state.SetStatus(StateStatus(StateStatusFailed))
		return ctx, data, err
	}

	// Mark step as completed
	s.state.AddCompletedStep(s.id)
	s.state.SetWorkflowData(data)
	s.state.SetStatus(StateStatus(StateStatusComplete))

	return ctx, data, nil
}

// GetState returns the current workflow state
func (s *stepImplementation) GetState() StateInterface {
	return s.state
}

// SetState sets the workflow state
func (s *stepImplementation) SetState(state StateInterface) {
	s.state = state
}

// State helper methods
func (s *stepImplementation) IsRunning() bool {
	return s.state.GetStatus() == StateStatusRunning
}

func (s *stepImplementation) IsPaused() bool {
	return s.state.GetStatus() == StateStatusPaused
}

func (s *stepImplementation) IsCompleted() bool {
	return s.state.GetStatus() == StateStatusComplete
}

func (s *stepImplementation) IsFailed() bool {
	return s.state.GetStatus() == StateStatusFailed
}

func (s *stepImplementation) IsWaiting() bool {
	return s.state.GetStatus() == "" // Initial state before running
}

// // Visualize returns a DOT graph representation of the step
// func (s *stepImplementation) Visualize() string {
// 	var color string
// 	switch s.state.GetStatus() {
// 	case StateStatusRunning:
// 		color = "yellow"
// 	case StateStatusComplete:
// 		color = "green"
// 	case StateStatusFailed:
// 		color = "red"
// 	case StateStatusPaused:
// 		color = "orange"
// 	default:
// 		color = "gray"
// 	}

// 	name := s.GetName()
// 	if name == "" {
// 		name = s.GetID()
// 	}

// 	return fmt.Sprintf(`digraph {
//     node [shape=box, style=filled, fillcolor=%s];
//     "%s" [label="%s"];
// }`, color, s.GetID(), name)
// }
