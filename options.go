package wf

import "context"

// Nameable is an interface for types that can have a name
type Nameable interface {
	SetName(name string)
}

// Identifiable is an interface for types that can have an ID
type Identifiable interface {
	SetID(id string)
}

// WithName is a generic option that sets the name of any type that implements Nameable
func WithName(name string) func(Nameable) {
	return func(n Nameable) {
		n.SetName(name)
	}
}

// WithID is a generic option that sets the ID of any type that implements Identifiable
func WithID(id string) func(Identifiable) {
	return func(i Identifiable) {
		i.SetID(id)
	}
}

// StepOption is a function that configures a Step
// This is a type alias for backward compatibility
// Deprecated: Use functional options directly instead
type StepOption = func(StepInterface)

// WithHandler sets the handler function for a step
func WithHandler(handler func(context.Context, map[string]any) (context.Context, map[string]any, error)) func(StepInterface) {
	return func(s StepInterface) {
		s.SetHandler(handler)
	}
}

// RunnableAdder is an interface that defines the RunnableAdd method
// RunnableAdder is an interface for types that can add runnable nodes
type RunnableAdder interface {
	RunnableAdd(node ...RunnableInterface)
}

// DependencyAdder is an interface for types that can add dependencies between nodes
type DependencyAdder interface {
	DependencyAdd(dependent RunnableInterface, dependency ...RunnableInterface)
}

// WithRunnables adds multiple runnable nodes to a Pipeline or Dag.
// Nil nodes are filtered out and not added.
func WithRunnables(nodes ...RunnableInterface) func(RunnableAdder) {
	return func(ra RunnableAdder) {
		// Filter out nil nodes
		var validNodes []RunnableInterface
		for _, node := range nodes {
			if node != nil {
				validNodes = append(validNodes, node)
			}
		}
		if len(validNodes) > 0 {
			ra.RunnableAdd(validNodes...)
		}
	}
}

// WithDependency adds a dependency between nodes in a DAG.
// The dependent node will only execute after all dependency nodes have completed successfully.
// This can be used to create a fluent API when building DAGs.
//
// Example:
//   dag := NewDag(
//       WithName("My DAG"),
//       WithRunnables(step1, step2, step3),
//       WithDependency(step2, step1),  // step2 depends on step1
//       WithDependency(step3, step2),  // step3 depends on step2
//   )
func WithDependency(dependent RunnableInterface, dependencies ...RunnableInterface) func(DependencyAdder) {
	return func(da DependencyAdder) {
		if dependent != nil && len(dependencies) > 0 {
			var validDeps []RunnableInterface
			for _, dep := range dependencies {
				if dep != nil {
					validDeps = append(validDeps, dep)
				}
			}
			if len(validDeps) > 0 {
				da.DependencyAdd(dependent, validDeps...)
			}
		}
	}
}
