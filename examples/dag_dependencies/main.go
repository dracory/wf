package main

import (
	"context"
	"fmt"
)

func main() {
	_, data, err := NewDependenciesDag().Run(context.Background(), map[string]any{})

	if err != nil {
		fmt.Printf("Error running DAG: %v\n", err)
		return
	}

	// Access values from the context
	finalPrice := data["finalPrice"]
	stepsCompleted := data["stepsCompleted"]

	fmt.Printf("Final price: %d\n", finalPrice)
	fmt.Printf("Steps completed: %v\n", stepsCompleted)
}
