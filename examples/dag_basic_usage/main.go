package basic_usage

import (
	"context"
	"fmt"
)

func main() {
	dag := NewMultipleIncrementDag()

	_, data, err := dag.Run(context.Background(), map[string]any{
		"value": 0,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	value := data["value"]
	fmt.Printf("Value: %v\n", value)
}
