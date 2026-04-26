// Package pointersref is a personal reference package for pointers and
// memory management in Go.
//
// CORE FACTS TO REMEMBER:
//
//  1. A POINTER is a variable whose value is the memory address of
//     another variable. Its type is `*T` where T is the pointed-to type.
//  2. `&x` takes the ADDRESS of x (gives you a pointer).
//     `*p` DEREFERENCES p (gives you the value at that address).
//     Same symbol `*`, two uses: TYPE modifier vs DEREFERENCE operator.
//  3. The zero value of any pointer type is `nil`. Dereferencing a nil
//     pointer is a runtime PANIC, not a compile error.
//  4. Slices, maps, channels, and strings are ALREADY reference-like --
//     you almost never wrap them in pointers.
//  5. Go does ESCAPE ANALYSIS at compile time to decide whether a value
//     lives on the stack or the heap. You don't manage this yourself.
//  6. Go has a CONCURRENT GARBAGE COLLECTOR. Heap objects are reclaimed
//     automatically when they become unreachable from roots (globals +
//     goroutine stacks).
//
// Each section below is a small, runnable example with explanatory
// comments. Read top-to-bottom or jump to the topic you need to refresh.
package pointersref

import (
	"fmt"
	"runtime"
)

// =============================================================================
// SECTION 1: POINTER BASICS -- `&` AND `*`
// =============================================================================

// ShowAddress returns the address of an int variable. This is pure
// illustration -- we never do anything useful with the printed address,
// but it proves that variables live somewhere concrete in memory.
//
// Note: you can't rely on the actual address staying the same across runs,
// platforms, or even within a single program due to stack growth.
func ShowAddress() string {
	x := 42
	p := &x // p is of type *int -- "pointer to int"
	// %p formats a pointer as a hex address like 0xc00001c030.
	return fmt.Sprintf("x=%d, address of x (=&x) = %p", x, p)
}

// ReadThrough demonstrates the `*p` dereference operator in read mode.
// Given a pointer, `*p` gets you the value at the address.
//
// Go AUTO-DEREFERENCES when you use `.` on a pointer to a struct, so
// you won't often write `*p` for struct fields -- but for primitives
// you need it explicitly.
func ReadThrough() (int, int) {
	x := 42
	p := &x
	return x, *p // both are 42 -- one direct, one through the pointer
}

// WriteThrough demonstrates `*p = ...` in write mode -- modifying the
// original value THROUGH the pointer. This is the whole point of
// pointers: shared, mutable access.
func WriteThrough() int {
	x := 42
	p := &x
	*p = 100 // write through the pointer -- x itself changes
	return x // 100, not 42
}

// Aliasing shows that two pointer variables can point to the same
// underlying data. Modifying through one is visible through the other.
// This is SHARING, and it's exactly what pointers are for.
func Aliasing() (int, int) {
	x := 42
	p := &x
	q := p // q and p hold the same address

	*q = 999 // modify through q

	return *p, *q // both 999 -- same underlying int
}

// =============================================================================
// SECTION 2: NIL POINTERS
// =============================================================================

// ZeroValue shows that the zero value of a pointer type is `nil`.
// Any pointer you declare without initializing starts out nil, and
// dereferencing it would PANIC at runtime.
func ZeroValue() (pointerIsNil bool) {
	var p *int
	return p == nil // true -- uninitialized pointers are nil
}

// SafeDeref demonstrates the defensive check for nil before dereferencing.
// If p is nil, we return a sentinel (0, false). Otherwise we dereference
// and return (value, true). The caller inspects `ok` before using `val`.
//
// This `(value, ok)` pattern should feel familiar from the functions
// lesson -- it's Go's typical way to signal "maybe has a value."
func SafeDeref(p *int) (val int, ok bool) {
	if p == nil {
		return 0, false
	}
	return *p, true
}

// =============================================================================
// SECTION 3: CREATING POINTERS -- &, new, and constructors
// =============================================================================

// WithAddress is the most common way to create a pointer: take the
// address of an existing local variable with `&`.
func WithAddress() *int {
	x := 42
	return &x // x escapes to the heap; see Section 8
}

// WithNew shows `new(T)`, which allocates a zero-valued T and returns
// a *T. Used less often than & in practice -- idiomatic Go prefers
// `&T{...}` for structs because you usually want to initialize fields.
func WithNew() *int {
	p := new(int) // p is *int, pointing to a fresh zero int
	*p = 42
	return p
}

// WithStructLiteral is the most common construction form you'll see in
// real Go code: `&Type{...}` to allocate and initialize in one step.
type Config struct {
	Host string
	Port int
}

