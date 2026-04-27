// Package genericsref is a personal reference package for generics in Go.
//
// CORE MENTAL MODEL:
//
// Generics let you write a function or type ONCE, with a placeholder type
// (called a "type parameter") that the compiler fills in for you when you
// call the function. You get one source-of-truth implementation, full
// type safety, and no interface boxing.
//
// FIVE THINGS TO REMEMBER:
//
//  1. Type parameters go in [SQUARE BRACKETS] before the regular params:
//     func Max[T any](a, b T) T  -- T is the placeholder.
//  2. Each type parameter has a CONSTRAINT that limits which types are
//     allowed. `any` means "anything"; `comparable` means "supports ==";
//     a union like `int | float64` means "one of these"; an interface
//     name means "satisfies this interface".
//  3. The constraint also determines what OPERATIONS are allowed on T
//     inside the function. Looser constraint = fewer operations possible.
//  4. TYPE INFERENCE often lets you skip the [T] when calling: Max(1, 2)
//     instead of Max[int](1, 2). You explicitly specify it only when
//     inference fails.
//  5. TYPES can be generic too, not just functions. Stack[T any] is a
//     generic container; Stack[int] and Stack[string] are different types.
//
// USE WHEN: same logic, multiple types, want to preserve type info.
// DON'T USE WHEN: a regular function or an interface fits cleanly.
package generics

import (
	"cmp"
	"fmt"
)

// =============================================================================
// SECTION 1: A FIRST GENERIC FUNCTION -- Max
// =============================================================================

// Max returns the larger of two values. The constraint `cmp.Ordered`
// (from the standard library `cmp` package, Go 1.21+) covers any type
// where < and > are defined: ints, floats, and strings.
//
// Without generics, you'd need MaxInt, MaxFloat, MaxString -- three
// near-identical functions. With generics, one definition works for all.
//
// Read the signature piece by piece:
//
//	func Max         <- name
//	[T cmp.Ordered]  <- type parameter T, constrained to ordered types
//	(a, b T)         <- two parameters of type T
//	T                <- returns a T
func Max[T cmp.Ordered](a, b T) T {
	if a > b { // `>` is allowed because cmp.Ordered guarantees it
		return a
	}
	return b
}

// Min is the mirror image. Same constraint, same shape.
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// =============================================================================
// SECTION 2: A SLIGHTLY BIGGER EXAMPLE -- Sum
// =============================================================================

// Number is a USER-DEFINED CONSTRAINT. It's an interface that lists the
// types we want to accept. Inside this constraint we list every numeric
// type that supports `+`. The `~` prefix is a small extension: it means
// "this type OR any type defined as `type X T` from it" (so MyInt, where
// `type MyInt int`, also qualifies). For most beginner uses you can
// just ignore the `~`.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Sum adds up all elements of a slice. Works for any numeric type.
// The expression `var total T` gives us a fresh zero value of T (0 for
// numeric types). This is how you create "the zero of T" without knowing
// what T actually is.
func Sum[T Number](nums []T) T {
	var total T
	for _, v := range nums {
		total += v // `+=` is allowed: every type in Number supports it
	}
	return total
}

// =============================================================================
// SECTION 3: THE `comparable` CONSTRAINT -- Contains
// =============================================================================

// Contains reports whether `target` appears in `items`. The constraint
// `comparable` is BUILT IN to Go and means "supports == and !=". It
// covers int, string, bool, structs of comparable fields, pointers, and
// channels. It excludes slices, maps, and functions (those don't support ==).
//
// Try removing `comparable` and using `any` instead. The compiler will
// reject `v == target` with: "cannot compare v == target (operator ==
// not defined on T)". That's the constraint doing its job: it tells the
// compiler what operations are safe inside the function body.
func Contains[T comparable](items []T, target T) bool {
	for _, v := range items {
		if v == target {
			return true
		}
	}
	return false
}

// IndexOf is the same idea, returning the position (or -1).
func IndexOf[T comparable](items []T, target T) int {
	for i, v := range items {
		if v == target {
			return i
		}
	}
	return -1
}

// =============================================================================
// SECTION 4: TYPE INFERENCE -- WHEN YOU CAN OMIT [T]
// =============================================================================

// Pair is generic over TWO type parameters: K (the key type) and V (the
// value type). Both are unconstrained (`any`), so they can be anything.
// The function returns a string built from both.
//
// `K, V any` is shorthand for `K any, V any` -- same as parameter
// declarations.
func Pair[K, V any](k K, v V) string {
	return fmt.Sprintf("%v=%v", k, v)
}

// CallExamples demonstrates the call-side behavior of inference. These
// are equivalent at runtime; only the call syntax differs. We wrap each
// result with fmt.Sprintf so the slice can hold them all as strings for
// display.
func CallExamples() []string {
	return []string{
		Pair("count", 42),                                               // K, V inferred as string, int
		Pair[string, int]("count", 42),                                  // explicit -- rare, only if inference fails
		fmt.Sprintf("Max(...) = %v", Max(3.14, 2.71)),                   // T inferred as float64
		fmt.Sprintf("Max[float64](...) = %v", Max[float64](3.14, 2.71)), // explicit version
	}
}

// =============================================================================
// SECTION 5: GENERIC HIGHER-ORDER FUNCTIONS -- Map, Filter
// =============================================================================
//
// These are the canonical "you'll write them once and use them forever"
// utilities. The standard library has `slices.Map` and `slices.Filter`
// equivalents now (Go 1.21+); we redefine them here to show the pattern.

