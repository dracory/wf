# Dependencies Example

This example demonstrates how to create steps with complex dependencies and ensure they execute in the correct order, using a data processing scenario.

## Overview

The example shows:
1. How to create multiple steps with different operations
2. How to set up complex dependencies between steps
3. How the steps are automatically sorted and executed in dependency order
4. How to maintain proper execution flow through dependencies

## Key Concepts

- Step Dependencies: Using `DependencyAdd` to specify that one step depends on another
- Topological Sorting: The steps are automatically sorted to ensure dependencies are respected
- Execution Order: Steps are executed in the order determined by their dependencies
- Complex Dependencies: Steps can depend on multiple other steps

## Data Processing Flow

The example demonstrates a simple data processing pipeline:

1. `Set Initial Data`: Initializes an array of numbers [1, 2, 3]
2. `Process Data`: Multiplies each number by 2 (result: [2, 4, 6])
3. `Calculate Sum`: Calculates the sum of the processed numbers (result: 12)

## Running the Example

To run this example:

```bash
# Run the main program
go run main.go

# Run the tests
go test -v
```

## Expected Output

The program will output:
```
Final sum: 12
Steps completed: [Set Initial Data Process Data Calculate Sum]
```

## Data Processing Process

The steps perform the following operations in order:
1. `Set Initial Data`: Creates an array with numbers [1, 2, 3]
2. `Process Data`: Multiplies each number by 2 (result: [2, 4, 6])
3. `Calculate Sum`: Adds up all numbers (result: 12)

The final sum is calculated as:
1. Start with initial numbers: [1, 2, 3]
2. Double each number: [2, 4, 6]
3. Calculate sum: 2 + 4 + 6 = 12

## Implementation Details

### Step Dependencies
- `Process Data` depends on `Set Initial Data`
- `Calculate Sum` depends on `Process Data`

### DAG Structure
The DAG is created with three steps:
1. `Set Initial Data` - Creates the initial data
2. `Process Data` - Processes the data
3. `Calculate Sum` - Calculates the final result

The dependencies ensure that each step only runs after its dependencies have completed.
