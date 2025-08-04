package main

import (
	"context"
	"fmt"
)

func main() {
	// Create and run the DAG
	dag := NewErrorHandlingDag()

	_, data, err := dag.Run(context.Background(), map[string]any{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Access the value from the context
	value := data["value"]
	fmt.Printf("Value: %v\n", value)
}
