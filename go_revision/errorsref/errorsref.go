// Package errorsref is a personal reference package for error handling in Go.
//
// CORE MENTAL MODEL:
//
// In Go, errors are VALUES, not control-flow events like exceptions. A
// function that can fail returns an extra `error` value, and the caller
// checks it explicitly with `if err != nil { ... }`. There is no
// try/catch -- every failure is visible right where it happens.
//
// FIVE THINGS TO REMEMBER:
//
//  1. The `error` type is just an interface with one method: `Error() string`.
//     Any type with that method is an error.
//  2. `errors.New(msg)` creates a basic error. `fmt.Errorf("...", args)`
//     creates a formatted one. Use `fmt.Errorf("... %w", err)` to WRAP
//     an existing error with extra context, preserving the original.
//  3. Use `errors.Is(err, target)` to check if a sentinel error is
//     anywhere in the wrap chain. Use `errors.As(err, &target)` to
//     pull out a specific custom error TYPE from the chain.
//  4. `defer fn()` schedules `fn` to run when the surrounding function
//     returns, regardless of how. Used for cleanup that must always
//     happen (close files, unlock mutexes, log errors).
//  5. `panic` halts execution and unwinds the stack; `recover` (only
//     inside a deferred function) catches a panic. Reserve panic for
//     truly unrecoverable bugs, NOT regular failures.
//
// USE plain errors for almost everything. Use panic only for "this should
// never happen" situations.
package errorsref

import (
	"errors"
	"fmt"
	"strings"
)

// =============================================================================
// SECTION 1: BASIC ERROR HANDLING -- (result, error) PATTERN
// =============================================================================

// Divide returns the quotient of a/b, or an error if b is zero.
// This is the canonical Go function shape: real result + error.
//
// The caller MUST check the error before using the result. If err != nil,
// the result is meaningless (we return 0 here as a placeholder).
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		// errors.New is the simplest way to create an error from a string.
		// It returns a value of an unexported type whose Error() method
		// just returns the string you passed.
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil // nil error means "everything is fine"
}

// Lookup returns a value from a small map, or an error if not found.
// Demonstrates fmt.Errorf for FORMATTED error messages -- like fmt.Sprintf
// but it returns an `error` instead of a string.
//
// Use fmt.Errorf when you want to interpolate dynamic data into the
// error message (the missing key, in this case).
func Lookup(m map[string]int, key string) (int, error) {
	v, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("key %q not found in map", key)
	}
	return v, nil
}

// =============================================================================
// SECTION 2: SENTINEL ERRORS -- EXPORTED ERROR VALUES
// =============================================================================

// ErrUserNotFound is a SENTINEL ERROR. It's a specific error VALUE that
// callers can compare against to detect a particular failure mode.
//
// Convention: sentinel errors are exported (capitalized), prefixed with
// "Err", and declared at the package level. Examples in stdlib:
// io.EOF, sql.ErrNoRows, os.ErrNotExist.
//
// Use sentinels when callers might legitimately want to react to THIS
// SPECIFIC failure differently from other failures.
var ErrUserNotFound = errors.New("user not found")

// ErrInvalidEmail is another sentinel for a different error case.
var ErrInvalidEmail = errors.New("invalid email address")

// FindUser is a tiny demo function that returns ErrUserNotFound for
// missing users. The caller can use `errors.Is(err, ErrUserNotFound)`
// (covered below) to detect this specific case.
func FindUser(id int) (string, error) {
	users := map[int]string{1: "Asha", 2: "Bo"}
	name, ok := users[id]
	if !ok {
		return "", ErrUserNotFound // bare sentinel, no wrapping
	}
	return name, nil
}

// =============================================================================
// SECTION 3: CUSTOM ERROR TYPES -- STRUCTURED ERRORS
// =============================================================================

// ValidationError is a CUSTOM error type. It satisfies the error interface
// because it has an Error() string method. Beyond that, it carries
// structured data (Field and Message) that callers can read directly --
// useful for things like building API error responses.
//
// Convention: name types ending in "Error" (e.g. ValidationError,
// PathError, NetworkError).
type ValidationError struct {
	Field   string
	Message string
}

// Error makes *ValidationError satisfy the error interface. The receiver
// is a pointer because the type is normally used via pointers (so that
// errors.As can match it). All ValidationError methods should use
// pointer receivers for consistency.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %q: %s", e.Field, e.Message)
}

