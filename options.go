package wf

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
