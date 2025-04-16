package futures

// State represents the possible states of a Future
type State int

const (
	Pending   State = iota // Initial state
	Running                // Task is currently executing
	Fulfilled              // Task completed successfully
	Rejected               // Task completed with an error
)
