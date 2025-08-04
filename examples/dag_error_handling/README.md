# Error Handling Example

This example demonstrates error handling in a DAG using the steps package. It shows how errors can be caught, handled gracefully, and how conditional dependencies can be used to control step execution.

## Overview

The example consists of three steps:

1. `Set Initial Value`: Sets an initial value of 1
2. `Process Data`: Multiplies the value by 2
3. `Intentional Error`: Intentionally fails with an error

The steps are arranged in a DAG with dependencies:

```
Set Initial Value -> Process Data -> Intentional Error
```

## Key Concepts

- **Error Handling**: Shows how errors are propagated through the DAG
- **Step Dependencies**: Demonstrates how steps depend on each other
- **Context Management**: Shows how data is passed between steps
- **Graceful Failure**: Demonstrates how the DAG can fail gracefully while still completing successful steps

## Implementation Details

### Step Implementation

1. **Set Initial Value**
   - Sets an initial value of 1
   - Stores the value in the context

2. **Process Data**
   - Retrieves the value from the context
   - Multiplies the value by 2
   - Returns an error if the value is not found

3. **Intentional Error**
   - Always returns an error
   - Used to demonstrate error handling

### DAG Structure

The DAG is created with three steps:

1. `Set Initial Value` - Creates the initial value
2. `Process Data` - Processes the value
3. `Intentional Error` - Intentionally fails

The dependencies ensure that each step only runs after its dependencies have completed.

## Running the Example

To run the example:

```bash
# Run the main program
go run main.go

# Run the tests
go test -v
```

## Expected Output

The program will output:
```
Error: intentional error
```

The error occurs in the `Intentional Error` step, but the previous steps (`Set Initial Value` and `Process Data`) still complete successfully. This demonstrates how the DAG can handle errors gracefully while still allowing successful steps to complete their work.

## Error Handling Process

The error handling process works as follows:

1. `Set Initial Value` sets the initial value to 1
2. `Process Data` multiplies the value by 2 (result: 2)
3. `Intentional Error` fails with an error

The final value is processed successfully before the error occurs.

## Best Practices

This example demonstrates several best practices for error handling in DAGs:

1. Use proper error handling in each step
2. Handle errors at the step level
3. Use proper type assertions for data handling
4. Maintain clear step dependencies
5. Keep error messages descriptive and specific

## Testing

To run the tests:

```bash
go test -v
```

The tests verify that:
1. The error is properly propagated
2. The value is correctly processed by the successful steps
3. The error message matches expectations
4. The DAG fails gracefully
