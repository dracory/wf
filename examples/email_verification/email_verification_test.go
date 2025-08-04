package email_verification_test

// import (
// 	"context"
// 	"testing"

// 	"github.com/dracory/wf"
// )

// func TestEmailVerificationWorkflow(t *testing.T) {
// 	// Create workflow
// 	dag := NewEmailVerificationWorkflow()

// 	// Initialize data
// 	ctx := context.Background()
// 	// Add the DAG to the context so the SendEmailStep can access it
// 	ctx = context.WithValue(ctx, dagKey, dag)
// 	data := map[string]any{
// 		"email": "test@example.com",
// 	}

// 	// Start workflow
// 	ctx, data, err := dag.Run(ctx, data)
// 	if err != nil {
// 		t.Fatalf("workflow failed: %v", err)
// 	}

// 	// Verify email was sent and code was generated
// 	if code, ok := data["verificationCode"].(string); !ok {
// 		t.Fatal("verification code not generated")
// 	} else if len(code) != 6 {
// 		t.Fatalf("expected 6-digit code, got %s", code)
// 	}

// 	// Verify workflow is paused
// 	if !dag.IsPaused() {
// 		t.Fatal("workflow is not paused")
// 	}

// 	// Save workflow state
// 	state := dag.GetState()
// 	stateJSON, err := state.ToJSON()
// 	if err != nil {
// 		t.Fatalf("failed to save state: %v", err)
// 	}
// 	t.Logf("Saved workflow state: %s", string(stateJSON))

// 	// Create a new workflow instance
// 	newDag := NewEmailVerificationWorkflow()

// 	// Load saved state
// 	newState := wf.NewState()
// 	if err := newState.FromJSON(stateJSON); err != nil {
// 		t.Fatalf("failed to load state: %v", err)
// 	}
// 	newDag.SetState(newState)

// 	// Verify the new workflow is paused
// 	if !newDag.IsPaused() {
// 		t.Fatal("new workflow is not paused")
// 	}

// 	// Set the entered code to match the verification code
// 	data["enteredCode"] = data["verificationCode"]

// 	// Resume workflow
// 	_, _, err = newDag.Resume(ctx, data)
// 	if err != nil {
// 		t.Fatalf("workflow resume failed: %v", err)
// 	}

// 	// Verify workflow completed
// 	if !newDag.IsCompleted() {
// 		t.Fatal("workflow did not complete")
// 	}

// 	// Verify the verification was successful
// 	if verified, ok := data["verified"].(bool); !ok || !verified {
// 		t.Fatal("email verification failed")
// 	}
// }

// func TestEmailVerificationWorkflowInvalidCode(t *testing.T) {
// 	// Create workflow
// 	dag := NewEmailVerificationWorkflow()

// 	// Initialize data
// 	ctx := context.Background()
// 	// Add the DAG to the context so the SendEmailStep can access it
// 	ctx = context.WithValue(ctx, dagKey, dag)
// 	data := map[string]any{
// 		"email": "test@example.com",
// 	}

// 	// Start workflow
// 	ctx, data, err := dag.Run(ctx, data)
// 	if err != nil {
// 		t.Fatalf("workflow failed: %v", err)
// 	}

// 	// Verify workflow is paused
// 	if !dag.IsPaused() {
// 		t.Fatal("workflow is not paused")
// 	}

// 	// Save workflow state
// 	state := dag.GetState()
// 	stateJSON, err := state.ToJSON()
// 	if err != nil {
// 		t.Fatalf("failed to save state: %v", err)
// 	}

// 	// Create a new workflow instance
// 	newDag := NewEmailVerificationWorkflow()

// 	// Load saved state
// 	newState := wf.NewState()
// 	if err := newState.FromJSON(stateJSON); err != nil {
// 		t.Fatalf("failed to load state: %v", err)
// 	}
// 	newDag.SetState(newState)

// 	// Set an invalid code
// 	data["enteredCode"] = "000000"

// 	// Resume workflow - should fail due to invalid code
// 	_, _, err = newDag.Resume(ctx, data)
// 	if err == nil {
// 		t.Fatal("expected error for invalid verification code")
// 	}
// }

// func TestEmailVerificationWorkflowE2E(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		email       string
// 		wantErr     bool
// 		description string
// 	}{
// 		{
// 			name:        "Valid email verification",
// 			email:       "test@example.com",
// 			wantErr:     false,
// 			description: "Should complete successfully with valid email",
// 		},
// 		{
// 			name:        "Empty email",
// 			email:       "",
// 			wantErr:     true,
// 			description: "Should fail with empty email",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Create workflow
// 			dag := NewEmailVerificationWorkflow()

