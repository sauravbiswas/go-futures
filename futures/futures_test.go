package futures_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sauravbiswasiupr/go-futures/futures"
	"github.com/stretchr/testify/assert"
)

func TestMultipleChainedFutures(t *testing.T) {
	// Create a future and chain multiple computations
	future := futures.NewFuture(func() (any, error) {
		time.Sleep(100 * time.Millisecond)
		return "Start", nil
	}).Then(func(res any) (any, error) {
		time.Sleep(200 * time.Millisecond)
		str, ok := res.(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", res)
		}
		return str + " → Step 1", nil
	}).Then(func(res any) (any, error) {
		time.Sleep(300 * time.Millisecond)
		str, ok := res.(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", res)
		}
		return str + " → Step 2", nil
	}).Then(func(res any) (any, error) {
		str, ok := res.(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", res)
		}
		return str + " → Step 3", nil
	})

	// Start the chain
	future.Start()

	// Wait for the result
	result, err := future.Result()

	// Using testify assertions
	assert.NoError(t, err)
	expected := "Start → Step 1 → Step 2 → Step 3"
	assert.Equal(t, expected, result)
}

func TestErrorPropagation(t *testing.T) {
	// Create a future that returns an error
	expectedError := "intentional failure in first future"
	future := futures.NewFuture(func() (any, error) {
		return nil, fmt.Errorf("%s", expectedError)
	}).Then(func(res any) (any, error) {
		// This should never execute due to the error in the first future
		time.Sleep(500 * time.Millisecond)
		t.Error("This code should not be reached!")
		return "This should not be returned", nil
	})

	// Start the chain
	future.Start()

	// Wait for the result
	result, err := future.Result()

	// The error from the first future should be propagated
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError)
	assert.Nil(t, result)
}
