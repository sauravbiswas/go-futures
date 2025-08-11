# go-futures

> 💡 A lightweight, composable concurrency library for Go, inspired by the elegance of Ruby's `concurrent-ruby` futures and promises.

---
## 🌟 Overview

Golang is famous for its first-class support for concurrency via goroutines and channels. However, writing structured, composable, and callback-enabled asynchronous code can often lead to verbose or tightly-coupled designs.

Libraries like [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup), [tunny](https://github.com/Jeffail/tunny), and [ants](https://github.com/panjf2000/ants) offer great pooling and job scheduling features — but they do **not** provide a **composable, future-based API** that can be chained, inspected, or observed with callbacks.

That’s where `go-futures` comes in.

---

## 🚀 Why `go-futures`?

Many existing concurrency libraries in Go **miss the following capabilities**:

| Capability                         | Missing in others | Included in `go-futures` |
|-----------------------------------|--------------------|---------------------------|
| Chainable future-like execution   | ✅                 | ✅                        |
| Success and failure callbacks     | ✅                 | ✅                        |
| Lightweight, zero-dependency core | ❌                 | ✅                        |
| Ability to inspect state          | ✅                 | ✅                        |
| Designed for flexibility          | ❌                 | ✅                        |


## 🔗 Why Chaining Async Operations Matters (Especially in Go)

Chaining asynchronous operations — like `.then().then()` — enables developers to express **sequential dependencies** between async tasks in a clear, composable way. This is especially useful in systems that handle **concurrent IO, network calls, or CPU-heavy operations**.

### ✅ Benefits of Chaining in Async Systems:

1. **Readability & Structure**  
   Chained syntax naturally expresses “do this, then that, then that...” — reducing nesting and callback hell.

2. **Error Propagation**  
   Errors can be caught and handled **in one place**, making the logic more maintainable and predictable.

3. **Reusability & Composition**  
   Functions returning futures/promises can be **composed and reused** like LEGO blocks, enabling flexible pipelines.

---

### ⚙️ Why This Is Useful in Go

Go emphasizes concurrency with goroutines and channels, but lacks **native support for chaining futures or promises**. This often leads to:

- Deeply nested code (`go func() { ... }` inside more `go func() { ... }`)
- Manual error propagation between goroutines
- Complex coordination logic using `sync.WaitGroup`, channels, or context timeouts

By introducing **chaining mechanisms**, Go developers can:

- **Compose concurrent workflows** more declaratively
- Avoid callback nesting and excessive channel usage
- **Encapsulate retry logic, timeouts, and result propagation** within the chain

---

###  Example: Chained Async Operations in Go

Here's an example using a hypothetical `futures` package that supports chaining like `Promise` or `CompletableFuture`:

```go
	// Create a future and chain multiple computations
	future := futures.NewFuture(func() (any, error) {
		// Call some API which is blocking
        time.Sleep(3 * time.Second)
		return "Start", nil
	}).Then(func(res any) (any, error) {
		// Call some I/O operation, here calling sleep
        time.sleep(2 * time.Second)
		str, ok := res.(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", res)
		}
		return str + " → Step 1", nil
	}).Then(func(res any) (any, error) {
		// Call some other I/O operation
        time.Sleep(4 * time.Second)
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
```

### ✨ Features
* Simple Future[T] abstraction with type safety

* `.Then(...)` chaining

* `.OnSuccess(...)` and `.OnFailure(...)` callbacks

* State introspection (Pending, Running, Fulfilled, Rejected)

* Fully tested with go test


### 🛠️ Installation

```bash
go get github.com/sauravbiswas/go-futures
```

### 🧪 Running Tests
```bash
go test ./futures -v
```

### 📚 Inspired By
* concurrent-ruby

* Go’s simplicity and goroutines

* Composable future/promise patterns from modern programming languages
