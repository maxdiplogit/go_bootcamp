// Package contextref is a personal reference package for the context
// package in Go.
//
// CORE MENTAL MODEL:
//
// `context.Context` is an INTERFACE for cancellation, deadlines, and
// request-scoped values. Contexts form a TREE: a root (Background) at
// the top, with derived children (WithCancel, WithTimeout, WithDeadline,
// WithValue). Cancellation propagates DOWN the tree -- if you cancel
// a parent, all descendants are canceled too. Cancellation never goes UP.
//
// Internally, the cancellation mechanism is just a channel: `ctx.Done()`
// returns a channel that's CLOSED when the context is canceled. Goroutines
// detect cancellation by reading from this channel (usually in a select).
//
// FIVE THINGS TO REMEMBER:
//
//  1. Always start with `context.Background()` at the top of a request
//     or main function. Derive children from it as you go.
//  2. `WithCancel`, `WithTimeout`, `WithDeadline` return (ctx, cancel).
//     ALWAYS call `defer cancel()` -- forgetting it leaks resources.
//  3. Pass context as the FIRST parameter, named `ctx`, by VALUE (not
//     pointer): `func Foo(ctx context.Context, ...)`.
//  4. Inside a goroutine, watch ctx.Done() in a `select`:
//     case <-ctx.Done(): return
//     This is how you actually stop work when canceled.
//  5. Use `WithValue` sparingly -- only for cross-cutting concerns like
//     request IDs and trace spans. Don't smuggle real arguments through it.
package contextref

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// SECTION 1: WHY: CHANNELS-ONLY VS CONTEXT
// =============================================================================

// CancelWithRawChannel demonstrates the OLD WAY (using a raw channel)
// to cancel a goroutine. It works, but doesn't scale: each function
// would need to take this channel as a parameter, and there's no built-in
// way to handle deadlines or compose cancellations.
//
// We close `stop` to signal cancellation. The worker goroutine watches
// the channel in a select and returns when the channel becomes readable
// (a closed channel is always readable).
func CancelWithRawChannel() string {
	stop := make(chan struct{})
	results := make(chan string, 1)

	go func() {
		// Loop until we receive a stop signal.
		for i := 0; ; i++ {
			select {
			case <-stop:
				results <- fmt.Sprintf("worker stopped after %d iterations", i)
				return
			default:
				time.Sleep(5 * time.Millisecond) // pretend to do work
			}
		}
	}()

	// Let the worker run a bit, then signal stop by closing the channel.
	time.Sleep(20 * time.Millisecond)
	close(stop)

	return <-results
}

