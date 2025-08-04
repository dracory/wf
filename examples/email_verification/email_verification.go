package email_verification

import (
	"context"
	"fmt"

	"github.com/dracory/wf"
)

// contextKey is a custom type for context keys
type contextKey string

const dagKey contextKey = "dag"

// NewSendEmailStep creates a step that sends a verification email
func NewSendEmailStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Send Verification Email")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		email, ok := data["email"].(string)
		if !ok {
			return ctx, data, fmt.Errorf("email is required")
		}

		// Validate email
		if email == "" {
			return ctx, data, fmt.Errorf("email cannot be empty")
		}

		code := generateVerificationCode()
		data["verificationCode"] = code

		// In a real application, this would send an actual email
		fmt.Printf("Sending verification code %s to %s\n", code, email)

		return ctx, data, nil
	})
	return step
}

// NewWaitForVerificationStep creates a step that waits for user input
func NewWaitForVerificationStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Wait for Verification")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		// In a real application, this would wait for user input
		// For this example, we'll just pass through
		return ctx, data, nil
	})
	return step
}

// NewVerifyCodeStep creates a step that verifies the entered code
func NewVerifyCodeStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Verify Code")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		// Check if we have an entered code
		enteredCode, ok := data["enteredCode"].(string)
		if !ok {
			// If no code entered yet, pause the workflow
			if dag, ok := ctx.Value(dagKey).(wf.DagInterface); ok {
				if dag.IsRunning() {
					if err := dag.Pause(); err != nil {
						return ctx, data, fmt.Errorf("failed to pause workflow: %v", err)
					}
					fmt.Println("Workflow paused waiting for verification code")
				}
			}
			return ctx, data, nil
		}

		expectedCode := data["verificationCode"].(string)

		if enteredCode != expectedCode {
			return ctx, data, fmt.Errorf("invalid verification code")
		}

		data["verified"] = true
		return ctx, data, nil
	})
	return step
}

// NewCompleteStep creates a step that completes the workflow
func NewCompleteStep() wf.StepInterface {
	step := wf.NewStep()
	step.SetName("Complete")
	step.SetHandler(func(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
		fmt.Println("Email verification completed successfully!")
		return ctx, data, nil
	})
	return step
}

// generateVerificationCode generates a random 6-digit code
func generateVerificationCode() string {
	// Use a fixed seed for testing purposes
	return "123456"
}

// NewEmailVerificationWorkflow creates a workflow for email verification
func NewEmailVerificationWorkflow() wf.DagInterface {
	dag := wf.NewDag()
	dag.SetName("Email Verification Workflow")

	// Create steps
	sendEmail := NewSendEmailStep()
	waitForVerification := NewWaitForVerificationStep()
	verifyCode := NewVerifyCodeStep()
	complete := NewCompleteStep()

	// Add steps to DAG
	dag.RunnableAdd(sendEmail, waitForVerification, verifyCode, complete)

	// Set up dependencies
	dag.DependencyAdd(waitForVerification, sendEmail)
	dag.DependencyAdd(verifyCode, waitForVerification)
	dag.DependencyAdd(complete, verifyCode)

	return dag
}

// RunEmailVerificationExample demonstrates the email verification workflow
func RunEmailVerificationExample() error {
	// Create workflow
	dag := NewEmailVerificationWorkflow()

	// Initialize data
	ctx := context.Background()
	// Add the DAG to the context so the SendEmailStep can access it
	ctx = context.WithValue(ctx, dagKey, dag)
	data := map[string]any{
		"email": "user@example.com",
	}

	// Start workflow
	ctx, data, err := dag.Run(ctx, data)
	if err != nil {
		return fmt.Errorf("workflow failed: %v", err)
	}

	// Save workflow state
	state := dag.GetState()
	stateJSON, err := state.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to save state: %v", err)
	}
	fmt.Printf("Saved workflow state: %s\n", string(stateJSON))

	// Create a new workflow instance
	newDag := NewEmailVerificationWorkflow()

	// Load saved state
	newState := wf.NewState()
	if err := newState.FromJSON(stateJSON); err != nil {
		return fmt.Errorf("failed to load state: %v", err)
	}
	newDag.SetState(newState)

	// Resume workflow with the updated context
	_, _, err = newDag.Resume(ctx, data)
	if err != nil {
		return fmt.Errorf("workflow resume failed: %v", err)
	}

	return nil
}
