package wf

import (
	"errors"
	"sort"
)

// visitNode performs a depth-first search to detect cycles and build the topological order
func visitNode(node RunnableInterface, graph map[RunnableInterface][]RunnableInterface, visited map[RunnableInterface]bool, tempMark map[RunnableInterface]bool, result *[]RunnableInterface) error {
	if tempMark[node] {
		return errors.New("cycle detected")
	}
	if visited[node] {
		return nil
	}

	tempMark[node] = true
	for _, dep := range graph[node] {
		if err := visitNode(dep, graph, visited, tempMark, result); err != nil {
			return err
		}
	}

	tempMark[node] = false
	visited[node] = true
	*result = append([]RunnableInterface{node}, *result...)
	return nil
}

// topologicalSort performs a topological sort on the dependency graph
// The sorting is based on the ID of the nodes, in order to make the sorting deterministic
//
// Parameters:
// - graph: The dependency graph
// Returns:
// - A slice of RunnableInterface in topological order
// - An error if a cycle is detected
func topologicalSort(graph map[RunnableInterface][]RunnableInterface) ([]RunnableInterface, error) {
	visited := make(map[RunnableInterface]bool)
	tempMark := make(map[RunnableInterface]bool)
	result := []RunnableInterface{}

	// Start with any node (since the graph is connected)
	for node := range graph {
		if err := visitNode(node, graph, visited, tempMark, &result); err != nil {
			return nil, err
		}
	}

	// First reverse the result to maintain topological order
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	// Create a map to track the number of dependencies for each node
	dependencyCount := make(map[RunnableInterface]int)
	for node := range graph {
		dependencyCount[node] = 0
	}
	for _, deps := range graph {
		for _, dep := range deps {
			dependencyCount[dep]++
		}
	}

	// Sort the result to make it deterministic while preserving dependency order
	// First sort by dependency count (ascending)
	// Then sort by name to maintain a consistent order
	// Then sort by ID to maintain a consistent order
	sort.SliceStable(result, func(i, j int) bool {
		// Get the current nodes
		current := result[i]
		compare := result[j]

		// Check if current node depends on compare node
		dependsOnCompare := false
		for _, dep := range graph[current] {
			if dep == compare {
				dependsOnCompare = true
				break
			}
		}

		// Check if compare node depends on current node
		dependsOnCurrent := false
		for _, dep := range graph[compare] {
			if dep == current {
				dependsOnCurrent = true
				break
			}
		}

		// If there's a dependency relationship, respect it
		if dependsOnCompare {
			return false // current should come after compare
		}
		if dependsOnCurrent {
			return true // current should come before compare
		}

		// If both nodes have the same ID, sort by name
		// should never happen
		if current.GetID() == compare.GetID() {
			return current.GetName() < compare.GetName()
		}

		// Otherwise, sort by ID
		return current.GetID() < compare.GetID()
	})

	return result, nil
}

// buildDependencyGraph builds a graph of runner dependencies
func buildDependencyGraph(runnables map[string]RunnableInterface, dependencies map[string][]string) map[RunnableInterface][]RunnableInterface {
	graph := make(map[RunnableInterface][]RunnableInterface)

	// Add all nodes to the graph
	for _, node := range runnables {
		graph[node] = []RunnableInterface{}
	}

	// Add all dependencies
	for dependentID, dependencyIDs := range dependencies {
		dependent, ok := runnables[dependentID]
		if !ok {
			continue
		}

		for _, dependencyID := range dependencyIDs {
			dependency, ok := runnables[dependencyID]
			if !ok {
				continue
			}

			// Add dependency
			graph[dependent] = append(graph[dependent], dependency)
		}
	}

	return graph
}
