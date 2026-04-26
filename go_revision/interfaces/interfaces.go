// Package interfacesref is a personal reference package for interfaces in Go.
//
// CORE MENTAL MODEL:
//
// An interface is a description of CAPABILITY. It says "any value that has
// these methods is acceptable here." Types don't declare that they
// implement an interface -- the compiler just checks whether their methods
// match. This is called STRUCTURAL TYPING.
//
// Internally, an interface value is a 2-word structure:
//
//	(type_descriptor_pointer, data_pointer)
//
// Method calls on an interface look up the method via the type descriptor
// and call it with the data pointer as the receiver. That's the whole
// mechanism.
//
// FIVE THINGS TO REMEMBER:
//
//  1. BASICS: define methods on a type; that type satisfies any interface
//     whose methods it has. No `implements` keyword.
//  2. EMPTY INTERFACE: `any` (= interface{}) is satisfied by every type.
//     Useful for "I don't know the type yet" cases. Use sparingly.
//  3. EMBEDDING: one interface can embed others to compose larger ones.
//     io.ReadWriter is just Reader + Writer.
//  4. TYPE ASSERTION: `v.(T)` recovers a concrete type from an interface.
//     Always use the two-value form: `t, ok := v.(T)`.
//  5. TYPE SWITCH: a `switch v.(type)` block lets you handle many possible
//     concrete types cleanly.
//
// Each section is small and runnable. Read top to bottom or jump around.
package interfaces

import (
	"fmt"
	"strings"
)

// =============================================================================
// SECTION 1: BASICS -- DEFINING AND SATISFYING AN INTERFACE
// =============================================================================

// Greeter is a single-method interface. ANY type with a `Greet() string`
// method satisfies Greeter automatically -- no declaration needed.
type Greeter interface {
	Greet() string
}

// EnglishSpeaker satisfies Greeter implicitly: it has a `Greet() string`
// method, so the compiler treats it as a Greeter wherever one is expected.
//
// Notice: EnglishSpeaker doesn't say "implements Greeter" anywhere. Go
// figures it out structurally, by looking at the methods.
type EnglishSpeaker struct {
	Name string
}

func (e EnglishSpeaker) Greet() string {
	return "Hello, " + e.Name
}

// SpanishSpeaker is a completely different type that ALSO happens to have
// a Greet() string method, so it ALSO satisfies Greeter.
type SpanishSpeaker struct {
	Name string
}

func (s SpanishSpeaker) Greet() string {
	return "Hola, " + s.Name
}

// SayHi accepts ANY Greeter. The function doesn't care whether you pass
// an EnglishSpeaker, a SpanishSpeaker, or some type that doesn't even
// exist yet. As long as it has a Greet() string method, it works.
//
// This is the entire point of interfaces: write functions in terms of
// behavior, not concrete types, and they instantly work with anything
// that has the right methods.
func SayHi(g Greeter) string {
	return g.Greet()
}

// =============================================================================
// SECTION 2: WHY INTERFACES MATTER -- ONE FUNCTION, MANY TYPES
// =============================================================================

// Shape is a tiny interface describing "anything with an area." Three
// completely unrelated structs can satisfy it just by having Area().
type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 { return r.Width * r.Height }

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 { return 3.14159 * c.Radius * c.Radius }

type Triangle struct {
	Base, Height float64
}

func (t Triangle) Area() float64 { return 0.5 * t.Base * t.Height }

// TotalArea works for ANY mix of Shapes. We can pass a slice containing
// rectangles, circles, and triangles all together -- something we couldn't
// do if the function required one specific concrete type.
//
// This is the polymorphism payoff: write the function once, plug in any
// types that satisfy the interface. New shape next month? It'll work
// without changing this function.
func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// =============================================================================
// SECTION 3: THE EMPTY INTERFACE -- `any`
// =============================================================================

// AnyDescription accepts a value of any type at all. `any` (= interface{})
// has zero methods, so every type trivially satisfies it.
//
// The cost: inside the function, we know nothing about the value. To do
// anything with it, we have to inspect the type at runtime via assertion
// or type switch (see next section).
//
// Use sparingly. With generics now in Go 1.18+, `any` for "I want
// multiple types" is usually the wrong choice -- generics give you type
// safety that `any` throws away.
func AnyDescription(v any) string {
	return fmt.Sprintf("got value=%v of type=%T", v, v)
	// %T is a format verb that prints the dynamic type of an interface value.
	// It works because interface values carry their type info around with them.
}

// =============================================================================
// SECTION 4: TYPE ASSERTIONS
// =============================================================================

// AssertString tries to recover a string from an `any`. Always use the
// TWO-VALUE form `s, ok := v.(string)` -- it sets ok=false on mismatch
// instead of panicking.
//
// The single-value form `s := v.(string)` panics on mismatch. You will
// rarely want that.
func AssertString(v any) (string, bool) {
	s, ok := v.(string)
	return s, ok
}

// AssertInt is the same idea for int.
func AssertInt(v any) (int, bool) {
	n, ok := v.(int)
	return n, ok
}

// AssertGreeter shows that you can also assert for INTERFACE types, not
// just concrete types. Here we ask "does this `any` value satisfy the
// Greeter interface?"
//
// This is how you discover capabilities at runtime: "I have some value;
// can it Greet?"
func AssertGreeter(v any) (Greeter, bool) {
	g, ok := v.(Greeter)
	return g, ok
}

