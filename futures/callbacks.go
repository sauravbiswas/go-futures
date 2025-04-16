package futures

// OnSuccess registers a callback function to be called when the future completes successfully.
func (f *Future[T]) OnSuccess(cb func(T)) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	// If the future is already fulfilled, execute the callback immediately
	if f.state == Fulfilled {
		cb(f.result)
		return
	}

	f.onSuccess = append(f.onSuccess, cb)
}

// OnFailure registers a callback function to be called when the future completes with an error.
func (f *Future[T]) OnFailure(cb func(error)) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	// If the future is already rejected, execute the callback immediately
	if f.state == Rejected {
		cb(f.err)
		return
	}

	f.onFailure = append(f.onFailure, cb)
}
