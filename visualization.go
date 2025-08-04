package wf

import (
	"fmt"
	"slices"
	"strings"
)

// --- Constants (Color, Node Style, Edge Style) remain the same ---

// Color constants for DOT graph visualization
const (
	colorWhite  = "#ffffff" // Default node fill
	colorRed    = "#F44336" // Failed status
	colorYellow = "#FFC107" // Paused status
	colorBlue   = "#2196F3" // Running status
	colorGreen  = "#4CAF50" // Completed status/edges
	colorGrey   = "#9E9E9E" // Default edge, fallback fill
)

// Node style constants
const (
	nodeStyleSolid  = "solid"
	nodeStyleFilled = "filled"
)

// Edge style constants
const (
	edgeStyleSolid = "solid"
	// edgeStyleDashed = "dashed" // Example if needed later
)

// --- Structs (DotNodeSpec, DotEdgeSpec) remain the same ---

// DotNodeSpec represents a node in the DOT graph
type DotNodeSpec struct {
	Name        string
	DisplayName string
	Tooltip     string
	Shape       string
	Style       string // Use nodeStyleSolid or nodeStyleFilled
	FillColor   string
	// FontColor is handled conditionally by dotTemplateFuncs
}

// DotEdgeSpec represents an edge in the DOT graph
type DotEdgeSpec struct {
	FromNodeName string
	ToNodeName   string
	Tooltip      string
	Style        string // Use edgeStyleSolid, etc.
	Color        string
}

// --- Visualize Methods ---

// Visualize returns a DOT graph representation of the pipeline.
func (p *pipelineImplementation) Visualize() string {
	if len(p.nodes) == 0 {
		return dotTemplateFuncs([]*DotNodeSpec{}, []*DotEdgeSpec{})
	}

	nodes := make([]*DotNodeSpec, 0, len(p.nodes))
	edges := make([]*DotEdgeSpec, 0, len(p.nodes)-1)

	status, currentStepID, _ := getWorkflowStateInfo(p.state) // Use helper

	// Determine current step index once for edge coloring
	currentStepIndex := -1
	if currentStepID != "" {
		currentStepIndex = slices.IndexFunc(p.nodes, func(n RunnableInterface) bool {
			return n.GetID() == currentStepID
		})
	}

	// --- Create Node Specs ---
	for i, node := range p.nodes {
		isCurrentStep := currentStepID == node.GetID()
		// Use pipeline-specific node styling helper
		nodeStyle, fillColor := getPipelineNodeStyleAndColor(p.state, node.GetID(), i, len(p.nodes), isCurrentStep)
		nodes = append(nodes, createDotNodeSpec(node, nodeStyle, fillColor)) // Use helper
	}

	// --- Create Edge Specs ---
	for i := 1; i < len(p.nodes); i++ {
		fromNode := p.nodes[i-1]
		toNode := p.nodes[i]
		edgeStyle := edgeStyleSolid
		edgeColor := colorGrey // Default

		// Edge coloring logic (remains specific to pipeline)
		if status == StateStatusComplete ||
			(status == StateStatusRunning && currentStepIndex != -1 && i <= currentStepIndex) {
			edgeColor = colorGreen
		}
		edges = append(edges, createDotEdgeSpec(fromNode, toNode, edgeStyle, edgeColor)) // Use helper
	}

	return dotTemplateFuncs(nodes, edges)
}

// Visualize returns a DOT graph representation of the DAG.
func (d *Dag) Visualize() string {
	if len(d.runnables) == 0 {
		return dotTemplateFuncs([]*DotNodeSpec{}, []*DotEdgeSpec{})
	}

	status, currentStepID, completedSteps := getWorkflowStateInfo(d.state)
	nodes := d.createDagNodeSpecs(currentStepID, completedSteps)
	edges := d.createDagEdgeSpecs(status, completedSteps)

	return dotTemplateFuncs(nodes, edges)
}

// Visualize returns a DOT graph representation of the step.
func (s *stepImplementation) Visualize() string {
	nodeStyle, fillColor := nodeStyleSolid, colorWhite // Default

	status, _, _ := getWorkflowStateInfo(s.state) // Use helper for consistency

	// Determine style and color based on status
	switch status {
	case StateStatusRunning:
		nodeStyle = nodeStyleFilled
		fillColor = colorBlue
	case StateStatusComplete: // Single step complete is green
		nodeStyle = nodeStyleFilled
		fillColor = colorGreen
	case StateStatusFailed:
		nodeStyle = nodeStyleFilled
		fillColor = colorRed
	case StateStatusPaused:
		nodeStyle = nodeStyleFilled
		fillColor = colorYellow
	}

	// FIX: Use createDotNodeSpec to ensure consistent node creation including tooltip
	nodeSpec := createDotNodeSpec(s, nodeStyle, fillColor)

	edges := []*DotEdgeSpec{} // Steps have no edges

	return dotTemplateFuncs([]*DotNodeSpec{nodeSpec}, edges)
}

// --- Helper Functions ---

