# WF Examples

This directory contains example code demonstrating various features of the wf package. Each example is in its own folder with a main program and tests.

## Examples

### Basic Usage
- [examples/dag_basic_usage/main.go](dag_basic_usage/main.go): Shows how to create and run a simple step
- [examples/dag_basic_usage/basic_usage_test.go](dag_basic_usage/basic_usage_test.go): Tests for the basic usage example

### Dependencies
- [examples/dag_dependencies/main.go](dag_dependencies/main.go): Demonstrates how to create steps with dependencies
- [examples/dag_dependencies/dependencies_test.go](dag_dependencies/dependencies_test.go): Tests for the dependencies example

### Error Handling
- [examples/dag_error_handling/main.go](dag_error_handling/main.go): Shows how errors are propagated through the step chain
- [examples/dag_error_handling/error_handling_test.go](dag_error_handling/error_handling_test.go): Tests for the error handling example

### Conditional Logic
- [examples/dag_conditional_logic/main.go](dag_conditional_logic/main.go): Demonstrates conditional logic using DAGs and pipelines
- [examples/dag_conditional_logic/conditional_logic_test.go](dag_conditional_logic/conditional_logic_test.go): Tests for the conditional logic example

## Running Examples

To run an example, navigate to its directory and use:

```bash
# Run the main program
go run main.go

# Run the tests
go test -v
