package basic_usage

import (
	"context"
	"testing"
)

func TestBasicUsage(t *testing.T) {
	dag := NewMultipleIncrementDag()

	_, data, err := dag.Run(context.Background(), map[string]any{"value": 0})
	if err != nil {
		t.Errorf("Error running DAG: %v", err)
		return
	}

	// Verify the value
	value := data["value"].(int)
	if value != 4 {
		t.Errorf("Expected value 4, got %v", value)
	}
}
