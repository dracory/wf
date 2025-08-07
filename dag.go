package wf

import (
	"context"
	"errors"
	"slices"

	"github.com/gouniverse/uid"
)

type Dag struct {
	// id of the dag
	id string

	// name of the dag
	name string

	// runnable sequence (IDs)
	runnableSequence []string

	// runnables (ID, RunnableInterface)
	runnables map[string]RunnableInterface

	// dependencies (DependentID, DependencyIDs []string)
	dependencies map[string][]string

	// current state of the workflow
	state StateInterface
}

// NewDag creates a new DAG with the given options
func NewDag(opts ...func(Nameable)) DagInterface {
	dag := &Dag{
		id:               uid.HumanUid(),
		name:             "New DAG",
		runnableSequence: make([]string, 0),
		runnables:        make(map[string]RunnableInterface),
		dependencies:     make(map[string][]string),
		state:            NewState(),
	}

	// Apply all options (WithName, etc.)
	for _, opt := range opts {
		opt(dag)
	}

	return dag
}

func (d *Dag) GetID() string {
	return d.id
}

func (d *Dag) SetID(id string) {
	d.id = id
}

func (d *Dag) GetName() string {
	return d.name
}

func (d *Dag) SetName(name string) {
	d.name = name
}

// RunnableAdd adds a single node to the DAG.
func (d *Dag) RunnableAdd(node ...RunnableInterface) {
	for _, n := range node {
		if n == nil {
			continue
		}
		id := n.GetID()
		if id == "" {
			id = uid.HumanUid()
			n.SetID(id)
		}

		// Check for duplicate ID
		if _, exists := d.runnables[id]; exists {
			// Generate a new ID if there's a conflict
			newID := uid.HumanUid()
			n.SetID(newID)
			id = newID
		}

		d.runnables[id] = n
		if !slices.Contains(d.runnableSequence, id) {
			d.runnableSequence = append(d.runnableSequence, id)
		}
	}
}

// RunnableRemove removes a node from the DAG.
func (d *Dag) RunnableRemove(node RunnableInterface) bool {
	id := node.GetID()
	if id == "" {
		return false
	}

	if _, exists := d.runnables[id]; !exists {
		return false
	}

	// Remove from runnables
	delete(d.runnables, id)

	// Remove from runnableSequence
	for i, seqID := range d.runnableSequence {
		if seqID == id {
			d.runnableSequence = append(d.runnableSequence[:i], d.runnableSequence[i+1:]...)
			break
		}
	}

	// Remove dependencies
	delete(d.dependencies, id)

	// Remove this node from other nodes' dependencies
	for depID, depList := range d.dependencies {
		for i, dep := range depList {
			if dep == id {
				d.dependencies[depID] = append(depList[:i], depList[i+1:]...)
				break
			}
		}
	}

	return true
}

// RunnableList returns all runnable nodes in the DAG.
func (d *Dag) RunnableList() []RunnableInterface {
	result := make([]RunnableInterface, 0, len(d.runnables))
	for _, node := range d.runnables {
		result = append(result, node)
	}
	return result
}

// Run executes all nodes in the DAG in the correct order
func (d *Dag) Run(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
	// If we have a saved state, use it
	if d.state.GetStatus() == StateStatus(StateStatusPaused) {
		return d.resumeFromState(ctx, data)
	}

	// Initialize new state
	d.state = NewState()
	d.state.SetStatus(StateStatus(StateStatusRunning))
	d.state.SetWorkflowData(data)

	// Build dependency graph
	graph := buildDependencyGraph(d.runnables, d.dependencies)

	// Get execution order
	order, err := topologicalSort(graph)
	if err != nil {
		d.state.SetStatus(StateStatus(StateStatusFailed))
		return ctx, data, err
	}

	// Execute steps in order
	for _, node := range order {
		// Skip completed steps
		if slices.Contains(d.state.GetCompletedSteps(), node.GetID()) {
			continue
		}

		// Update current step
		d.state.SetCurrentStepID(node.GetID())

		// Execute step
		ctx, data, err = node.Run(ctx, data)
		if err != nil {
			d.state.SetStatus(StateStatus(StateStatusFailed))
			return ctx, data, err
		}

		// Mark step as completed
		d.state.AddCompletedStep(node.GetID())
		d.state.SetWorkflowData(data)
	}

	d.state.SetStatus(StateStatus(StateStatusComplete))
	return ctx, data, nil
}