// MapSlice transforms each element of a slice using the supplied function.
// TWO type parameters: T is the input type, U is the output type. They
// can be the same (e.g. doubling ints) or different (e.g. converting
// ints to strings).
//
// Notice the function parameter `f func(T) U` is just a regular Go
// function value -- the only generic part is that its types are spelled
// in terms of T and U. That's all generics do here: parameterize types.
func MapSlice[T, U any](in []T, f func(T) U) []U {
	out := make([]U, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}

// FilterSlice keeps only the elements for which `keep(v)` is true. Single
// type parameter T -- input and output have the same element type.
func FilterSlice[T any](in []T, keep func(T) bool) []T {
	out := make([]T, 0, len(in))
	for _, v := range in {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

// Reduce collapses a slice into a single value using an accumulator
// function. Like MapSlice, has TWO type parameters: T for input
// elements, U for the accumulator/result.
func Reduce[T, U any](in []T, initial U, f func(U, T) U) U {
	acc := initial
	for _, v := range in {
		acc = f(acc, v)
	}
	return acc
}

// =============================================================================
// SECTION 6: A GENERIC TYPE -- Stack[T]
// =============================================================================

// Stack is a generic LIFO container. The `[T any]` after the type name
// declares the type parameter. From this line on, T can be used inside
// the struct definition wherever you'd normally use a concrete type.
type Stack[T any] struct {
	items []T
}

// NewStack is the conventional constructor. The type parameter list
// `[T any]` appears before the parameters of the function, just as it
// does for any generic function. The return type uses T.
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{items: []T{}}
}

// Push adds a value. Note the receiver: `*Stack[T]`, NOT `*Stack`. When
// defining a method on a generic type, you keep the parameter in the
// receiver type. You don't repeat the constraint -- T is already known
// from the type definition.
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop removes and returns the top value. The (T, bool) return is the
// classic "value, ok" pattern. When the stack is empty, we return a
// zero value of T -- created with `var zero T`. This is the standard
// way to get "the zero of T" inside a generic function.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	last := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return last, true
}

// Len returns the current size. Trivial -- shown so you see a method
// that doesn't use T in its signature, and that's also fine.
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// Peek returns the top value without removing it.
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// =============================================================================
// SECTION 7: A GENERIC TYPE WITH TWO PARAMETERS -- Pair, KeyValueStore
// =============================================================================

// KVPair is a tiny generic struct holding two values of (potentially)
// different types. The type parameter list `[K, V any]` appears right
// after the type name, separated by commas.
type KVPair[K, V any] struct {
	Key   K
	Value V
}

// String makes KVPair satisfy fmt.Stringer no matter what K and V are.
// This works because `%v` works for any type.
func (p KVPair[K, V]) String() string {
	return fmt.Sprintf("%v=%v", p.Key, p.Value)
}

// KeyValueStore is a tiny generic key/value container. The constraint
// on K is `comparable` because we use K as a map key (Go map keys must
// support ==). V is unconstrained.
type KeyValueStore[K comparable, V any] struct {
	data map[K]V
}

// NewKVStore is the constructor.
func NewKVStore[K comparable, V any]() *KeyValueStore[K, V] {
	return &KeyValueStore[K, V]{data: make(map[K]V)}
}

func (s *KeyValueStore[K, V]) Set(k K, v V) {
	s.data[k] = v
}

func (s *KeyValueStore[K, V]) Get(k K) (V, bool) {
	v, ok := s.data[k]
	return v, ok
}

func (s *KeyValueStore[K, V]) Len() int {
	return len(s.data)
}

// =============================================================================
// SECTION 8: A GENERIC INTERFACE -- describing methods over a type
// =============================================================================

// A generic interface is just an interface that uses type parameters.
// The methods are described in terms of T.
//
// Container[T] says: "any type that has Add(T), Get(int)(T,bool), and Size()".
// This is sometimes useful when you want to write functions that take
// "any container of T", regardless of whether the underlying type is a
// stack, queue, list, etc.
type Container[T any] interface {
	Add(v T)
	Get(i int) (T, bool)
	Size() int
}

// SimpleList is a concrete generic type that satisfies Container[T] for
// any T. We add it just to show that the pieces fit together.
type SimpleList[T any] struct {
	items []T
}

func NewSimpleList[T any]() *SimpleList[T] {
	return &SimpleList[T]{}
}

func (l *SimpleList[T]) Add(v T) {
	l.items = append(l.items, v)
}

func (l *SimpleList[T]) Get(i int) (T, bool) {
	if i < 0 || i >= len(l.items) {
		var zero T
		return zero, false
	}
	return l.items[i], true
}

func (l *SimpleList[T]) Size() int {
	return len(l.items)
}

// SumContainer is a function that takes ANY Container[T] holding numeric
// items and sums them. Notice the two type parameters: T for the
// element type, and the container is constrained to Container[T].
//
// This shows generics + interfaces working together: T is generic, but
// we constrain it to Number so we can use `+`, AND we accept any thing
// that satisfies Container[T] (could be SimpleList, could be your own
// custom container later).
func SumContainer[T Number](c Container[T]) T {
	var total T
	for i := 0; i < c.Size(); i++ {
		v, _ := c.Get(i)
		total += v
	}
	return total
}
