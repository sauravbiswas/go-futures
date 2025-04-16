package futures

import (
	"sync"
)

// Future represents the result of an asynchronous computation.
type Future[T any] struct {
	task      func() (T, error)
	mutex     sync.Mutex
	result    T
	err       error
	state     State
	done      chan struct{}
	next      interface{} // Link to the next future in the chain (any type)
	onSuccess []func(T)
	onFailure []func(error)
	started   bool
}

// NewFuture creates a new Future instance.
func NewFuture[T any](task func() (T, error)) *Future[T] {
	return &Future[T]{
		task:      task,
		state:     Pending,
		done:      make(chan struct{}),
		onSuccess: []func(T){},
		onFailure: []func(error){},
		started:   false,
	}
}

// Then chains a new computation step to the current Future.
func (f *Future[T]) Then(nextTask func(T) (any, error)) *Future[any] {
	// Make sure the current future is started
	if !f.started {
		go f.Start()
	}

	nextFuture := NewFuture(func() (any, error) {
		// Wait for the parent future to complete
		result, err := f.Result()
		if err != nil {
			return nil, err
		}

		// Execute the next task with the result from the parent
		return nextTask(result)
	})

	// Link futures for debugging/tracing
	f.mutex.Lock()
	f.next = nextFuture
	f.mutex.Unlock()

	return nextFuture
}

// Start starts the future by executing the task asynchronously.
func (f *Future[T]) Start() {
	f.mutex.Lock()
	if f.state != Pending || f.started {
		f.mutex.Unlock()
		return
	}
	f.started = true
	f.state = Running
	f.mutex.Unlock()

	go func() {
		res, err := f.task()

		f.mutex.Lock()
		if f.state != Running { // Check if already completed somehow
			f.mutex.Unlock()
			return
		}

		if err != nil {
			f.err = err
			f.state = Rejected
			// Execute failure callbacks
			callbacks := make([]func(error), len(f.onFailure))
			copy(callbacks, f.onFailure)
			f.mutex.Unlock()

			// Execute callbacks outside the lock
			for _, cb := range callbacks {
				cb(err)
			}
		} else {
			f.result = res
			f.state = Fulfilled
			// Execute success callbacks
			callbacks := make([]func(T), len(f.onSuccess))
			copy(callbacks, f.onSuccess)
			f.mutex.Unlock()

			// Execute callbacks outside the lock
			for _, cb := range callbacks {
				cb(res)
			}
		}

		// Signal completion after callbacks
		close(f.done)
	}()
}

// Result returns the result of the future computation.
// It starts the future if it hasn't been started yet.
func (f *Future[T]) Result() (T, error) {
	// Auto-start if not already started
	if !f.started {
		f.Start()
	}

	<-f.done // Wait for completion

	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.result, f.err
}

// State returns the current state of the future.
func (f *Future[T]) State() State {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.state
}

// GetDone returns the done channel (used internally for chaining)
func (f *Future[T]) GetDone() chan struct{} {
	return f.done
}

// GetNext returns the next future in the chain (used internally for chaining)
func (f *Future[T]) GetNext() interface{} {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.next
}
