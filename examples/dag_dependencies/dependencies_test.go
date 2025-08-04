package main

import (
	"context"
	"testing"
)

func TestDependencies(t *testing.T) {
	_, data, err := NewDependenciesDag().Run(context.Background(), map[string]any{})

	if err != nil {
		t.Errorf("Error running DAG: %v", err)
		return
	}

	// Verify the processed numbers
	numbers := data["numbers"]
	if numbers == nil {
		t.Errorf("Expected numbers to be present in data")
		return
	}

	processedNumbers := numbers.([]int)
	expectedNumbers := []int{2, 4, 6}
	if !slicesEqual(processedNumbers, expectedNumbers) {
		t.Errorf("Expected processed numbers to be %v, got %v", expectedNumbers, processedNumbers)
	}

	// Verify the sum
	sum := data["sum"]
	if sum == nil {
		t.Errorf("Expected sum to be present in data")
		return
	}

	if sum.(int) != 12 {
		t.Errorf("Expected sum to be 12, got %d", sum)
	}
}

// Helper function to compare slices
func slicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