// ValidateEmail returns a *ValidationError if the email looks invalid.
// We don't go through fmt.Errorf or wrapping -- we just construct the
// struct and return it. The caller can either treat it as a plain error
// (just print it) or extract its fields with errors.As (see Section 5).
func ValidateEmail(email string) error {
	if email == "" {
		return &ValidationError{Field: "email", Message: "must not be empty"}
	}
	if !strings.Contains(email, "@") {
		return &ValidationError{Field: "email", Message: "must contain @"}
	}
	return nil
}

// =============================================================================
// SECTION 4: WRAPPING ERRORS WITH %w
// =============================================================================

// LoadUserConfig demonstrates ERROR WRAPPING. As an error travels up the
// call stack, each layer wraps it with extra context using fmt.Errorf
// and the special %w verb.
//
// %w is different from %v:
//
//	%v -- includes the error's TEXT in the message (no link to original)
//	%w -- includes the error's text AND keeps a link to the original
//	      error so errors.Is and errors.As can walk the chain
//
// Always use %w when wrapping. %v is only correct when you want to
// "obscure" the original (rare).
func LoadUserConfig(userID int) (string, error) {
	_, err := FindUser(userID)
	if err != nil {
		// Wrap: "loadUserConfig 5: user not found"
		// The original ErrUserNotFound is preserved INSIDE this new error,
		// reachable via errors.Is.
		return "", fmt.Errorf("LoadUserConfig %d: %w", userID, err)
	}
	// ... pretend we load config from somewhere ...
	return "", nil
}

// StartApp wraps the error one more time. Now the chain looks like:
//
//	"StartApp: LoadUserConfig 5: user not found"
//
// The original ErrUserNotFound is still findable at the bottom of the chain.
func StartApp() error {
	_, err := LoadUserConfig(99) // 99 doesn't exist
	if err != nil {
		return fmt.Errorf("StartApp: %w", err)
	}
	return nil
}

// =============================================================================
// SECTION 5: UNWRAPPING -- errors.Is AND errors.As
// =============================================================================

// CheckIfNotFound demonstrates errors.Is. It walks the wrap chain
// looking for a specific sentinel value. Returns true if `target` is
// anywhere in the chain.
//
// Use errors.Is (NOT plain ==) whenever the error might have been
// wrapped on its way to you. Plain `err == ErrUserNotFound` would fail
// against a wrapped error -- the outer error and ErrUserNotFound are
// different values.
func CheckIfNotFound(err error) bool {
	return errors.Is(err, ErrUserNotFound)
}

// ExtractValidationError demonstrates errors.As. It walks the wrap
// chain looking for an error of a specific TYPE, and assigns it to
// the variable you pass.
//
// errors.As takes a pointer to the variable -- it needs to write to it.
// If the chain contains a *ValidationError anywhere, ve becomes that
// error and we can read its Field and Message.
//
// Returns ("", "", false) if no *ValidationError is in the chain.
func ExtractValidationError(err error) (field, message string, ok bool) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ve.Field, ve.Message, true
	}
	return "", "", false
}

// =============================================================================
// SECTION 6: defer -- CLEANUP THAT ALWAYS RUNS
// =============================================================================

// DeferOrder demonstrates that deferred calls run in LIFO order
// (Last In, First Out). The last `defer` registered runs first.
//
// Returns a string showing the order so the demo can show it cleanly
// without printing.
func DeferOrder() string {
	var b strings.Builder
	defer b.WriteString(" 1") // 3rd to run
	defer b.WriteString(" 2") // 2nd to run
	defer b.WriteString(" 3") // 1st to run

	b.WriteString("body")
	// Returning here triggers the deferred calls in REVERSE order.
	// We need to inspect b AFTER the defers, but b.String() called
	// inside the function captures only "body" so far. To work around
	// this, we use a named return + a final defer, demonstrated next.
	return b.String() // returns "body" -- defers haven't fired yet
}

// DeferOrderNamed shows the trick: use a named return so a deferred
// function can MODIFY the return value just before it's sent back to
// the caller. The defers do their work, and then the named-return
// `result` is what the caller sees.
//
// This is also how `defer` interacts with named returns to wrap errors
// (the technique we mentioned briefly in functions).
func DeferOrderNamed() (result string) {
	defer func() { result += " (defer 1 ran last)" }()
	defer func() { result += " (defer 2 ran second)" }()
	defer func() { result += " (defer 3 ran first)" }()
	result = "body"
	return
}

