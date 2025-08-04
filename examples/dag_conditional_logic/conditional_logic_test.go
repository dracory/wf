package main

import (
	"testing"
)

func TestConditionalLogic(t *testing.T) {
	// Create test cases
	testCases := []struct {
		name           string
		orderType      string
		totalAmount    float64
		expectedSteps  []string
		expectedAmount float64
	}{
		{"Digital Order", "digital", 100.0, []string{"ProcessOrder", "ApplyDiscount", "CalculateTax"}, 108.0},
		{"Physical Order", "physical", 100.0, []string{"ProcessOrder", "ApplyDiscount", "AddShipping", "CalculateTax"}, 114.0},
		{"Subscription Order", "subscription", 100.0, []string{"ProcessOrder", "ApplyDiscount", "CalculateTax"}, 108.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := RunConditionalExample(tc.orderType, tc.totalAmount)
			if err != nil {
				t.Errorf("Error running DAG: %v", err)
				return
			}

			stepsExecuted := data["stepsExecuted"].([]string)
			if !equalSlices(stepsExecuted, tc.expectedSteps) {
				t.Errorf("Expected steps %v, got %v", tc.expectedSteps, stepsExecuted)
			}

			totalAmount := data["totalAmount"].(float64)
			if totalAmount != tc.expectedAmount {
				t.Errorf("Expected total amount %.2f, got %.2f", tc.expectedAmount, totalAmount)
			}
		})
	}
}

func TestConditionalLogicWithPipelines(t *testing.T) {
	// Create test cases
	testCases := []struct {
		name           string
		orderType      string
		totalAmount    float64
		expectedSteps  []string
		expectedAmount float64
	}{
		{"Digital Order", "digital", 100.0, []string{"ProcessOrder", "ApplyDiscount", "CalculateTax"}, 108.0},
		{"Physical Order", "physical", 100.0, []string{"ProcessOrder", "ApplyDiscount", "AddShipping", "CalculateTax"}, 114.0},
		{"Subscription Order", "subscription", 100.0, []string{"ProcessOrder", "ApplyDiscount", "CalculateTax"}, 108.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := RunConditionalExampleWithPipelines(tc.orderType, tc.totalAmount)
			if err != nil {
				t.Errorf("Error running DAG: %v", err)
				return
			}

			// Verify steps executed
			stepsExecuted := data["stepsExecuted"].([]string)
			if !equalSlices(stepsExecuted, tc.expectedSteps) {
				t.Errorf("Expected steps %v, got %v", tc.expectedSteps, stepsExecuted)
			}

			// Verify final amount
			totalAmount := data["totalAmount"].(float64)
			if totalAmount != tc.expectedAmount {
				t.Errorf("Expected total amount %.2f, got %.2f", tc.expectedAmount, totalAmount)
			}
		})
	}
}

// equalSlices checks if two slices are equal
func equalSlices(a, b []string) bool {
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
