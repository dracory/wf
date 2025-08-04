# Conditional Logic Example

This example demonstrates how to implement conditional logic using DAGs and pipelines in the Dracory steps package. It shows two approaches to handling different order types (digital, physical, subscription) with varying processing requirements.

## Key Features

1. DAG Implementation:
   - Uses dependencies to maintain step ordering
   - Conditionally adds steps based on order type
   - Maintains proper execution flow through dependencies

2. Pipeline Implementation:
   - Groups related steps into logical pipelines
   - Simplifies conditional logic by adding steps to pipelines
   - Maintains proper execution order within pipelines
   - Reduces code duplication

## Order Types and Processing

### Digital Orders
- ProcessOrder → ApplyDiscount → CalculateTax
- No shipping cost
- 10% discount applied
- 20% tax applied

### Physical Orders
- ProcessOrder → ApplyDiscount → AddShipping → CalculateTax
- Fixed $5 shipping cost
- 10% discount applied
- 20% tax applied

### Subscription Orders
- ProcessOrder → ApplyDiscount → CalculateTax
- No shipping cost
- 10% discount applied
- 20% tax applied

## Running the Example

To run this example:

```bash
# Run the tests
# Tests both DAG and pipeline implementations
go test -v

# Run the main program
# Demonstrates both implementations
go run main.go
```

## Implementation Details

### DAG Implementation
- Uses direct dependencies between steps
- Requires careful dependency management
- More verbose for complex conditional logic

### Pipeline Implementation
- Groups related steps into pipelines
- Simplifies conditional logic
- More maintainable for complex scenarios
- Easier to modify step order within pipelines

## Example Output

When running `main.go`, you'll see:
1. Results from the DAG implementation
2. Results from the pipeline implementation
3. Both should produce identical results for each order type
