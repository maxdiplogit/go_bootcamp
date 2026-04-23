// Package funcsref is a personal reference package for functions in Go.
//
// CORE FACTS TO REMEMBER:
//
//  1. Go ALWAYS passes by value. Pointers let you "pass by reference" by
//     copying an address instead of the data itself.
//  2. Functions can return MULTIPLE values, idiomatically used for
//     (result, error) and (value, ok) patterns.
//  3. Functions are FIRST-CLASS values: assignable, passable, returnable,
//     storable in slices/maps/structs.
//  4. CLOSURES capture variables from their surrounding scope and keep
//     them alive after the surrounding function returns.
//  5. VARIADIC functions accept any number of arguments via `...T`.
//  6. NAMED RETURNS are pre-declared zero-value variables that a bare
//     `return` will send back; useful for short funcs and defer interaction.
//
// Each section below is a small, runnable example with explanatory comments.
// Use it as a quick lookup when you forget how something works.
package funcsref

import (
	"fmt"
	"strings"
)

// =============================================================================
// SECTION 1: BASICS -- PARAMETERS AND RETURN
// =============================================================================

// Add is the simplest possible function: takes two ints, returns their sum.
// Note the parameter SHORTHAND: when consecutive params share a type, you
// can list the type once at the end. `a, b int` == `a int, b int`.
func Add(a, b int) int {
	return a + b
}

// Greet has no return value -- the return type is omitted entirely (not
// `void`, no return type at all). Such functions are used for side effects
// like printing, writing to a file, modifying state via pointer, etc.
//
// In real production code you'd return an error to report failure, but
// for a simple greeter this is fine.
func Greet(name string) string {
	return "Hello, " + name + "!"
}

// NoArgsNoReturn shows the simplest signature -- nothing in, nothing out.
// Useful only for side effects. We return a string here just so the demo
// can show the output, but conceptually this models the void/void shape.
func NoArgsNoReturn() string {
	return "I take no args and (conceptually) return nothing."
}

// =============================================================================
// SECTION 2: CALL BY VALUE vs CALL BY POINTER
// =============================================================================

// FailedReset is a deliberate teaching example showing why call-by-value
// can surprise you. Even though we write `n = 0` inside the function, the
// caller's variable is untouched -- because `n` here is a COPY of the
// caller's value, not the same variable.
//
// We RETURN n at the end so you can see what the function thinks the
// value is (it sees 0), versus what the caller will see (still the
// original).
func FailedReset(n int) int {
	n = 0
	return n // function sees 0, but the original argument is unchanged
}

// SuccessfulReset uses a POINTER parameter. The caller passes &x, which
// is the address of x. Inside the function, `n` is a copy of that address,
// but both addresses point to the SAME underlying int. Writing through
// the pointer (*n = 0) modifies the original.
//
// Mental model: passing &x is like handing someone a paper with your home
// address on it. You both have separate copies of the paper, but they
// both point to the same house, and changes made to the house are visible
// to everyone who has the address.
func SuccessfulReset(n *int) {
	*n = 0 // dereference: write through the pointer to modify the target
}

// CounterStruct demonstrates the same value/pointer distinction with a
// struct. Structs are copied IN FULL when passed by value, which is both
// wasteful for large structs AND prevents mutation of the original.
type CounterStruct struct {
	Count int
}

// IncrementValue takes a struct by value -- the WHOLE struct is copied.
// Modifications inside this function affect only the copy. The caller's
// struct is unchanged.
func IncrementValue(c CounterStruct) {
	c.Count++
}

// IncrementPointer takes a pointer to the struct. Only 8 bytes are copied
// (the address), and modifications go through the pointer to the original.
//
// Note Go's auto-dereference: `c.Count++` is shorthand for `(*c).Count++`.
// The compiler inserts the dereference for you.
func IncrementPointer(c *CounterStruct) {
	c.Count++
}