// CancelWithContext demonstrates the SAME idea using context.WithCancel.
// The behavior is identical, but we get a standard interface that scales
// to deadlines, multiple cancellation sources, and value propagation
// without us having to design any of that ourselves.
func CancelWithContext() string {
	// context.Background() is the empty root context. We derive a
	// cancelable child from it.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // always pair with defer cancel

	results := make(chan string, 1)

	go func() {
		for i := 0; ; i++ {
			select {
			case <-ctx.Done(): // this channel closes when cancel() is called
				results <- fmt.Sprintf("worker stopped after %d iterations: %v", i, ctx.Err())
				return
			default:
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()

	time.Sleep(20 * time.Millisecond)
	cancel() // signals all goroutines watching ctx.Done() to stop

	return <-results
}

// =============================================================================
// SECTION 2: WithCancel -- MANUAL CANCELLATION
// =============================================================================

// SimulateBackgroundJob runs a "job" that loops, doing work, and stops
// when told via context. Returns how many work units it completed.
//
// This is the canonical worker shape: the for-select pattern.
//   - The default case does one unit of work.
//   - The ctx.Done() case is checked on every iteration.
//
// As soon as the context is canceled (cancel() called, or timeout, or
// any ancestor context being canceled), the next iteration exits.
func SimulateBackgroundJob(ctx context.Context) int {
	count := 0
	for {
		select {
		case <-ctx.Done():
			return count
		default:
			count++
			time.Sleep(5 * time.Millisecond)
		}
	}
}

// RunWithManualCancel demonstrates calling SimulateBackgroundJob and
// canceling it after a short period. Shows the explicit (cancel, ctx)
// pair returned by WithCancel.
func RunWithManualCancel() int {
	ctx, cancel := context.WithCancel(context.Background())

	// Schedule cancellation after 30ms. We could also just call cancel()
	// directly when we want to stop -- this just delays it slightly.
	time.AfterFunc(30*time.Millisecond, cancel)

	count := SimulateBackgroundJob(ctx)
	return count
}

// =============================================================================
// SECTION 3: WithTimeout -- AUTOMATIC CANCELLATION AFTER A DURATION
// =============================================================================

// FetchWithTimeout simulates an operation that might take too long. We
// use ctx.Done() inside a select to abort early if the deadline fires.
//
// This is the SHAPE of every "I/O with timeout" function in real Go code:
//   - Set up the context with a deadline (caller's responsibility usually)
//   - Inside the function, race the operation against ctx.Done()
//   - Return whichever comes first
//
// The error returned is exactly ctx.Err() -- which will be
// context.DeadlineExceeded if the timeout fired, or context.Canceled
// if someone canceled manually. The caller can use errors.Is to check.
func FetchWithTimeout(ctx context.Context, simulatedDelay time.Duration) (string, error) {
	select {
	case <-time.After(simulatedDelay):
		return "fetched data", nil
	case <-ctx.Done():
		return "", ctx.Err() // DeadlineExceeded or Canceled
	}
}

// RunFetchTooSlow shows the timeout firing before the operation completes.
// Returns the error so the demo can verify it's context.DeadlineExceeded.
func RunFetchTooSlow() error {
	// Give the operation only 50ms.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// But the operation needs 200ms -- it'll be cut short.
	_, err := FetchWithTimeout(ctx, 200*time.Millisecond)
	return err
}

// RunFetchInTime shows the operation completing before the timeout fires.
// Returns the result and any error.
func RunFetchInTime() (string, error) {
	// Give the operation 200ms.
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Operation takes only 30ms -- finishes well before timeout.
	return FetchWithTimeout(ctx, 30*time.Millisecond)
}

// =============================================================================
// SECTION 4: ERROR INSPECTION -- DEADLINE vs CANCELED
// =============================================================================

// ClassifyContextError shows how to distinguish between the two main
// reasons a context can be done:
//   - context.Canceled        -- someone called cancel() explicitly
//   - context.DeadlineExceeded -- a deadline or timeout fired
//
// We use errors.Is because in real code the error might have been wrapped
// by the time you receive it, and errors.Is correctly walks the chain.
func ClassifyContextError(err error) string {
	switch {
	case err == nil:
		return "no error"
	case errors.Is(err, context.Canceled):
		return "canceled by user"
	case errors.Is(err, context.DeadlineExceeded):
		return "deadline exceeded"
	default:
		return "some other error: " + err.Error()
	}
}

// =============================================================================
// SECTION 5: PROPAGATION -- THE TREE STRUCTURE
// =============================================================================

// ParentChildPropagation demonstrates that canceling a PARENT context
// automatically cancels the CHILD. This is the tree behavior you read
// about; it's why context is so powerful.
//
// Setup:
//
//	parent (cancelable) ── child (cancelable)
//
// We cancel the PARENT. Both ctx.Done() channels close. Both goroutines
// watching either context will return.
//
// We return the errors from both contexts after cancellation so the
// demo can verify both got canceled.
func ParentChildPropagation() (parentErr, childErr error) {
	parent, cancelParent := context.WithCancel(context.Background())
	defer cancelParent()

	child, cancelChild := context.WithCancel(parent) // derived from parent
	defer cancelChild()

	// Cancel the PARENT only.
	cancelParent()

	// Give propagation a tiny moment (in practice it's instant, but in
	// rare cases on a busy system the goroutines that close child.Done()
	// might run a microsecond later -- this Sleep is just defensive
	// for demo purposes).
	time.Sleep(1 * time.Millisecond)

	return parent.Err(), child.Err()
}

// FanOutWithCancel demonstrates the most common real-world use of context:
// one parent context spawns multiple worker goroutines, all of which
// stop when ANY of these things happens:
//   - The parent is canceled
//   - The deadline fires
//   - One of the workers explicitly calls the cancel function (rare)
//
// This is the structural pattern of pretty much every HTTP handler in
// production Go code: handler gets a context, spawns helpers (DB query,
// upstream API call, etc.), and they all share the same context.
func FanOutWithCancel(numWorkers int) []int {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	// Each worker reports how many iterations it completed.
	results := make([]int, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			count := 0
			for {
				select {
				case <-ctx.Done():
					results[id] = count
					return
				default:
					count++
					time.Sleep(2 * time.Millisecond)
				}
			}
		}(i)
	}
	wg.Wait()
	return results
}

// =============================================================================
// SECTION 6: WithValue -- REQUEST-SCOPED DATA (USE SPARINGLY)
// =============================================================================

// requestIDKey is the key used for storing/retrieving the request ID.
// CONVENTION: keys for context.WithValue should be of an UNEXPORTED CUSTOM
// TYPE. Why? Because if everyone used `string` keys, two packages could
// accidentally collide on the same string ("user_id"). An unexported
// custom type is unique to your package and impossible to spoof.
type contextKey int

const (
	requestIDKey contextKey = iota
	userKey
)

// AttachRequestID returns a NEW context (a derived one) that has the
// given request ID attached. The original ctx is unchanged.
//
// Note: WithValue takes a parent context and returns a child. The child
// has all the parent's behavior (cancellation, deadline) PLUS this value.
func AttachRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// AttachUser is the same idea for a user value.
func AttachUser(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, userKey, name)
}

// GetRequestID retrieves the request ID from the context, if any. The
// type assertion `.(string)` recovers the concrete type. If the key
// isn't present, ctx.Value returns nil and the assertion fails -- we
// use the two-value form to detect that gracefully.
func GetRequestID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(requestIDKey).(string)
	return v, ok
}