// NewConfig is a typical constructor function: named `New<Type>`,
// returns a pointer, sets sensible defaults for fields the caller
// didn't specify.
func NewConfig(host string) *Config {
	return &Config{
		Host: host,
		Port: 8080, // default
	}
}

// =============================================================================
// SECTION 4: POINTERS WITH STRUCTS
// =============================================================================

// Point is a tiny struct used to demonstrate value vs pointer behavior.
type Point struct {
	X, Y int
}

// MoveValue modifies a COPY of the Point. The caller sees no change.
// This is call-by-value for a struct: the entire struct is copied into
// the parameter when the function is invoked.
func MoveValue(p Point, dx, dy int) Point {
	p.X += dx
	p.Y += dy
	return p // we have to return the modified copy -- otherwise the
	// caller's Point is untouched
}

// MovePointer modifies the ORIGINAL point via a pointer. The caller
// sees the change directly, no return value needed.
//
// Note the auto-dereference: `p.X` is shorthand for `(*p).X`. Go inserts
// the dereference for you when you use `.` on a pointer to a struct.
func MovePointer(p *Point, dx, dy int) {
	p.X += dx
	p.Y += dy
}

// ShareStruct shows that two pointers can refer to the same struct.
// Changes via one are visible via the other -- SHARED state.
func ShareStruct() (first, second Point) {
	p := &Point{X: 1, Y: 2}
	q := p // q is another pointer to the SAME Point, not a copy

	q.X = 999 // modify via q

	return *p, *q // both show X=999, Y=2
}

// ReturnLocalPointer is safe in Go -- and SHOULD BE UNFAMILIAR if you're
// coming from C. In C, returning the address of a local variable is a
// bug: the local goes out of scope and the pointer dangles. In Go, the
// compiler detects that `p` escapes and allocates it on the heap instead
// of the stack, so the pointer remains valid.
//
// You never write malloc/free in Go. Just return addresses freely; the
// compiler and GC handle the rest.
func ReturnLocalPointer() *Point {
	p := Point{X: 10, Y: 20}
	return &p // `p` is moved to the heap because its address escapes
}

// =============================================================================
// SECTION 5: POINTERS WITH SLICES
// =============================================================================

// DoubleSliceContents mutates the elements of the slice. Notice we pass
// []int, NOT *[]int -- the slice header is copied, but its internal
// pointer still points at the same backing array, so mutations are
// visible to the caller.
func DoubleSliceContents(s []int) {
	for i := range s {
		s[i] *= 2
	}
}

// AppendWithoutPointer is a TEACHING example of a pitfall. `append` may
// reallocate when capacity runs out; the new header is only visible
// inside this function because the caller passed a HEADER by value.
//
// The result: the caller's slice is NOT changed. Run this to see.
func AppendWithoutPointer(s []int, v int) {
	s = append(s, v) // may or may not reallocate -- either way, caller won't see it
}

// AppendReturning is the IDIOMATIC Go fix: return the new slice. The
// caller reassigns: `nums = AppendReturning(nums, 4)`. This is how
// the built-in `append` itself works, and it's the standard pattern.
func AppendReturning(s []int, v int) []int {
	return append(s, v)
}

// AppendViaPointer is the LESS COMMON fix: take *[]int and modify in
// place. Sometimes useful when you have multiple return values and
// adding "the new slice" to the list would be awkward -- but usually
// AppendReturning is cleaner.
func AppendViaPointer(s *[]int, v int) {
	*s = append(*s, v)
}

// Item is used to demonstrate the difference between []Item and []*Item
// in terms of mutation behavior inside a range loop.
type Item struct {
	Name  string
	Count int
}

// IncrementAllValues tries to increment Count on every Item in a slice
// of VALUES. This is a CLASSIC BUG: the range loop gives you a COPY of
// each item, so `item.Count++` modifies the copy, not the slice.
//
// We keep this function here so you can actually see the broken behavior
// and build an intuition for when you need to change approach.
func IncrementAllValues(items []Item) {
	for _, item := range items {
		item.Count++ // modifies the local copy -- DOES NOT AFFECT items
	}
}

// IncrementAllValuesFixed is the correct version using the INDEX. We
// access the actual slice element via `items[i]` and modify it directly.
func IncrementAllValuesFixed(items []Item) {
	for i := range items {
		items[i].Count++
	}
}

// IncrementAllPointers uses a SLICE OF POINTERS. The range loop gives
// us a copy of each POINTER, but both copies point at the same Item,
// so modifying via the pointer works fine.
//
// This is often nicer than index-based mutation when the loop body is
// complex. And it sidesteps the "value copy" gotcha entirely.
func IncrementAllPointers(items []*Item) {
	for _, item := range items {
		item.Count++ // `item` is *Item, auto-dereference for `.Count`
	}
}