// escapeDotString escapes characters in a string for use in DOT labels/tooltips.
func escapeDotString(s string) string {
	quoted := fmt.Sprintf("%q", s)
	if len(quoted) >= 2 {
		// Remove the outer quotes added by %q
		return quoted[1 : len(quoted)-1]
	}
	return "" // Return empty string if input was empty or invalid
}

// getWorkflowStateInfo safely extracts common state information from a StateInterface.
// It returns default values if the state is nil.
func getWorkflowStateInfo(state StateInterface) (status StateStatus, currentStepID string, completedSteps []string) {
	status = StateStatus("") // Default empty status
	currentStepID = ""
	completedSteps = []string{} // Initialize as empty slice

	if state != nil {
		status = state.GetStatus()
		currentStepID = state.GetCurrentStepID()
		// Only get completed steps if state exists and is not nil
		if cs := state.GetCompletedSteps(); cs != nil {
			completedSteps = cs
		}
	}
	return status, currentStepID, completedSteps
}

// getNodeStyleAndColor determines the base style and fill color for a node
// primarily based on whether it's the *current* step and the workflow's status.
func getNodeStyleAndColor(state StateInterface, isCurrentStep bool) (style, fillColor string) {
	// Default styles for non-current steps or nil state
	style = nodeStyleSolid
	fillColor = colorWhite

	// Only apply special styling if it's the current step and state is valid
	if state == nil || !isCurrentStep {
		return style, fillColor
	}

	status := state.GetStatus()

	// Determine fill color based on status for the current step
	switch status {
	case StateStatusFailed:
		style = nodeStyleFilled
		fillColor = colorRed
	case StateStatusPaused:
		style = nodeStyleFilled
		fillColor = colorYellow
	case StateStatusRunning:
		style = nodeStyleFilled
		fillColor = colorBlue
	// Add other statuses for current step if needed, otherwise defaults apply
	default:
		// Keep default white/solid for current step in other states (e.g., Complete)
		// Or adjust if a different style is desired for a current-but-completed step.
	}

	return style, fillColor
}

// getPipelineNodeStyleAndColor determines the style and fill color for a node within a pipeline.
// It considers the overall pipeline status, the node's position, and whether it's the current step.
func getPipelineNodeStyleAndColor(state StateInterface, nodeID string, index, nodeCount int, isCurrentStep bool) (style, fillColor string) {
	// If it's the current step, use the general current step styling logic
	if isCurrentStep {
		return getNodeStyleAndColor(state, true)
	}

	// Default style for non-current steps
	style = nodeStyleSolid
	fillColor = colorWhite

	// Pipeline-specific logic for NON-CURRENT steps:
	status, _, completedSteps := getWorkflowStateInfo(state)

	// Style completed steps green if the pipeline is running or complete
	// In a pipeline, completion implies sequential execution up to that point.
	// We check `index < nodeCount-1` because the *last* node in a completed pipeline
	// might not need the green 'completed' look, depending on desired visualization.
	// Adjust this condition if the last node should also be green when complete.
	isCompleted := slices.Contains(completedSteps, nodeID)

	if (status == StateStatusRunning || status == StateStatusComplete) && isCompleted {
		// Make completed steps green, except potentially the very last one if status is Complete
		if status != StateStatusComplete || index < nodeCount-1 {
			style = nodeStyleFilled
			fillColor = colorGreen
		}
	}
	// Add more conditions here if needed for other states/styles

	return style, fillColor
}

// getDagNodeStyleAndColor determines the style and fill color for a node within a DAG.
// It considers the overall DAG status, whether the node is completed, and if it's the current step.
func getDagNodeStyleAndColor(state StateInterface, nodeID string, isCurrentStep bool, completedSteps []string) (style, fillColor string) {
	// If it's the current step, use the general current step styling logic
	if isCurrentStep {
		return getNodeStyleAndColor(state, true)
	}

	// Default style for non-current steps
	style = nodeStyleSolid
	fillColor = colorWhite

	if state == nil {
		return style, fillColor // Return default if state is nil
	}

	status := state.GetStatus()

	// Style completed steps green only when the DAG is actively running
	// In a completed or failed DAG, completed steps might revert to default or another style.
	if status == StateStatusRunning && slices.Contains(completedSteps, nodeID) {
		style = nodeStyleFilled
		fillColor = colorGreen
	}
	// Add more conditions here if needed for other states (e.g., different style for completed nodes in a Failed DAG)

	return style, fillColor
}

// createDotNodeSpec creates a DotNodeSpec struct with common defaults and provided style/color.
func createDotNodeSpec(node RunnableInterface, style, fillColor string) *DotNodeSpec {
	name := node.GetName()
	if name == "" {
		name = node.GetID() // Use ID if name is empty
	}
	return &DotNodeSpec{
		Name:        node.GetID(),
		DisplayName: name,
		Shape:       "box", // Common shape, can be customized if needed
		Style:       style,
		FillColor:   fillColor,
		Tooltip:     fmt.Sprintf("Step: %s", name), // Default tooltip
	}
}

