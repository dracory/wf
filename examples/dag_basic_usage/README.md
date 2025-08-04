# Basic Usage Example

This example demonstrates the basic usage of the steps package by creating and running a DAG (Directed Acyclic Graph) with multiple steps in sequence.

## Overview

The example shows:
1. How to create a step using the `Step` function
2. How to create a DAG (Directed Acyclic Graph) using `NewDag()`
3. How to add multiple steps to the DAG
4. How to run the DAG with context and data

## Key Concepts

- Step Creation: Using the `Step` function to create steps that perform operations
- DAG Structure: Using `NewDag()` to create a directed acyclic graph of steps
- Step Addition: Using `RunnableAdd()` to add steps to the DAG
- Execution: Running the DAG with context and data map

## Example Code

Here's a quick overview of the example code:

```go
// Create a new DAG
func NewMultipleIncrementDag() steps.DagInterface {
    dag := steps.NewDag()
    dag.SetName("Multiple Increment DAG")
    
    // Add 4 increment steps
    for i := 0; i < 4; i++ {
        dag.RunnableAdd(NewIncrementStep())
    }
    
    return dag
}
```

## Running the Example

To run this example:

```bash
# Run the main program
# This will execute a DAG with 4 increment steps
# The initial value is 0 and it will be incremented 4 times
# The final value should be 4

go run main.go

# Run the tests
# The test verifies that the DAG correctly increments the value 4 times
go test -v
```

## Expected Output

The program will output:
```
Value: 4