// SliceMutation is the surprising case. Slices, maps, and channels are
// REFERENCE-LIKE in their behavior even when passed by value, because
// the value of a slice is a small header {pointer, length, capacity},
// and the underlying data is shared.
//
// So this function modifies the caller's slice contents, even though
// "Go always passes by value." Both rules are true at once -- the slice
// header was copied, but the data it refers to is shared.
//
// (This will fully click when you cover slices formally.)
func SliceMutation(s []int) {
	for i := range s {
		s[i] *= 2
	}
}

// =============================================================================
// SECTION 3: MULTIPLE RETURN VALUES
// =============================================================================

// Divide returns BOTH the quotient and remainder. Multiple return values
// are listed in parentheses. The caller unpacks them with multi-assign:
//
//	q, r := Divide(17, 5)
func Divide(a, b int) (int, int) {
	return a / b, a % b
}

// SafeDivide adds the famous Go (value, error) pattern. In real code you'd
// return a proper `error` type, but we haven't formally covered error
// handling yet, so we use a (value, ok bool) variant -- still the same
// shape and idea.
//
// The caller checks ok BEFORE using the value:
//
//	q, ok := SafeDivide(10, 0)
//	if !ok { ... handle ... }
func SafeDivide(a, b int) (int, bool) {
	if b == 0 {
		return 0, false // zero is just a placeholder; the bool tells you it's invalid
	}
	return a / b, true
}

// MinMax computes the smallest and largest values in a slice in one pass.
// Returning both at once is more efficient than calling separate Min and
// Max, and it's a more honest signature -- both values come from the
// same data.
func MinMax(nums []int) (int, int, bool) {
	if len(nums) == 0 {
		return 0, 0, false // can't compute min/max of nothing
	}
	min, max := nums[0], nums[0]
	for _, v := range nums[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max, true
}

// =============================================================================
// SECTION 4: NAMED RETURN VALUES
// =============================================================================

// SplitName demonstrates named returns. The signature declares
// `(first, last string)`, which:
//  1. Pre-declares `first` and `last` as local variables (initialized to "").
//  2. Documents the meaning of each return value (shows up in IDE tooltips).
//  3. Allows a "naked return" -- `return` with no arguments -- which sends
//     back the current values of the named returns.
//
// We could equivalently write `return first, last` explicitly. The naked
// return is a convenience, not a requirement.
func SplitName(full string) (first, last string) {
	parts := strings.SplitN(full, " ", 2)
	first = parts[0]
	if len(parts) > 1 {
		last = parts[1]
	}
	return // naked return -- sends current values of `first` and `last`
}

// ParseRange shows when named returns shine: short functions where the
// names ARE the documentation. The signature alone tells you exactly
// what each value means, with no comment needed.
func ParseRange(s string) (start, end int, ok bool) {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return // returns 0, 0, false (the zero values)
	}
	// Hand-parse small integers so we don't pull in strconv yet.
	// (We'll revisit this once you cover packages and strconv.)
	for _, c := range parts[0] {
		if c < '0' || c > '9' {
			return
		}
		start = start*10 + int(c-'0')
	}
	for _, c := range parts[1] {
		if c < '0' || c > '9' {
			return
		}
		end = end*10 + int(c-'0')
	}
	ok = true
	return
}

// =============================================================================
// SECTION 5: VARIADIC FUNCTIONS
// =============================================================================