// Pause pauses the workflow execution
func (d *Dag) Pause() error {
	if d.state.GetStatus() != StateStatus(StateStatusRunning) {
		return errors.New("workflow is not running")
	}
	d.state.SetStatus(StateStatus(StateStatusPaused))
	return nil
}

// Resume resumes the workflow execution from the last saved state
func (d *Dag) Resume(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
	if d.state.GetStatus() != StateStatus(StateStatusPaused) {
		return ctx, data, errors.New("workflow is not paused")
	}
	return d.resumeFromState(ctx, data)
}

// resumeFromState resumes the workflow from the saved state
func (d *Dag) resumeFromState(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
	// Update data with saved state
	savedData := d.state.GetWorkflowData()
	for k, v := range savedData {
		data[k] = v
	}

	// Build dependency graph
	graph := buildDependencyGraph(d.runnables, d.dependencies)

	// Get execution order
	order, err := topologicalSort(graph)
	if err != nil {
		d.state.SetStatus(StateStatus(StateStatusFailed))
		return ctx, data, err
	}

	// Find the current step
	currentStepID := d.state.GetCurrentStepID()
	var currentStepIndex int
	for i, node := range order {
		if node.GetID() == currentStepID {
			currentStepIndex = i
			break
		}
	}

	// Execute remaining steps
	d.state.SetStatus(StateStatus(StateStatusRunning))
	for i := currentStepIndex; i < len(order); i++ {
		node := order[i]

		// Skip completed steps
		if slices.Contains(d.state.GetCompletedSteps(), node.GetID()) {
			continue
		}

		// Update current step
		d.state.SetCurrentStepID(node.GetID())

		// Execute step
		ctx, data, err = node.Run(ctx, data)
		if err != nil {
			d.state.SetStatus(StateStatus(StateStatusFailed))
			return ctx, data, err
		}

		// Mark step as completed
		d.state.AddCompletedStep(node.GetID())
		d.state.SetWorkflowData(data)
	}

	d.state.SetStatus(StateStatus(StateStatusComplete))
	return ctx, data, nil
}

// GetState returns the current workflow state
func (d *Dag) GetState() StateInterface {
	return d.state
}

// SetState sets the workflow state
func (d *Dag) SetState(state StateInterface) {
	d.state = state
}

// State helper methods
func (d *Dag) IsRunning() bool {
	return d.state.GetStatus() == StateStatusRunning
}

func (d *Dag) IsPaused() bool {
	return d.state.GetStatus() == StateStatusPaused
}

func (d *Dag) IsCompleted() bool {
	return d.state.GetStatus() == StateStatusComplete
}

func (d *Dag) IsFailed() bool {
	return d.state.GetStatus() == StateStatusFailed
}

func (d *Dag) IsWaiting() bool {
	return d.state.GetStatus() == "" // Initial state before running
}

// DependencyAdd adds a dependency between two nodes.
func (d *Dag) DependencyAdd(dependent RunnableInterface, dependency ...RunnableInterface) {
	dependentID := dependent.GetID()
	for _, dep := range dependency {
		depID := dep.GetID()
		d.dependencies[dependentID] = append(d.dependencies[dependentID], depID)
	}
}

// DependencyList returns all dependencies for a given node.
func (d *Dag) DependencyList(ctx context.Context, node RunnableInterface, data map[string]any) []RunnableInterface {
	dependencies := []RunnableInterface{}

	// Get all direct dependencies
	dependentID := node.GetID()
	if deps, ok := d.dependencies[dependentID]; ok {
		for _, depID := range deps {
			dep, ok := d.runnables[depID]
			if !ok {
				continue
			}

			// Add regular dependency
			dependencies = append(dependencies, dep)
		}
	}

	return dependencies
}