// =============================================================================
// SECTION 5: TYPE SWITCH
// =============================================================================

// Describe handles many possible concrete types cleanly. The special
// syntax is `switch x := v.(type) { ... }` -- inside each case, `x` is
// already typed correctly, so no further assertion is needed.
//
// Use a type switch (instead of repeated `v.(T)` assertions) whenever
// you have an interface value that could be one of several types.
func Describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil value"
	case int:
		return fmt.Sprintf("int with double = %d", x*2)
	case string:
		return fmt.Sprintf("string of length %d", len(x))
	case bool:
		return fmt.Sprintf("bool = %v", x)
	case Greeter:
		// Note: this case is reached only if NONE of the above matched.
		// Cases are tested in order, top to bottom.
		return "greeter says: " + x.Greet()
	default:
		return fmt.Sprintf("unhandled type %T", x)
	}
}

// =============================================================================
// SECTION 6: EMBEDDING INTERFACES
// =============================================================================

// Reader and Writer are tiny single-method interfaces. (These echo the
// real `io.Reader` and `io.Writer` -- the same pattern.)
type Reader interface {
	Read() string
}

type Writer interface {
	Write(s string)
}

// ReadWriter is built by EMBEDDING Reader and Writer. Anything that
// satisfies both Reader AND Writer automatically satisfies ReadWriter.
//
// This is the only common use of interface embedding: composing small
// interfaces into larger ones. Don't overthink it.
type ReadWriter interface {
	Reader
	Writer
}

// MemoryStore is a tiny in-memory thing that has both Read and Write.
// Because it has both methods, it ALSO satisfies ReadWriter automatically.
type MemoryStore struct {
	Data string
}

func (m *MemoryStore) Read() string   { return m.Data }
func (m *MemoryStore) Write(s string) { m.Data = s }

// UseReadWriter takes a ReadWriter -- something that can both read and
// write. It writes to the store and reads back what it stored.
func UseReadWriter(rw ReadWriter) string {
	rw.Write("hello, embedded interfaces")
	return rw.Read()
}

// =============================================================================
// SECTION 7: POINTER RECEIVERS AND INTERFACE SATISFACTION
// =============================================================================

// Counter is an interface with a mutating method.
type Counter interface {
	Increment()
	Value() int
}

// IntCounter has Increment with a POINTER receiver because Increment must
// modify the counter's state. Value uses a value receiver because reading
// doesn't need to mutate.
//
// Convention: when ANY method on a type uses a pointer receiver, make
// ALL of them use pointer receivers, for consistency. We don't follow
// that convention here because we want to show you the rule it covers.
type IntCounter struct {
	N int
}

func (c *IntCounter) Increment() { c.N++ }
func (c IntCounter) Value() int  { return c.N }

// NewIntCounter is the conventional constructor. It returns *IntCounter
// because the caller will (almost certainly) need to call Increment,
// which requires a pointer.
//
// THE KEY POINT for interfaces:
//   - *IntCounter satisfies Counter (it has both methods in its method set)
//   - IntCounter (value) does NOT satisfy Counter (Increment requires *)
//
// So you must pass &IntCounter{} or NewIntCounter() to a Counter param,
// not IntCounter{}.
func NewIntCounter() *IntCounter {
	return &IntCounter{}
}

// IncrementMany takes a Counter (interface). The compiler checks at the
// call site that whatever you pass satisfies Counter. If you forget to
// pass a pointer when needed, you'll get a clear compile error.
func IncrementMany(c Counter, times int) int {
	for i := 0; i < times; i++ {
		c.Increment()
	}
	return c.Value()
}

// =============================================================================
// SECTION 8: A SMALL REALISTIC EXAMPLE -- A LOG SINK
// =============================================================================
//
// This pulls together the basics in one place: a single `Log` function
// that works against any log destination via an interface.

// LogSink is an interface for "anything I can write a log line to."
// One method, simple and focused -- the Go style.
type LogSink interface {
	Log(line string)
}

// ConsoleSink prints lines to standard output. Real code would use
// fmt.Println; we use a string builder so the demo can capture the
// output.
type ConsoleSink struct {
	Buffer *strings.Builder
}

func (c *ConsoleSink) Log(line string) {
	c.Buffer.WriteString("[console] " + line + "\n")
}

// MemorySink stores lines in a slice. Useful for tests -- you can write
// log lines and then assert they appeared.
type MemorySink struct {
	Lines []string
}

func (m *MemorySink) Log(line string) {
	m.Lines = append(m.Lines, line)
}

// PrefixSink is a SINK THAT WRAPS ANOTHER SINK. It adds a prefix and
// forwards to the inner sink. This is the "decorator" pattern --
// possible cheaply because PrefixSink itself satisfies LogSink, so it
// can stand in anywhere a LogSink is expected.
type PrefixSink struct {
	Prefix string
	Inner  LogSink
}

func (p *PrefixSink) Log(line string) {
	p.Inner.Log(p.Prefix + ": " + line)
}

// LogAll writes the same lines to a sink. It works for ANY LogSink --
// console, memory, prefixed, or some sink you write tomorrow.
//
// This is the function we wanted from the start: write one function,
// plug in any destination.
func LogAll(sink LogSink, lines []string) {
	for _, line := range lines {
		sink.Log(line)
	}
}