// =============================================================================
// SECTION 6: POINTERS WITH MAPS
// =============================================================================

// AddToMap mutates a map via a plain (non-pointer) parameter. This works
// because maps, like slices, are reference-like: passing a map lets the
// function modify its entries directly.
//
// You almost never want `*map[K]V` -- just pass the map.
func AddToMap(m map[string]int, key string, value int) {
	m[key] = value
}

// IncrementMapValue demonstrates the famous "can't take address of map
// value" restriction. The following line WILL NOT COMPILE:
//
//	m["a"].Count++    // ERROR: cannot assign to struct field
//
// The reason: the map's internal storage can move entries around during
// rehashing, so any "address" would be unstable. Go simply disallows it.
//
// Workaround 1: read-modify-write. Copy out, mutate, put back.
func IncrementMapValue(m map[string]Item, key string) {
	it := m[key] // copy the value out (zero value if key missing)
	it.Count++
	m[key] = it // write the modified copy back
}

// IncrementMapPointer demonstrates Workaround 2: store POINTERS in the
// map. `m[key]` returns *Item; the auto-dereference lets us modify the
// pointed-to struct directly. This is usually the cleanest solution when
// you want mutable struct values in a map.
func IncrementMapPointer(m map[string]*Item, key string) {
	if p, ok := m[key]; ok {
		p.Count++ // modifies the Item that p points to
	}
}

// =============================================================================
// SECTION 7: NIL-POINTER PANIC DEMONSTRATION
// =============================================================================

// DerefNilRecovered deliberately dereferences a nil pointer to show what
// the panic looks like. We use `recover()` inside a deferred function
// to catch the panic so the program doesn't crash. (We'll cover defer/
// recover formally when you hit error handling in the roadmap; here it's
// just a wrapper so you can safely observe the panic message.)
//
// In real code, you should PREVENT nil derefs with guards, not catch
// them with recover. This function is purely demonstrative.
func DerefNilRecovered() (msg string) {
	defer func() {
		if r := recover(); r != nil {
			msg = fmt.Sprintf("recovered panic: %v", r)
		}
	}()

	var p *int
	_ = *p // PANIC: runtime error: invalid memory address or nil pointer dereference
	return "this line never runs"
}

// =============================================================================
// SECTION 8: MEMORY -- STACK VS HEAP (OBSERVATION ONLY)
// =============================================================================

// MakeOnHeap returns a pointer to a struct created inside the function.
// Because the pointer escapes (it's returned to the caller), the Go
// compiler's escape analysis will place this struct on the HEAP.
//
// You can verify this yourself by running:
//
//	go build -gcflags="-m" ./...
//
// The output will include something like:
//
//	"moved to heap: p"
//
// The caller never needs to know whether the allocation was on the stack
// or the heap -- the pointer works the same either way. But the choice
// affects allocation speed and GC pressure.
func MakeOnHeap() *Point {
	p := Point{X: 7, Y: 11}
	return &p
}

// NoEscape returns a Point BY VALUE. The local `p` is copied out at the
// return, and its original lives only within this function's stack frame.
// Escape analysis will typically keep this entirely on the stack -- no
// heap allocation, no GC involvement.
func NoEscape() Point {
	p := Point{X: 7, Y: 11}
	return p
}

// MemStats returns a small snapshot of the current memory state. Useful
// for getting a rough feel for allocation behavior.
//
// Fields (in bytes):
//
//	HeapAlloc    = current bytes allocated on the heap and still in use
//	TotalAlloc   = cumulative bytes allocated over the program's lifetime
//	NumGC        = number of completed GC cycles since start
//
// Call this before and after doing work to see how much got allocated.
func MemStats() (heapAlloc, totalAlloc uint64, numGC uint32) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc, m.TotalAlloc, m.NumGC
}

// ForceGC triggers an immediate garbage collection cycle. You almost
// never need this in production -- the runtime handles GC scheduling
// well on its own. It's here for observation in demos.
func ForceGC() {
	runtime.GC()
}

// =============================================================================
// SECTION 9: GC-FRIENDLY vs GC-HOSTILE ALLOCATION PATTERNS
// =============================================================================

// sink is a package-level variable we write to from demo functions to
// defeat the compiler's dead-store elimination. Without something like
// this, `make([]byte, 1024)` followed by no real use of `buf` can be
// optimized away entirely, and the demo would show no allocation at all.
var sink []byte