// 			// Initialize data
// 			ctx := context.Background()
// 			// Add the DAG to the context so the SendEmailStep can access it
// 			ctx = context.WithValue(ctx, dagKey, dag)
// 			data := map[string]any{
// 				"email": tt.email,
// 			}

// 			// Start workflow
// 			ctx, data, err := dag.Run(ctx, data)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("workflow failed: %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			// Skip further checks if we expected an error
// 			if tt.wantErr {
// 				return
// 			}

// 			// Verify email was sent and code was generated
// 			if code, ok := data["verificationCode"].(string); !ok {
// 				t.Fatal("verification code not generated")
// 			} else if len(code) != 6 {
// 				t.Fatalf("expected 6-digit code, got %s", code)
// 			}

// 			// Verify workflow is paused
// 			if !dag.IsPaused() {
// 				t.Fatal("workflow is not paused")
// 			}

// 			// Save workflow state
// 			state := dag.GetState()
// 			stateJSON, err := state.ToJSON()
// 			if err != nil {
// 				t.Fatalf("failed to save state: %v", err)
// 			}
// 			t.Logf("Saved workflow state: %s", string(stateJSON))

// 			// Create a new workflow instance
// 			newDag := NewEmailVerificationWorkflow()

// 			// Load saved state
// 			newState := wf.NewState()
// 			if err := newState.FromJSON(stateJSON); err != nil {
// 				t.Fatalf("failed to load state: %v", err)
// 			}
// 			newDag.SetState(newState)

// 			// Verify the new workflow is paused
// 			if !newDag.IsPaused() {
// 				t.Fatal("new workflow is not paused")
// 			}

// 			// Set the entered code to match the verification code
// 			data["enteredCode"] = data["verificationCode"]

// 			// Resume workflow
// 			_, _, err = newDag.Resume(ctx, data)
// 			if err != nil {
// 				t.Fatalf("workflow resume failed: %v", err)
// 			}

// 			// Verify workflow completed
// 			if !newDag.IsCompleted() {
// 				t.Fatal("workflow did not complete")
// 			}

// 			// Verify the verification was successful
// 			if verified, ok := data["verified"].(bool); !ok || !verified {
// 				t.Fatal("email verification failed")
// 			}
// 		})
// 	}
// }

// // TestEmailVerificationWorkflowPauseResume demonstrates the pause/resume functionality
// func TestEmailVerificationWorkflowPauseResume(t *testing.T) {
// 	// Create workflow
// 	dag := NewEmailVerificationWorkflow()

// 	// Initialize data
// 	ctx := context.Background()
// 	// Add the DAG to the context so the SendEmailStep can access it
// 	ctx = context.WithValue(ctx, dagKey, dag)
// 	data := map[string]any{
// 		"email": "test@example.com",
// 	}

// 	// Start workflow
// 	ctx, data, err := dag.Run(ctx, data)
// 	if err != nil {
// 		t.Fatalf("workflow failed: %v", err)
// 	}

// 	// Verify email was sent and code was generated
// 	if code, ok := data["verificationCode"].(string); !ok {
// 		t.Fatal("verification code not generated")
// 	} else if len(code) != 6 {
// 		t.Fatalf("expected 6-digit code, got %s", code)
// 	}

// 	// Verify workflow is paused
// 	if !dag.IsPaused() {
// 		t.Fatal("workflow is not paused")
// 	}

// 	// Save workflow state
// 	state := dag.GetState()
// 	stateJSON, err := state.ToJSON()
// 	if err != nil {
// 		t.Fatalf("failed to save state: %v", err)
// 	}
// 	t.Logf("Saved workflow state: %s", string(stateJSON))

// 	// Create a new workflow instance
// 	newDag := NewEmailVerificationWorkflow()

// 	// Load saved state
// 	newState := wf.NewState()
// 	if err := newState.FromJSON(stateJSON); err != nil {
// 		t.Fatalf("failed to load state: %v", err)
// 	}
// 	newDag.SetState(newState)

// 	// Verify the new workflow is paused
// 	if !newDag.IsPaused() {
// 		t.Fatal("new workflow is not paused")
// 	}

// 	// Set the entered code to match the verification code
// 	data["enteredCode"] = data["verificationCode"]

// 	// Resume workflow
// 	_, _, err = newDag.Resume(ctx, data)
// 	if err != nil {
// 		t.Fatalf("workflow resume failed: %v", err)
// 	}

// 	// Verify workflow completed
// 	if !newDag.IsCompleted() {
// 		t.Fatal("workflow did not complete")
// 	}

// 	// Verify the verification was successful
// 	if verified, ok := data["verified"].(bool); !ok || !verified {
// 		t.Fatal("email verification failed")
// 	}
// }