// SumAll is variadic -- the `...int` parameter accepts ANY number of int
// arguments at the call site, including zero. Inside the function,
// `nums` is just a regular []int.
//
//	SumAll()              -> 0   (empty slice)
//	SumAll(1)             -> 1
//	SumAll(1, 2, 3)       -> 6
//	SumAll(1, 2, 3, 4, 5) -> 15
//
// To pass an EXISTING slice, use the spread syntax at the call site:
//
//	nums := []int{1, 2, 3}
//	SumAll(nums...)       -> 6
//
// Without the `...`, Go would think you're passing one []int argument.
func SumAll(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// LogPrefixed shows mixing a fixed parameter with a variadic one. The
// variadic must be LAST -- you can have any number of fixed parameters
// before it, but the `...` parameter must come at the end of the list.
//
// We return the formatted string here (rather than printing) so the demo
// can verify the output.
func LogPrefixed(level string, parts ...string) string {
	return fmt.Sprintf("[%s] %s", level, strings.Join(parts, " "))
}

// =============================================================================
// SECTION 6: FUNCTIONS AS VALUES (FIRST-CLASS)
// =============================================================================

// IntOp is a TYPE ALIAS for `func(int, int) int`. Aliasing function types
// makes signatures more readable when you pass them around or store them
// in a struct. (You'll see this everywhere in idiomatic Go.)
type IntOp func(int, int) int

// Apply takes an IntOp value and calls it. Notice that we treat `op`
// like any other variable -- this is what "first-class function" means:
// functions are values you can pass around like ints or strings.
func Apply(a, b int, op IntOp) int {
	return op(a, b)
}

// AddOp and SubOp are FUNCTION VALUES of type IntOp. They're declared
// at the package level just like any other variable, and assigned a
// function literal as their value.
//
// (This is purely for demonstration; in practice you'd often just write
// the function literal inline at the call site.)
var AddOp IntOp = func(a, b int) int { return a + b }
var SubOp IntOp = func(a, b int) int { return a - b }

// OpsByName returns a map from operator names to operations. This is the
// classic "dispatch table" pattern -- a clean alternative to long
// if/else or switch chains for selecting behavior by a key.
func OpsByName() map[string]IntOp {
	return map[string]IntOp{
		"+": func(a, b int) int { return a + b },
		"-": func(a, b int) int { return a - b },
		"*": func(a, b int) int { return a * b },
		"/": func(a, b int) int {
			if b == 0 {
				return 0 // sentinel; real code would return an error
			}
			return a / b
		},
	}
}

// =============================================================================
// SECTION 7: ANONYMOUS FUNCTIONS
// =============================================================================

// AnonymousImmediate shows the IIFE pattern -- declare a function and
// immediately invoke it. The trailing `(5, 7)` is the function call.
//
// This is occasionally useful for scoping a small computation that you
// don't want to leak helper variables into the surrounding scope.
func AnonymousImmediate() int {
	return func(a, b int) int {
		return a*a + b*b
	}(5, 7) // immediately call with (5, 7)
}

// FilterStrings takes a slice and an inline predicate function. The
// predicate is an anonymous function passed at the call site --
// FilterStrings doesn't care if it's anonymous or named.
//
// This is the foundation of higher-order utilities: Filter, Map, Reduce,
// Sort.Slice, etc. all take functions as parameters.
func FilterStrings(in []string, keep func(string) bool) []string {
	out := make([]string, 0, len(in))
	for _, s := range in {
		if keep(s) {
			out = append(out, s)
		}
	}
	return out
}

// =============================================================================
// SECTION 8: CLOSURES
// =============================================================================

// MakeCounter returns a closure: a function that captures (closes over)
// the local variable `count`. Even after MakeCounter returns, the inner
// function STILL HAS ACCESS to its `count` -- it's kept alive by the
// runtime for as long as the returned function exists.
//
// Each call to MakeCounter creates a FRESH `count` (a new local variable),
// so different counters don't share state.
func MakeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// MakeMultiplier shows the "configuration captured at creation time"
// closure pattern. The factor is fixed when MakeMultiplier is called;
// the returned function applies that factor on every call.
//
// This lets you customize behavior without using a struct or class --
// the captured variable IS your "state."
func MakeMultiplier(factor int) func(int) int {
	return func(n int) int {
		return n * factor
	}
}

// MakeAccount returns multiple closures that share the SAME captured
// state. `bal` lives in MakeAccount's scope; both `deposit` and `balance`
// close over it, so they see and modify the same variable.
//
// This is one way to model encapsulated, mutable state in Go without
// using a struct -- though for anything non-trivial, a struct with
// methods is usually clearer. This pattern is mostly useful for small
// callback-driven cases.
func MakeAccount(initial float64) (deposit func(float64), balance func() float64) {
	bal := initial
	deposit = func(amount float64) { bal += amount }
	balance = func() float64 { return bal }
	return
}

// =============================================================================
// SECTION 9: RECURSION
// =============================================================================

// Factorial is the classic recursion example. The function calls itself
// with a smaller input until it hits the base case (n <= 1).
//
// Go does NOT optimize tail calls, so deep recursion can grow the stack.
// For factorial of small n, this is perfectly fine; for problems that
// recurse millions of levels deep, prefer an iterative version.
func Factorial(n int) int {
	if n <= 1 {
		return 1 // base case -- the recursion stops here
	}
	return n * Factorial(n-1)
}

// Fibonacci shows recursion with TWO recursive calls. Note that the naive
// version is exponentially slow because it recomputes the same values
// over and over -- Fib(5) calls Fib(3) twice, Fib(2) three times, etc.
//
// In real code, you'd add memoization or use iteration. We keep the
// naive version here because it's the clearest demonstration of recursion
// branching, and it's still fast enough for small n.
func Fibonacci(n int) int {
	if n < 2 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// =============================================================================
// SECTION 10: A REALISTIC EXAMPLE -- A TINY CALCULATOR
// =============================================================================

// Calculator is a small example that pulls in MOST of the topics from
// this package:
//   - First-class functions (operations stored in a map)
//   - Variadic input (any number of operations to apply)
//   - Closures (preset operations capture a "factor")
//   - Multiple return values (PreloadedCalculator returns two things)
//   - Pointer receivers (we mutate History in place)
//   - Named returns (in Summary)
//
// The calculator stores a running total and a history of operations.
type Calculator struct {
	Total   int
	History []string
}

// NewCalculator is the conventional constructor pattern: a function named
// New<Type> that returns a pointer to a zero-valued (or initially
// configured) instance.
//
// Why return a pointer? Because callers will mutate the calculator via
// methods, and we want those mutations to stick. If we returned a value,
// every method call would operate on a copy.
func NewCalculator(initial int) *Calculator {
	return &Calculator{
		Total:   initial,
		History: []string{fmt.Sprintf("init=%d", initial)},
	}
}

// Apply runs an operation against the current total. The operation is
// passed in as a function value -- so callers can supply ANY operation,
// including custom closures, without us having to enumerate them.
//
// `name` is just for the history log so we can record what happened.
func (c *Calculator) Apply(name string, op func(int) int) {
	before := c.Total
	c.Total = op(c.Total)
	c.History = append(c.History,
		fmt.Sprintf("%s: %d -> %d", name, before, c.Total))
}

// NamedOp pairs a name with an operation, used by ApplyMany.
type NamedOp struct {
	Name string
	Op   func(int) int
}

// ApplyMany takes ANY number of NamedOp values via variadic. The variadic
// form makes it easy to chain operations without wrapping them in a
// slice at the call site.
func (c *Calculator) ApplyMany(ops ...NamedOp) {
	for _, no := range ops {
		c.Apply(no.Name, no.Op)
	}
}

// Summary returns a short human-readable snapshot of the calculator's
// state. Uses a NAMED RETURN to demonstrate that style for short
// reporting functions.
func (c *Calculator) Summary() (s string) {
	s = fmt.Sprintf("Total=%d, %d operations", c.Total, len(c.History)-1)
	return
}

// PreloadedCalculator returns a calculator together with a small set of
// "pre-built" named operations as closures. Demonstrates returning
// multiple values, closures over a shared `factor`, and the dispatch-
// table pattern in one go.
func PreloadedCalculator(initial, factor int) (*Calculator, map[string]func(int) int) {
	c := NewCalculator(initial)
	ops := map[string]func(int) int{
		"add":      func(x int) int { return x + factor }, // closes over `factor`
		"subtract": func(x int) int { return x - factor },
		"multiply": func(x int) int { return x * factor },
		"square":   func(x int) int { return x * x }, // doesn't use factor
	}
	return c, ops
}
