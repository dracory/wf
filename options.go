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