// DeferArgsCapture demonstrates that DEFERRED FUNCTION ARGUMENTS are
// evaluated WHEN THE DEFER IS REGISTERED, not when the function runs.
//
// We capture x=1 in the first defer's argument list. Even though x is
// later changed to 99, the deferred call was already given 1.
//
// In the second case, we wrap the call in a closure (no args at the
// outer call). The closure captures x by reference, so when it actually
// runs at function exit, it sees the latest value.
func DeferArgsCapture() (snapshot int, latest int) {
	x := 1

	// Captures x=1 at this line. The deferred Println later prints "1".
	defer func(captured int) {
		snapshot = captured
	}(x)

	// Closure: no args. x is captured by reference. Later this defer
	// will read the CURRENT value of x.
	defer func() {
		latest = x
	}()

	x = 99
	return // both defers run; named returns get filled in
}

// =============================================================================
// SECTION 7: panic AND recover
// =============================================================================

// MustDivide PANICS on division by zero instead of returning an error.
// Functions named "MustX" are a Go convention: they panic on failure
// rather than returning an error. Use them ONLY when failure indicates
// a programmer bug (the input should always be valid in correct code).
//
// Example from stdlib: regexp.MustCompile -- panics if the regex is
// invalid, because regex literals are usually compile-time constants
// and an invalid one is a programmer mistake.
func MustDivide(a, b int) int {
	if b == 0 {
		panic("MustDivide: division by zero")
	}
	return a / b
}

// SafeDivide demonstrates RECOVER. It calls MustDivide, but wraps the
// call in a deferred recover so a panic becomes a regular error
// returned to the caller.
//
// Mechanics:
//  1. We use a NAMED return `err` so the deferred function can set it.
//  2. The deferred function calls recover(). If a panic is in progress,
//     recover() returns the panic value (and stops the panic).
//     If no panic, recover() returns nil.
//  3. If we got a non-nil recover value, we convert it into an error
//     and assign it to err. The function then returns normally.
//
// This pattern is used at trust boundaries (e.g., HTTP handlers in
// long-running servers) so a single panic doesn't crash the whole process.
func SafeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			// r is whatever was passed to panic() -- could be a string,
			// an error, anything. We coerce to an error here.
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	result = MustDivide(a, b) // may panic; deferred recover catches it
	return result, nil
}

// =============================================================================
// SECTION 8: PUTTING IT TOGETHER -- A REGISTRATION FLOW
// =============================================================================
//
// A small example showing all the pieces working together: sentinel
// errors, custom error types, wrapping, and unwrapping.

// ErrUserExists is a sentinel for the "duplicate user" case.
var ErrUserExists = errors.New("user already exists")

// User is a simple user record.
type User struct {
	ID    int
	Email string
}

// UserStore is a tiny in-memory store of users.
type UserStore struct {
	users map[int]*User
}

// NewUserStore constructs an empty store.
func NewUserStore() *UserStore {
	return &UserStore{users: make(map[int]*User)}
}

// Register validates the input and adds a new user. Demonstrates:
//   - returning a *ValidationError directly when input is invalid
//   - returning a wrapped sentinel (ErrUserExists) when the user
//     already exists, with extra context about which user
//   - returning nil on success
func (s *UserStore) Register(id int, email string) error {
	// Custom error -- structured info for the caller to inspect.
	if err := ValidateEmail(email); err != nil {
		return err // already a *ValidationError
	}
	// Sentinel error -- with wrapping for context.
	if _, ok := s.users[id]; ok {
		return fmt.Errorf("Register id=%d: %w", id, ErrUserExists)
	}
	s.users[id] = &User{ID: id, Email: email}
	return nil
}

// HandleRegister is the kind of function you'd write in an HTTP handler:
// it calls Register and DECIDES how to respond based on the kind of
// error. This is where errors.Is and errors.As earn their keep.
func (s *UserStore) HandleRegister(id int, email string) string {
	err := s.Register(id, email)
	if err == nil {
		return fmt.Sprintf("OK: registered id=%d", id)
	}

	// Use errors.As to detect a structured validation error and pull
	// out its details for a 400-style response.
	var ve *ValidationError
	if errors.As(err, &ve) {
		return fmt.Sprintf("BadRequest: field=%s issue=%s", ve.Field, ve.Message)
	}

	// Use errors.Is to detect a specific sentinel for a 409-style response.
	if errors.Is(err, ErrUserExists) {
		return fmt.Sprintf("Conflict: %v", err)
	}

	// Catch-all for unexpected errors.
	return fmt.Sprintf("InternalError: %v", err)
}