// AllocateInLoop deliberately does a wasteful thing: it allocates a new
// 1KB buffer every iteration. This puts pressure on both the allocator
// and the GC. Great as a "don't do this" baseline.
//
// To prevent the compiler from eliminating the allocation, we write to
// the buffer and store a reference to it in a package-level `sink`. Both
// are "uses" that the compiler can't prove away.
//
// Returns the total number of bytes touched, for display in the demo.
func AllocateInLoop(iterations int) int {
	total := 0
	for i := 0; i < iterations; i++ {
		buf := make([]byte, 1024) // fresh heap allocation every iteration
		buf[0] = byte(i)          // write so the alloc can't be elided
		sink = buf                // visible outside -- forces escape to heap
		total += len(buf)
	}
	return total
}

// ReuseBuffer does the same total work but with a SINGLE allocation
// outside the loop. Each iteration reuses the same backing array,
// which is dramatically friendlier to the GC.
//
// Pattern rule of thumb: if you need a temporary buffer in a hot loop,
// allocate it once and reuse.
func ReuseBuffer(iterations int) int {
	buf := make([]byte, 1024) // one allocation for the whole function
	total := 0
	for i := 0; i < iterations; i++ {
		buf[0] = byte(i) // same kind of work, no new allocations
		sink = buf
		total += len(buf)
	}
	return total
}

// PreallocateSlice shows the simplest GC-friendly habit you'll ever
// learn: when you know the final size of a slice, ask for it up front
// with `make([]T, 0, n)`. This avoids the regrow-and-copy allocations
// that happen as `append` outgrows its capacity.
func PreallocateSlice(n int) []int {
	out := make([]int, 0, n) // capacity n, length 0
	for i := 0; i < n; i++ {
		out = append(out, i*i) // no reallocations -- we pre-sized the backing array
	}
	return out
}

// =============================================================================
// SECTION 10: A PRACTICAL EXAMPLE -- A LINKED LIST WITH POINTERS
// =============================================================================

// Node is a linked list node. A linked list is the canonical "you must
// use pointers" data structure: each node holds a reference to the next,
// and using values instead would mean each node contains an infinitely
// recursive copy of its successor -- which the language won't allow.
//
// (Try changing `*Node` to `Node` and see the compile error.)
type Node struct {
	Value int
	Next  *Node // pointer to the next node; nil means end of list
}

// LinkedList holds a pointer to the head node. Using a pointer here
// means methods on *LinkedList can replace the head entirely (e.g.,
// when prepending) without any weird gymnastics.
type LinkedList struct {
	Head *Node
}

// NewLinkedList is the conventional constructor. Returns a *LinkedList
// so methods can mutate it.
func NewLinkedList() *LinkedList {
	return &LinkedList{Head: nil}
}

// Prepend adds a new node at the FRONT of the list. O(1) operation.
// Notice how we use `l.Head` both as "read the old head" and as "assign
// the new head" -- pointer manipulation in one line.
func (l *LinkedList) Prepend(v int) {
	l.Head = &Node{Value: v, Next: l.Head}
}

// Append adds a new node at the END of the list. O(n) operation because
// we have to walk to the tail.
//
// The logic: if the list is empty, the new node BECOMES the head.
// Otherwise, walk until we find the node whose Next is nil, and attach
// the new node there.
func (l *LinkedList) Append(v int) {
	newNode := &Node{Value: v} // Next defaults to nil -- the zero value for pointers
	if l.Head == nil {
		l.Head = newNode
		return
	}
	// Walk to the tail. `current` is our moving cursor.
	current := l.Head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode
}

// ToSlice converts the list to a []int for easy display. It walks from
// head to tail, collecting values as it goes.
func (l *LinkedList) ToSlice() []int {
	var out []int
	for current := l.Head; current != nil; current = current.Next {
		out = append(out, current.Value)
	}
	return out
}

// Length returns the number of nodes, iteratively. We could also do this
// recursively, but iterative is Go's idiomatic style (no tail call
// optimization, remember).
func (l *LinkedList) Length() int {
	count := 0
	for current := l.Head; current != nil; current = current.Next {
		count++
	}
	return count
}

// RemoveFirst removes the given value's FIRST occurrence from the list.
// Returns true if something was removed, false otherwise.
//
// This function is the clearest demonstration of why linked-list code
// loves pointers: to remove a node, we just change the predecessor's
// Next pointer to skip over it. The GC cleans up the unreachable node.
//
// No manual free(), no reference counting, no memory leak. Pointers +
// GC = linked lists that look a lot like they do in textbooks but
// without the usual C-style manual memory management.
func (l *LinkedList) RemoveFirst(v int) bool {
	// Special case: removing the head.
	if l.Head != nil && l.Head.Value == v {
		l.Head = l.Head.Next
		return true
	}
	// Walk with two pointers: `prev` trails `current` by one step.
	for prev, current := l.Head, l.Head; current != nil; current = current.Next {
		if current.Value == v {
			prev.Next = current.Next // skip over `current`; GC will reclaim it
			return true
		}
		prev = current
	}
	return false
}
