package wf

// Nameable is an interface for types that can have a name
type Nameable interface {
	SetName(name string)
}

// WithName is a generic option that sets the name of any type that implements Nameable
func WithName(name string) func(Nameable) {
	return func(n Nameable) {
		n.SetName(name)
	}
}
