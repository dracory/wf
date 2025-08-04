package main

import (
	"context"
	"fmt"
)

func main() {
	// Create and run DAG
	dag, err := NewConditionalDag("digital", 100.0)
	if err != nil {
		fmt.Printf("Error creating DAG: %v\n", err)
		return
	}

	_, data, err := dag.Run(context.Background(), map[string]any{
		"orderType":     "digital",
		"totalAmount":   100.0,
		"stepsExecuted": []string{},
	})
	if err != nil {
		fmt.Printf("Error running DAG: %v\n", err)
		return
	}

	// Print results
	fmt.Printf("Using step by step execution:\n")
	fmt.Printf("Total amount: %.2f\n", data["totalAmount"].(float64))
	fmt.Println("Steps executed:", data["stepsExecuted"].([]string))

	// Create and run DAG with pipelines
	dag, err = NewConditionalDagWithPipelines("digital", 100.0)
	if err != nil {
		fmt.Printf("Error creating DAG with pipelines: %v\n", err)
		return
	}

	_, data, err = dag.Run(context.Background(), map[string]any{
		"orderType":     "digital",
		"totalAmount":   100.0,
		"stepsExecuted": []string{},
	})
	if err != nil {
		fmt.Printf("Error running DAG with pipelines: %v\n", err)
		return
	}

	// Print results
	fmt.Printf("Using pipelines:\n")
	fmt.Printf("Total amount: %.2f\n", data["totalAmount"].(float64))
	fmt.Println("Steps executed:", data["stepsExecuted"].([]string))
}
