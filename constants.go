package wf

// StateStatus represents the current status of a workflow
type StateStatus string

const (
	// State status constants
	StateStatusRunning  = "running"
	StateStatusPaused   = "paused"
	StateStatusComplete = "complete"
	StateStatusFailed   = "failed"
)