// GetUser is the same for the user value.
func GetUser(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userKey).(string)
	return v, ok
}

// SimulateRequestHandler is a tiny example of how WithValue is used in
// a typical web handler. The HTTP middleware attaches a request ID;
// downstream code can retrieve it for logging.
//
// CRUCIALLY: real arguments (the user ID, the request body, etc.)
// should NOT be passed through the context. Pass them as regular
// function arguments. WithValue is for cross-cutting concerns only --
// things you want to pass implicitly through every layer for logging
// or tracing, not real domain inputs.
func SimulateRequestHandler() string {
	// At the top of a request handler, we usually:
	//   1. Take the incoming context
	//   2. Attach a request ID and any auth info
	//   3. Pass it down through the call stack
	ctx := context.Background()
	ctx = AttachRequestID(ctx, "req-12345")
	ctx = AttachUser(ctx, "asha")

	// Deeper in the call stack, code can pull out what it needs:
	id, _ := GetRequestID(ctx)
	user, _ := GetUser(ctx)
	return fmt.Sprintf("handling request %s for user %s", id, user)
}

// =============================================================================
// SECTION 7: PUTTING IT ALL TOGETHER -- A REALISTIC HANDLER
// =============================================================================

// Database is a fake database with a deliberately-slow Query method
// that respects context cancellation.
type Database struct{}

// Query simulates a database query. It either completes after queryTime
// or aborts when the context is done -- whichever comes first.
//
// This is the SHAPE every real I/O operation in Go follows: takes a
// context, races its work against ctx.Done(), returns ctx.Err() if
// the context fires first.
func (d *Database) Query(ctx context.Context, queryTime time.Duration) (string, error) {
	select {
	case <-time.After(queryTime):
		return "query result", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// HandleRequest is a realistic mini-handler. It:
//  1. Receives a context (typically from net/http)
//  2. Adds a request-specific timeout to limit total work
//  3. Attaches a request ID for downstream logging
//  4. Calls a sub-operation (db query) that respects the context
//  5. Returns either the result or a context error
//
// Notice how the SAME ctx flows through every layer. If the user
// disconnects, if the deadline fires, if anything cancels -- everything
// downstream stops automatically.
func HandleRequest(parentCtx context.Context, requestID string) (string, error) {
	// Layer our own timeout on top of whatever the caller gave us.
	// If the parent already has a shorter deadline, the parent wins
	// (the child is bounded by the parent's deadline as well as its own).
	ctx, cancel := context.WithTimeout(parentCtx, 100*time.Millisecond)
	defer cancel()

	// Attach the request ID for downstream visibility.
	ctx = AttachRequestID(ctx, requestID)

	// Do the work. The Database.Query function respects context, so if
	// our 100ms timeout fires (or the parent cancels), the query is
	// cut short.
	db := &Database{}
	result, err := db.Query(ctx, 50*time.Millisecond)
	if err != nil {
		// Pull the request ID back out for the error message.
		id, _ := GetRequestID(ctx)
		return "", fmt.Errorf("request %s failed: %w", id, err)
	}

	id, _ := GetRequestID(ctx)
	return fmt.Sprintf("[req=%s] %s", id, result), nil
}