// createDotEdgeSpec creates a DotEdgeSpec struct representing a directed edge.
func createDotEdgeSpec(fromNode, toNode RunnableInterface, style, color string) *DotEdgeSpec {
	fromName := fromNode.GetName()
	if fromName == "" {
		fromName = fromNode.GetID()
	}
	toName := toNode.GetName()
	if toName == "" {
		toName = toNode.GetID()
	}
	return &DotEdgeSpec{
		FromNodeName: fromNode.GetID(),
		ToNodeName:   toNode.GetID(),
		Style:        style,
		Color:        color,
		Tooltip:      fmt.Sprintf("From %s to %s", fromName, toName), // Tooltip indicating direction
	}
}

// dotTemplateFuncs generates the final DOT language string from node and edge specifications.
func dotTemplateFuncs(nodes []*DotNodeSpec, edges []*DotEdgeSpec) string {
	var sb strings.Builder

	sb.WriteString("digraph {\n")
	sb.WriteString("\trankdir = \"LR\"; // Left-to-right layout\n")
	sb.WriteString("\tnode [fontname=\"Arial\", shape=box]; // Default node attributes\n")
	sb.WriteString("\tedge [fontname=\"Arial\"]; // Default edge attributes\n")
	sb.WriteString("\n")

	// Add Nodes
	sb.WriteString("\t// Nodes\n")
	for _, node := range nodes {
		// Use DisplayName if available, otherwise Name (which defaults to ID if original name was empty)
		label := node.DisplayName
		tooltip := node.Tooltip // Use pre-formatted tooltip

		// Build node attributes string
		attrs := fmt.Sprintf("label=\"%s\", style=%s, tooltip=\"%s\", fillcolor=\"%s\"",
			escapeDotString(label),
			node.Style,
			escapeDotString(tooltip),
			node.FillColor,
		)
		// Add fontcolor=white for filled nodes for better contrast
		if node.Style == nodeStyleFilled {
			attrs += ", fontcolor=\"white\""
		}

		sb.WriteString(fmt.Sprintf("\t\"%s\" [%s];\n",
			escapeDotString(node.Name), // Node ID must be unique
			attrs,
		))
	}
	sb.WriteString("\n")

	// Add Edges
	sb.WriteString("\t// Edges\n")
	for _, edge := range edges {
		// Build edge attributes string
		attrs := fmt.Sprintf("style=%s, tooltip=\"%s\", color=\"%s\"",
			edge.Style,
			escapeDotString(edge.Tooltip),
			edge.Color,
		)

		sb.WriteString(fmt.Sprintf("\t\"%s\" -> \"%s\" [%s];\n",
			escapeDotString(edge.FromNodeName),
			escapeDotString(edge.ToNodeName),
			attrs,
		))
	}

	sb.WriteString("}\n")
	return sb.String()
}

// createDagEdgeSpecs generates the list of DotEdgeSpec for a DAG based on its dependencies and state.
func (d *Dag) createDagEdgeSpecs(status StateStatus, completedSteps []string) []*DotEdgeSpec {
	edges := make([]*DotEdgeSpec, 0) // Initialize empty slice

	// Iterate through the DAG's defined dependencies
	for dependentID, dependencyIDs := range d.dependencies {
		dependent, depExists := d.runnables[dependentID]
		if !depExists {
			continue // Skip if the dependent node doesn't exist in the runnables map
		}

		for _, dependencyID := range dependencyIDs {
			dependency, depExists2 := d.runnables[dependencyID]
			if !depExists2 {
				continue // Skip if the dependency node doesn't exist
			}

			// Determine edge style and color based on DAG status and completion
			edgeStyle := edgeStyleSolid // Default style
			edgeColor := colorGrey      // Default color

			isSourceCompleted := slices.Contains(completedSteps, dependencyID)

			// Color edge green if:
			// 1. The entire DAG is complete.
			// 2. The DAG is running AND the source node of the edge is completed.
			if status == StateStatusComplete || (status == StateStatusRunning && isSourceCompleted) {
				edgeColor = colorGreen
			}
			// Add more conditions for other edge colors/styles if needed (e.g., for failed/paused states)

			// Create and add the edge specification
			edges = append(edges, createDotEdgeSpec(dependency, dependent, edgeStyle, edgeColor))
		}
	}
	return edges
}

// createDagNodeSpecs generates the list of DotNodeSpec for a DAG based on its state.
func (d *Dag) createDagNodeSpecs(currentStepID string, completedSteps []string) []*DotNodeSpec {
	nodes := make([]*DotNodeSpec, 0, len(d.runnables)) // Pre-allocate slice capacity

	// Iterate through all runnable nodes in the DAG
	for _, node := range d.runnables {
		isCurrentStep := currentStepID == node.GetID()

		// Use the DAG-specific helper function to determine style and color
		nodeStyle, fillColor := getDagNodeStyleAndColor(d.state, node.GetID(), isCurrentStep, completedSteps)

		// Create and add the node specification
		nodes = append(nodes, createDotNodeSpec(node, nodeStyle, fillColor))
	}
	return nodes
}
