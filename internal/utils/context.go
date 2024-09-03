package utils

import (
	"context"
	"time"
)

// NewTimeoutContext creates a new context with a specified timeout.
// It returns the context and the associated cancel function.
func NewTimeoutContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// NewCancelableContext creates a new context that can be canceled manually.
// It returns the context and the associated cancel function.
func NewCancelableContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

// SleepWithContext pauses the execution for the given duration or until the context is canceled.
// Useful for controlled delays within a context-aware process.
func SleepWithContext(ctx context.Context, duration time.Duration) error {
	select {
	case <-time.After(duration):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WaitForSignal waits for a signal to complete or until the context is canceled.
func WaitForSignal(ctx context.Context, signalChan <-chan struct{}) error {
	select {
	case <-signalChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
