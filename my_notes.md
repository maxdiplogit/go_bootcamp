# Go Programming — Complete Study Notes

A comprehensive personal reference covering everything from types and structs through concurrency patterns, with all the intuition-building Q&A along the way.

---

## Table of Contents

1. [Types and Structs](#1-types-and-structs)
2. [Methods](#2-methods)
3. [Interfaces — Part 1: Method Sets](#3-interfaces--part-1-method-sets)
4. [Struct Tags and JSON](#4-struct-tags-and-json)
5. [Struct Embedding](#5-struct-embedding)
6. [Strings, Bytes, and Runes](#6-strings-bytes-and-runes)
7. [Loops](#7-loops)
8. [Functions](#8-functions)
9. [Pointers](#9-pointers)
10. [Memory Management and Garbage Collection](#10-memory-management-and-garbage-collection)
11. [Interfaces — Part 2: How They Work Internally](#11-interfaces--part-2-how-they-work-internally)
12. [Generics](#12-generics)
13. [Interfaces vs Type Constraints (The Two Jobs of `interface`)](#13-interfaces-vs-type-constraints-the-two-jobs-of-interface)
14. [Error Handling](#14-error-handling)
15. [Concurrency Foundations](#15-concurrency-foundations)
16. [The `context` Package](#16-the-context-package)
17. [Concurrency Patterns](#17-concurrency-patterns)
18. [Channels Internals — Deep Dive](#18-channels-internals--deep-dive)
19. [Quick Reference Cards](#19-quick-reference-cards)

---

## 1. Types and Structs

### Why Types Exist

Memory is just bytes. The number `65`, the character `'A'`, and the byte `0x41` could all be the same bits — types tell the compiler what those bits *mean* and what operations are valid.

Go is statically and strongly typed: types are checked at compile time, with no implicit conversions.

```go
var a int = 5
var b float64 = 3.2
// var c = a + b  // COMPILE ERROR
var c = float64(a) + b  // explicit conversion required
```

### Named (Defined) Types

```go
type Celsius float64
type Fahrenheit float64

var bodyTemp Celsius = 37.0
var roomTemp Fahrenheit = 72.0
// bodyTemp = roomTemp           // ERROR — different types
bodyTemp = Celsius(roomTemp)     // explicit conversion
```

The compiler enforces semantic distinctions. (The Mars Climate Orbiter crashed because of an implicit unit conversion in C.)

### Structs — Grouping Related Data

```go
type Person struct {
    Name  string
    Age   int
    Email string
    City  string
}
```

#### Creating Structs

```go
// 1. Zero value — all fields get zero values
var p1 Person

// 2. Struct literal with field names (BEST)
p2 := Person{
    Name:  "Asha",
    Age:   28,
    Email: "asha@example.com",
}

// 3. Positional literal (FRAGILE — breaks if fields reorder)
p3 := Person{"Asha", 28, "asha@example.com", "Bengaluru"}

// 4. Pointer via &
p4 := &Person{Name: "Asha", Age: 28}

// 5. Pointer via new() — returns *Person with zero values
p5 := new(Person)
p5.Name = "Asha"  // auto-dereference
```

### Struct Alignment and Padding

CPUs read memory in chunks (4 or 8 bytes at a time). For correctness on some architectures, data must be **aligned** — an 8-byte value at an address divisible by 8, etc. The compiler inserts invisible **padding** to satisfy this.

Type alignments:

| Type | Size | Alignment |
|---|---|---|
| `bool`, `int8`, `uint8` | 1 | 1 |
| `int16`, `uint16` | 2 | 2 |
| `int32`, `uint32`, `float32` | 4 | 4 |
| `int64`, `uint64`, `float64` | 8 | 8 |
| `string` | 16 | 8 |
| pointer, `int`, `uint` | 8 | 8 |

Example:

```go
type Bad struct {
    a bool   // 1 byte
    b int64  // 8 bytes
    c bool   // 1 byte
}
// Total: 24 bytes!
```

Layout: `a` at offset 0, then 7 bytes of padding so `b` aligns at offset 8, then `c` at offset 16, then 7 bytes of trailing padding.

```go
type Good struct {
    b int64  // 8 bytes
    a bool   // 1 byte
    c bool   // 1 byte
}
// Total: 16 bytes
```

**Rule:** order fields from largest to smallest. Verify with `unsafe.Sizeof()`.

---

## 2. Methods

### Definition

A method is a function with a special receiver parameter that binds it to a type.

```go
type Rectangle struct {
    Width, Height float64
}

// Value receiver — operates on a COPY
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Pointer receiver — operates on the ACTUAL rectangle
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}
```

### Why Methods Instead of Functions

1. **Discoverability** — IDE autocomplete shows everything you can do with a `Rectangle`.
2. **Namespacing** — multiple types can have an `Area()` method without collision.
3. **Interfaces** — interfaces are defined in terms of methods.

### The Method Set Rule

- **Value of type `T`** — has access to value-receiver methods only.
- **Pointer `*T`** — has access to both value-receiver and pointer-receiver methods.

Why? If you have only a value, Go can take its address to call a pointer method *most of the time*. But when assigning to an interface, the value doesn't have a stable address, so pointer methods aren't in its method set.

**Convention:** if any method on a type uses a pointer receiver, make all of them pointer receivers for consistency.

### Call by Value vs Call by Pointer

Go **always** passes by value. The interesting question is: what gets copied?

```go
type Counter struct {
    Count int
}

// Modifies a COPY
func incrementValue(c Counter) {
    c.Count++  // throwaway change
}

// Modifies through pointer
func incrementPointer(c *Counter) {
    c.Count++  // auto-dereferenced from (*c).Count
}

c := Counter{Count: 0}
incrementValue(c)
fmt.Println(c.Count)   // 0 — original untouched

incrementPointer(&c)
fmt.Println(c.Count)   // 1 — modified via pointer
```

#### When to Use Pointers vs Values

Pointers when:
1. **Mutation** — function needs to modify the original.
2. **Avoiding expensive copies** — large structs (>~64 bytes).
3. **Signaling "no value"** — pointers can be nil.

Values when:
- Small types (a few fields).
- No mutation needed.
- No nil concerns.

---

## 3. Interfaces — Part 1: Method Sets

### Why Interfaces Exist

Imagine a function that writes data somewhere — a file, a network connection, a buffer, stdout. Without interfaces, you'd write five copies. With interfaces:

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}

func writeReport(w Writer, data string) {
    w.Write([]byte(data))
}
```

Now any type with the right `Write` method works. **Files, connections, buffers, stdout — they all already have it.**

### Defining an Interface

```go
type Greeter interface {
    Greet() string
}

type EnglishSpeaker struct{ Name string }
func (e EnglishSpeaker) Greet() string { return "Hello, " + e.Name }

func sayHi(g Greeter) {
    fmt.Println(g.Greet())
}

sayHi(EnglishSpeaker{Name: "Asha"})
```

**`EnglishSpeaker` never says "I implement Greeter."** No `implements` keyword. The compiler just checks: does the type have the right methods? Yes → it satisfies the interface. This is **structural typing**.

The huge benefit: types and interfaces are **decoupled**. A type written today can satisfy an interface defined tomorrow without modification.

### The Empty Interface — `any`

```go
type Anything interface{}  // same as `any`

func printIt(v any) {
    fmt.Println(v)
}
```

`any` is satisfied by every type. Use it sparingly — it throws away type safety. Generics are usually better now (Go 1.18+).

### Type Assertions

```go
var v any = "hello"

s := v.(string)         // panics on mismatch
s, ok := v.(string)     // safe two-value form: s="hello", ok=true
n, ok := v.(int)        // n=0, ok=false (no panic)
```

Always use the two-value form unless you're certain.

### Type Switches

```go
func describe(v any) string {
    switch x := v.(type) {
    case int:
        return fmt.Sprintf("int: %d", x*2)
    case string:
        return fmt.Sprintf("string of length %d", len(x))
    case Rectangle:
        return fmt.Sprintf("area: %f", x.Area())
    default:
        return "unknown"
    }
}
```

`v.(type)` only works inside a `switch`. Inside each case, `x` already has the right type.

### Pointer Receivers and Interface Satisfaction

```go
type Counter interface { Increment() }
type MyCounter struct{ N int }

func (c *MyCounter) Increment() { c.N++ }   // pointer receiver

var c Counter = MyCounter{}     // ERROR — value doesn't satisfy
var c Counter = &MyCounter{}    // OK
```

The rule: methods with pointer receivers are in the method set of `*T` only, not `T`. So a value can't satisfy an interface whose methods use pointer receivers.

### Embedding Interfaces

```go
type Reader interface { Read(p []byte) (int, error) }
type Writer interface { Write(p []byte) (int, error) }

type ReadWriter interface {
    Reader
    Writer
}
```

`ReadWriter` satisfies anything that has both Read and Write. The standard library uses this constantly: `io.ReadWriter`, `io.ReadCloser`, etc.

---

## 4. Struct Tags and JSON

### Why Struct Tags Exist

Go uses `PascalCase`. JSON APIs use `snake_case` or `camelCase`. SQL columns use other conventions. Struct tags are metadata that bridges these worlds.

```go
type User struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
}
```

Multiple tags can coexist:

```go
type User struct {
    ID       int    `json:"id" db:"user_id" validate:"required"`
    Email    string `json:"email" db:"email" validate:"required,email"`
    Password string `json:"-" db:"password_hash"`
}
```

Use **backticks**, not double quotes. No space around `:` in the tag.

### JSON Operations

```go
// Marshal (Go → JSON)
data, err := json.Marshal(user)
data, err := json.MarshalIndent(user, "", "  ")  // pretty

// Unmarshal (JSON → Go)
var u User
err := json.Unmarshal(jsonBytes, &u)  // pass pointer
```

### Tag Options

```go
type Product struct {
    Name        string   `json:"name"`
    Description string   `json:"description,omitempty"`  // skip if empty
    Tags        []string `json:"tags,omitempty"`
    Internal    string   `json:"-"`                       // never serialize
    BigInt      int64    `json:"big_int,string"`          // encode as string
}
```

### The Zero Value Trap

```go
type UpdateRequest struct {
    Active bool   `json:"active,omitempty"`  // BUG: false will be omitted
}
```

Fix with pointers:

```go
type UpdateRequest struct {
    Active *bool `json:"active,omitempty"`  // nil = not provided, &false = explicitly false
}
```

This is the standard pattern for PATCH-style partial-update APIs.

### Best Practices

1. **Separate API types from domain types** — never expose a `PasswordHash` field by accident.
2. **Validate after unmarshal** — JSON decoding doesn't validate.
3. **Use `DisallowUnknownFields()` for strict APIs** — catches client typos.
4. **Stream large payloads** with `json.NewEncoder(w).Encode(v)`.
5. **Use `json.RawMessage` for pass-through data** you don't need to parse.

---

## 5. Struct Embedding

### What It Solves

Instead of repeating `ID`, `CreatedAt`, `UpdatedAt` in every struct, factor them out and embed:

```go
type Metadata struct {
    ID        int64
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Article struct {
    Metadata        // embedded — no field name
    Body string
}

type Video struct {
    Metadata
    URL string
}
```

### How Embedding Works

The compiler doesn't merge fields. The embedded struct still exists as a unit. **Promotion** is just syntactic sugar — `article.Title` becomes `article.Metadata.Title` automatically.

```go
a := Article{
    Metadata: Metadata{ID: 1, CreatedAt: time.Now()},
    Body:     "...",
}

a.ID                // promoted access — same as a.Metadata.ID
a.Metadata.ID       // explicit path also works
```

### Initialization Quirk

You **cannot** initialize embedded fields inline:

```go
// FAILS:
a := Article{
    ID:    1,        // ERROR
    Body:  "...",
}

// CORRECT:
a := Article{
    Metadata: Metadata{ID: 1},
    Body:     "...",
}
```

### Field Shadowing

```go
type Article struct {
    Metadata
    Title string  // shadows Metadata.Title if it exists
}

a.Title           // outer wins
a.Metadata.Title  // inner still reachable
```

### Embedding + JSON

By default, embedded fields are **flattened** in JSON output:

```go
type Article struct {
    Metadata
    Body string `json:"body"`
}
// Output: {"id": 1, "created_at": "...", "body": "..."}
```

To force nesting, use a named field:

```go
type Article struct {
    Meta Metadata `json:"meta"`
    Body string   `json:"body"`
}
// Output: {"meta": {"id": 1, ...}, "body": "..."}
```

---

## 6. Strings, Bytes, and Runes

### The Core Mental Model

```
string  →  immutable sequence of BYTES (conventionally UTF-8)
byte    →  alias for uint8 (one raw byte, 0-255)
rune    →  alias for int32 (one Unicode code point)
```

### Why Bytes Aren't Enough

```go
s := "hello"
fmt.Println(len(s))  // 5 bytes
fmt.Println(s[0])    // 104 (byte for 'h')
```

But:

```go
s := "héllo"
fmt.Println(len(s))  // 6, not 5! 'é' takes 2 bytes in UTF-8

s := "héllo 世界 🙂"
fmt.Println(len(s))  // 17 bytes (1+2+1+1+1+1+3+3+1+4)
```

**`len()` counts BYTES, not characters.** `s[i]` returns a byte. Slicing works on byte boundaries.

### Indexing vs Iterating

```go
// Byte iteration — broken for non-ASCII
for i := 0; i < len(s); i++ {
    fmt.Println(s[i])  // raw byte values
}

// Rune iteration — decodes UTF-8 correctly
for i, r := range s {
    fmt.Printf("%d: %c\n", i, r)
}
```

`range` on a string gives `(byteIndex, rune)`. The byte index is **not sequential** — it's the offset where each rune starts.

### Counting Characters Correctly

```go
import "unicode/utf8"

s := "héllo 世界"
fmt.Println(len(s))                    // 12 bytes
fmt.Println(utf8.RuneCountInString(s)) // 8 characters
```

### Conversions

```go
// string ↔ []byte
b := []byte(s)        // copies bytes
s := string(b)        // copies bytes
// Byte-level — fine for ASCII, splits multi-byte chars

// string ↔ []rune
r := []rune(s)        // decodes UTF-8
s := string(r)        // encodes back
// Character-level — safe for any Unicode
```

### Reverse a String — The Classic Trap

```go
// BROKEN for non-ASCII
func reverseBroken(s string) string {
    b := []byte(s)
    for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
        b[i], b[j] = b[j], b[i]
    }
    return string(b)
}
// reverseBroken("héllo") → "oll\xa9\xc3h" (corrupted!)

// CORRECT
func reverse(s string) string {
    r := []rune(s)
    for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r)
}
```

### Efficient Concatenation

```go
// SLOW — quadratic
s := ""
for i := 0; i < 10000; i++ {
    s += "x"
}

// FAST — strings.Builder
var b strings.Builder
b.Grow(10000)  // pre-allocate
for i := 0; i < 10000; i++ {
    b.WriteString("x")
}
s := b.String()
```

### Rules of Thumb

1. `len(s)` is bytes, not characters.
2. `s[i]` is a byte. Use `range` for runes.
3. Don't mutate strings — convert to `[]byte` or `[]rune`, modify, convert back.
4. ASCII-only logic can use `[]byte` safely. User text should go through runes.
5. Use `strings.Builder` for loops that build strings.
6. `strings.Index` returns byte offsets, not character offsets.

---

## 7. Loops

### One Keyword, Four Shapes

Go has only `for`. Four forms:

```go
// 1. Three-part (classic counted loop)
for i := 0; i < 10; i++ { ... }

// 2. Condition-only (Go's `while`)
for x < 100 { ... }

// 3. Infinite (Go's `while true`)
for { ... }

// 4. Range-based
for i, v := range collection { ... }
```

### Range Specifics

**Slices/arrays:** `(index, value-copy)`

```go
nums := []int{10, 20, 30}
for i, v := range nums {
    v *= 2  // modifies copy, NOT slice!
}
// nums unchanged

// To mutate, use index:
for i := range nums {
    nums[i] *= 2
}
```

**Strings:** `(byteIndex, rune)` — decodes UTF-8.

**Maps:** `(key, value)` — iteration order is **randomized**.

```go
ages := map[string]int{"Asha": 28, "Bo": 30}
for k, v := range ages {
    fmt.Println(k, v)
}
```

For sorted iteration:

```go
keys := make([]string, 0, len(ages))
for k := range ages {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Println(k, ages[k])
}
```

**Integer (Go 1.22+):**

```go
for i := range 5 {
    fmt.Println(i)  // 0, 1, 2, 3, 4
}
```

### break, continue, Labels

```go
outer:
for i := 0; i < 5; i++ {
    for j := 0; j < 5; j++ {
        if i*j > 6 {
            break outer  // exits BOTH loops
        }
        if j > i {
            continue outer  // skip to outer's next iteration
        }
    }
}
```

### Loop Variable Capture (Go 1.22 Fix)

Pre-1.22, loop variables were shared across iterations:

```go
funcs := []func(){}
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() { fmt.Println(i) })
}
for _, f := range funcs { f() }
// Pre-1.22: 3, 3, 3
// 1.22+:    0, 1, 2  (per-iteration variable now)
```

The old workaround was `i := i` at the top of the loop body. You'll still see this pattern in older code.

### Common Patterns

```go
// Filter
var keep []int
for _, v := range nums {
    if v > 0 { keep = append(keep, v) }
}

// Sum
total := 0
for _, v := range nums { total += v }

// Max
max := nums[0]
for _, v := range nums[1:] {
    if v > max { max = v }
}

// Transform
doubled := make([]int, len(nums))
for i, v := range nums { doubled[i] = v * 2 }

// Group by
byKey := map[byte][]string{}
for _, name := range names {
    byKey[name[0]] = append(byKey[name[0]], name)
}
```

---

## 8. Functions

### Basics

```go
func add(a, b int) int {  // type shorthand: a int, b int → a, b int
    return a + b
}
```

No default parameters. No keyword arguments. No overloading.

### Multiple Return Values

```go
func divide(a, b float64) (float64, float64) {
    return a / b, math.Mod(a, b)
}

q, r := divide(17, 5)
```

The canonical Go pattern: `(value, error)` or `(value, ok)`.

```go
n, err := parseInt("42")
if err != nil { ... }

v, ok := someMap[key]
```

### Named Return Values

```go
func divide(a, b float64) (quotient, remainder float64) {
    quotient = a / b
    remainder = math.Mod(a, b)
    return  // naked return — sends current values
}
```

Useful for documentation, automatic zero-value initialization, and `defer` interaction.

### Variadic Functions

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums { total += n }
    return total
}

sum(1, 2, 3)         // 6
sum()                // 0

// Spread an existing slice
nums := []int{1, 2, 3}
sum(nums...)         // 6
```

The variadic must be **last**. You can have any number of fixed parameters before it.

### Functions as Values (First-Class)

```go
type IntOp func(int, int) int

var add IntOp = func(a, b int) int { return a + b }

func apply(a, b int, op IntOp) int {
    return op(a, b)
}

// Inline anonymous
result := apply(3, 4, func(a, b int) int { return a * b })
```

### Closures

A closure is a function that captures variables from its surrounding scope:

```go
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

c1 := makeCounter()
c1()  // 1
c1()  // 2
c1()  // 3

c2 := makeCounter()  // independent — fresh count
c2()  // 1
```

Common patterns:

```go
// Configuration captured at creation time
func multiplier(factor int) func(int) int {
    return func(n int) int { return n * factor }
}
double := multiplier(2)
triple := multiplier(3)

// Encapsulating private state
func makeAccount(initial float64) (deposit func(float64), balance func() float64) {
    bal := initial
    deposit = func(amt float64) { bal += amt }
    balance = func() float64 { return bal }
    return
}
```

---

## 9. Pointers

### What a Pointer Is

A pointer is a variable whose value is the memory address of another variable.

```go
x := 42
p := &x       // p holds the address of x
fmt.Println(p)   // 0xc00001c030 (or similar)
fmt.Println(*p)  // 42 (dereference)
*p = 100         // write through pointer
fmt.Println(x)   // 100
```

### The Two Operators

- `&x` — address of x (gives a pointer)
- `*p` — value at p (dereference) OR `*T` as a type (pointer-to-T)

Same symbol, two meanings:
- Before a **type**: `*int` means "pointer to int" (type modifier)
- Before an **expression**: `*p` means "dereference p" (operator)

### Nil Pointers

Zero value of any pointer type is `nil`. Dereferencing nil panics:

```go
var p *int
fmt.Println(p == nil)  // true
fmt.Println(*p)        // PANIC

// Defensive check
if p != nil {
    fmt.Println(*p)
}
```

### Three Ways to Create Pointers

```go
// 1. Take address
x := 42
p := &x

// 2. new()
p := new(int)
*p = 42

// 3. Constructor function
func NewCounter() *Counter {
    return &Counter{Count: 0}
}
c := NewCounter()
```

### Returning Local Pointers Is Safe

```go
// In C this would be a bug. In Go it's fine.
func NewPoint(x, y int) *Point {
    p := Point{X: x, Y: y}
    return &p   // safe — escape analysis moves p to heap
}
```

The Go compiler does **escape analysis** — if a value's address escapes the function, it's allocated on the heap automatically. You never write `malloc`/`free`.

### Pointers with Slices

Slices are already reference-like. The slice header (`{ptr, len, cap}`) is copied on pass, but the backing array is shared.

```go
func doubleAll(s []int) {
    for i := range s {
        s[i] *= 2  // affects caller's data
    }
}
```

The exception — when you need `*[]int`:

```go
// PROBLEM: append may reallocate; caller doesn't see new slice
func addValue(s []int, v int) {
    s = append(s, v)  // local change only
}

// FIX 1: return the slice (idiomatic)
func addValue(s []int, v int) []int {
    return append(s, v)
}

// FIX 2: pass *[]int (less common)
func addValue(s *[]int, v int) {
    *s = append(*s, v)
}
```

### Slice of Values vs Pointers

```go
// Slice of values — embedded in backing array
people1 := []Person{{Name: "Asha"}, {Name: "Bo"}}

// Slice of pointers — each Person allocated separately
people2 := []*Person{{Name: "Asha"}, {Name: "Bo"}}

// Range gotcha:
for _, p := range people1 {
    p.Age++  // BUG — p is a copy, doesn't affect slice
}
for i := range people1 {
    people1[i].Age++  // correct — index access
}

for _, p := range people2 {
    p.Age++  // works — p is *Person
}
```

### Pointers with Maps

Maps are reference-like — pass directly, not as `*map[K]V`.

**The restriction:** can't take address of map value:

```go
m := map[string]Counter{"a": {Count: 1}}
m["a"].Count++  // ERROR

// Workaround 1: read-modify-write
c := m["a"]
c.Count++
m["a"] = c

// Workaround 2: store pointers
m := map[string]*Counter{"a": {Count: 1}}
m["a"].Count++  // works
```

### Pointer Type Decision Table

| Type | Default | Use pointer when |
|---|---|---|
| `int`, `bool`, `float64` | value | Almost never |
| Small struct | value | Need to mutate |
| Large struct | pointer | Default — avoid copying |
| `string` | value | Never (already a header) |
| `[]T` (slice) | value | Need to replace whole slice |
| `map[K]V` | value | Never |
| `chan T` | value | Never |

---

## 10. Memory Management and Garbage Collection

### Stack vs Heap

**Stack** — fast, automatic. Each function call gets a stack frame. Allocation is essentially free; deallocation is automatic when the function returns.

**Heap** — slower, GC-managed. Lives until garbage collector determines it's unreachable.

**Go decides which one** via escape analysis. You don't manage this yourself.

```go
func stackAlloc() int {
    x := 42      // lives on stack
    return x     // value copied out; x discarded
}

func heapAlloc() *int {
    x := 42      // ESCAPES — address leaves the function
    return &x    // compiler moves x to heap
}
```

### Seeing Escape Analysis

```bash
go build -gcflags="-m" yourfile.go
```

Output like `moved to heap: x` tells you what the compiler decided.

### Garbage Collection

Go uses a concurrent, tri-color, mark-sweep collector:

1. **Mark** — starting from roots (globals, goroutine stacks), walk the reference graph. Mark everything reachable.
2. **Sweep** — reclaim memory of unmarked (unreachable) objects.

Most of this happens **concurrently** with your program. Pauses are typically sub-millisecond.

### What the GC Does NOT Do

- **Doesn't free immediately** when something becomes unreachable. Runs periodically.
- **Doesn't manage stack memory.** Stack auto-frees on function return.
- **Doesn't reference-count.** Determines reachability from roots.

### Memory Leaks in Go

Memory leaks happen when you keep references to things you don't need:

```go
var cache = map[string]*BigObject{}

func process(key string) {
    cache[key] = loadBigObject()  // grows forever unless you delete
}
```

Also: goroutine leaks (a goroutine blocked forever and never exits).

### GC-Friendly Patterns

```go
// Pre-allocate slices when size is known
out := make([]int, 0, len(in))

// Reuse buffers in hot loops
buf := make([]byte, 1024)
for i := 0; i < 1000000; i++ {
    process(buf)  // reuse, don't allocate
}

// Prefer values for small structs
// Small Point{X, Y int} fits in 16 bytes — pointer adds nothing
```

### Demonstrated Allocation Pressure

In the pointers reference, we measured:
- `AllocateInLoop(10000)` — allocates 1KB per iteration → ~10MB total alloc, 3 GC cycles
- `ReuseBuffer(10000)` — single 1KB allocation reused → ~1KB total alloc, 0 GC cycles

Same total work; orders of magnitude difference in GC pressure.

---

## 11. Interfaces — Part 2: How They Work Internally

### Q: How are interfaces stored in memory?

An interface value is **two pointers wide** (16 bytes on 64-bit):

```
interface value = (type_descriptor_pointer, data_pointer)
```

- **type descriptor**: a pointer to information about what concrete type is inside, including a method-lookup table (called `itab` internally).
- **data pointer**: a pointer to the actual value.

### Q: Are methods stored separately for every instance?

**No.** Method code is stored once, per type. There's one `Write` function for `*FileLogger` regardless of how many `FileLogger`s you create. What changes per call is the receiver passed in.

### Q: Why does the interface need to store type info?

To know which method to call. If you have an interface value `w` that could be holding a `*FileLogger` or a `*MemoryLogger`, calling `w.Write(...)` needs to find the right `Write` function. The type descriptor has that information.

### Q: What happens when we call a method on an interface?

Step by step:
1. Look at the first slot of the interface value — get the type descriptor pointer.
2. In the type descriptor, find the method's address (e.g., `Write`).
3. Take the second slot — the data pointer (the receiver).
4. Call the function at that address, passing the data pointer as the receiver.

This is the indirection that interfaces add — a small constant cost.

### Q: Is the data passed by value or pointer in the interface?

The data slot of the interface is **always a pointer**. Two cases:

**You pass `&value`:** interface stores your pointer directly. Both the original variable and the interface refer to the same struct.

**You pass `value` (not pointer):** Go copies `value` into a fresh location and stores a pointer to that copy. Mutations made through the interface affect the copy, not the original.

This is why pointer-receiver methods are not in the value type's method set — letting `var c Counter = MyCounter{}` would make `c.Increment()` modify a throwaway copy, silently breaking things.

### Q: Will `sayHi(e)` work when `e` is `EnglishSpeaker` value?

Yes, **provided `Greet()` has a value receiver**.

```go
e := EnglishSpeaker{Name: ""}
sayHi(e)  // works fine — prints "Hello, "
```

The compile error case from earlier was different: `Counter` had a pointer-receiver method, so a value couldn't satisfy it.

### The Mental Model in Five Lines

1. Method code is shared per type, not per instance.
2. An interface value is two pointers: type info + data.
3. Type info exists once per (type, interface) pair.
4. The data slot points at the actual thing (your pointer, or a copy).
5. Method calls do an indirect lookup: type info → method address → call with data pointer.

---

## 12. Generics

### Why Generics Exist

Without generics:

```go
func MaxInt(a, b int) int { ... }
func MaxFloat(a, b float64) float64 { ... }
func MaxString(a, b string) string { ... }
// Three near-identical functions
```

Interfaces don't quite work because (a) basic types don't have a `GreaterThan` method, and (b) you'd lose type information through the interface.

Generics let you write the function once, with placeholder types:

```go
func Max[T cmp.Ordered](a, b T) T {
    if a > b { return a }
    return b
}

Max(3, 7)             // T = int
Max(3.14, 2.71)       // T = float64
Max("hi", "world")    // T = string
```

The compiler stamps out specialized versions per concrete type. Full type safety, no boxing.

### Type Parameters

```go
func Max[T cmp.Ordered](a, b T) T  // single param
func Pair[K, V any](k K, v V)      // two params
```

Convention: single uppercase letters (`T`, `K`, `V`).

### Constraints

```go
[T any]                          // any type
[T comparable]                   // types that support ==, !=
[T int | float64 | string]       // union of types
[T constraints.Ordered]          // predefined constraint
[T Stringer]                     // interface name
```

The constraint determines:
1. Which types the caller can use.
2. Which operations work on `T` inside the function.

### Custom Constraints

```go
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
        ~float32 | ~float64
}

func Sum[T Number](nums []T) T {
    var total T
    for _, v := range nums {
        total += v
    }
    return total
}
```

The `~` prefix means "this type or any type defined from it." `~int` accepts `int` and any `type MyInt int`.

### Type Inference

```go
Max(1, 2)             // T inferred as int
Max[int](1, 2)        // explicit
```

Most calls don't need explicit type parameters.

### Generic Types

```go
type Stack[T any] struct {
    items []T
}

func NewStack[T any]() *Stack[T] {
    return &Stack[T]{}
}

func (s *Stack[T]) Push(v T) {
    s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    last := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return last, true
}

intStack := NewStack[int]()
strStack := NewStack[string]()
```

Note: `var zero T` produces the zero value of T. `Stack[int]` and `Stack[string]` are **different types**.

### Generic Higher-Order Functions

```go
func MapSlice[T, U any](in []T, f func(T) U) []U {
    out := make([]U, len(in))
    for i, v := range in {
        out[i] = f(v)
    }
    return out
}

func FilterSlice[T any](in []T, keep func(T) bool) []T {
    out := make([]T, 0, len(in))
    for _, v := range in {
        if keep(v) { out = append(out, v) }
    }
    return out
}
```

### When to Use Generics

- **Same code, different types:** generics
- **Different code, common behavior:** interfaces
- **One specific type:** plain function

Don't reach for generics first. Most application code doesn't need them. They shine in libraries and shared utilities.

---

## 13. Interfaces vs Type Constraints (The Two Jobs of `interface`)

The keyword `interface` does **two different jobs**, and the syntax is similar but the meaning is very different.

### Job 1: Method Set (Runtime)

```go
type Greeter interface {
    Greet() string
}
```

- Describes "any type with these methods."
- Used as variable, parameter, or return type.
- Creates real runtime objects (the 16-byte interface value).
- Method dispatch is dynamic.

### Job 2: Type Set / Constraint (Compile-Time)

```go
type Number interface {
    int | float64 | string
}
```

- Describes "any type from this list."
- Used **only** inside `[T constraint]` brackets in generics.
- Pure compile-time check — no runtime presence.
- Cannot be used to declare variables.

```go
var n Number = 42  // COMPILE ERROR: interface contains type constraints
```

### Mixed Constraints

You can combine both:

```go
type StringableNumber interface {
    ~int | ~float64       // type union
    String() string       // method requirement
}
```

A type satisfies this if it's int- or float-based AND has a `String()` method. **But because it has a type union, it's constraint-only — you can't declare variables of this type.**

### The Decision Tree

| Body of interface | Can be used as | 
|---|---|
| Only methods | Variable type AND constraint |
| Has type union (anywhere) | Constraint only |

### Realistic Mixed-Constraint Example

```go
type Number interface {
    ~int | ~int64 | ~float64
}

type Printable interface {
    Format() string
}

type FormattableNumber interface {
    Number       // type union
    Printable    // method requirement
}

func DescribeMax[T FormattableNumber](a, b T) string {
    bigger := a
    if b > a { bigger = b }       // > works because all numbers
    return "max = " + bigger.Format()  // Format() works because constraint
}
```

The constraints in the interface mean two things at compile time:
1. Restrict which types are allowed for the type parameter.
2. Tell the compiler which operations are valid on values of that type.

---

## 14. Error Handling

### Errors Are Values

No exceptions. Functions return an extra `error` value:

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil
}

result, err := divide(10, 0)
if err != nil {
    fmt.Println("error:", err)
    return
}
```

The `if err != nil { return err }` pattern is everywhere.

### The `error` Interface

```go
type error interface {
    Error() string
}
```

Any type with `Error() string` is an error.

### Creating Errors

```go
err := errors.New("simple message")
err := fmt.Errorf("formatted: %s with %d", name, count)
```

### Custom Error Types

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func validate(name string) error {
    if name == "" {
        return &ValidationError{Field: "name", Message: "must not be empty"}
    }
    return nil
}
```

### Sentinel Errors

```go
var ErrNotFound = errors.New("user not found")
var ErrInvalidInput = errors.New("invalid input")

func Get(id int) (*User, error) {
    if !exists(id) {
        return nil, ErrNotFound
    }
    return load(id), nil
}

// Caller checks:
user, err := Get(42)
if err == ErrNotFound { ... }    // works for unwrapped
if errors.Is(err, ErrNotFound) { ... }  // works even if wrapped — preferred
```

Convention: prefix with `Err`, exported, declared at package level.

### Wrapping Errors with `%w`

```go
func loadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("loadConfig %s: %w", path, err)
    }
    // ...
}
```

`%w` preserves the original error inside the new one. Each layer adds context.

`%v` includes the error text but doesn't preserve linkage. Use `%w` when wrapping; `%v` only when deliberately replacing.

### `errors.Is` and `errors.As`

```go
// Walk the wrap chain looking for a sentinel
if errors.Is(err, os.ErrNotExist) { ... }

// Walk the wrap chain looking for a specific TYPE
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println(pathErr.Path)
}
```

**Always prefer these over `==`** — they handle wrapping correctly.

### `defer`

Schedules a function call to run when the surrounding function returns:

```go
func process(path string) error {
    f, err := os.Open(path)
    if err != nil { return err }
    defer f.Close()  // runs when function returns
    // ...
    return nil
}
```

Rules:
1. **LIFO order** — last deferred runs first.
2. **Arguments evaluated when `defer` is called**, not when it runs.
3. **Closures capture by reference** — they see the latest values.

```go
defer fmt.Println("1")   // runs 3rd
defer fmt.Println("2")   // runs 2nd
defer fmt.Println("3")   // runs 1st

x := 1
defer fmt.Println(x)         // captures 1 NOW, prints 1 later
defer func() { fmt.Println(x) }()  // captures by reference, prints 2 later
x = 2
```

### Panic and Recover

```go
func mustDivide(a, b int) int {
    if b == 0 {
        panic("division by zero")
    }
    return a / b
}

func safeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from: %v", r)
        }
    }()
    return mustDivide(a, b), nil
}
```

`panic` halts execution and unwinds the stack. `recover` (only inside `defer`) catches it.

**Use sparingly.** `panic` is for "this should never happen" bugs. Regular failures use `error`.

`MustX` is a Go convention: panics on failure (e.g., `regexp.MustCompile`). Used when failure indicates a programmer bug.

---

## 15. Concurrency Foundations

### Concurrency vs Parallelism

- **Parallelism** = doing things at literally the same instant on different cores.
- **Concurrency** = structuring your program so multiple tasks can be in progress at once.

Go is fundamentally about concurrency. The runtime makes it parallel automatically when cores are available.

Cooks analogy: parallelism is two cooks; concurrency is one cook juggling several pots.

### Why Not OS Threads?

OS threads are expensive (~1 MB stack each). Switching is slow. Programming them is hard.

Go uses **goroutines** — much lighter (start at ~2 KB stack), managed by Go's runtime, multiplexed onto OS threads (M:N scheduling). You can have millions of goroutines with reasonable memory use.

### Goroutines

```go
go sayHello("Asha")  // schedules and returns immediately
```

The `go` keyword launches a function as a goroutine. The caller doesn't wait. Order of execution is non-deterministic.

```go
go func() {
    // anonymous goroutine
}()
```

### The "Wait For Them" Problem

```go
func main() {
    go fmt.Println("hello")
    // main exits — goroutine may never run
}
```

Solutions: `sync.WaitGroup`, channels, or `time.Sleep` (hack).

### Channels

A typed pipe for communicating between goroutines:

```go
ch := make(chan int)

go func() {
    ch <- 42       // send
}()

v := <-ch          // receive
```

Read the arrow as "data flows in this direction":
- `ch <- v` — send v into ch
- `v := <-ch` — receive into v

### Channels Are Synchronization

Unbuffered channels: every send blocks until a receive happens. Both goroutines synchronize at that moment.

This is the entire point — channels do **data transfer + coordination** with one mechanism.

### Buffered vs Unbuffered

```go
ch := make(chan int)     // unbuffered, capacity 0
ch := make(chan int, 3)  // buffered, capacity 3
```

**Unbuffered:** every send blocks until a receive. Tight synchronization, immediate backpressure.

**Buffered:** sends complete immediately if there's room. Decouples timing, smooths bursts.

When to use which:
- **Unbuffered (default):** when sender/receiver should synchronize.
- **Small buffer (1-16):** to smooth out variable timing.
- **Avoid huge buffers** — they hide bottlenecks.

### Closing Channels

```go
close(ch)
```

A closed channel:
- Cannot have anything sent on it (panics).
- Returns the zero value forever when received from.

```go
v, ok := <-ch  // ok=false once closed and drained
```

**Convention:** only the **sender** closes the channel. Never the receiver. Receiver-closing causes the sender's next send to panic.

For multiple senders, none of them closes — use a coordinator (WaitGroup + closer goroutine).

### `range` on a Channel

```go
for v := range ch {
    process(v)
}
```

Receives until the channel is closed and drained.

### Channels Can Be Function Parameters

```go
func doWork(jobs chan int, done chan bool) {
    // ...
    done <- true
}
```

Channels are reference-like — copying the variable copies the pointer to the underlying channel struct, both refer to the same channel.

### Directional Channels

```go
func produce() <-chan int { ... }    // returns receive-only
func square(in <-chan int) <-chan int  // takes receive-only, returns receive-only
func consume(out chan<- int) { ... }  // takes send-only
```

Documents intent and prevents misuse. Inside the goroutine that owns the channel, it's bidirectional; outside, it's unidirectional.

### `select`

```go
select {
case v := <-ch1:
    // got from ch1
case v := <-ch2:
    // got from ch2
case <-time.After(1 * time.Second):
    // timeout
case ch3 <- value:
    // sent to ch3
default:
    // no case ready (non-blocking)
}
```

`select` blocks until one case is ready. With multiple ready, picks randomly. Default makes it non-blocking.

### WaitGroup

```go
var wg sync.WaitGroup

for i := 0; i < 3; i++ {
    wg.Add(1)              // BEFORE the go statement
    go func(id int) {
        defer wg.Done()    // decrement when done
        // work
    }(i)
}

wg.Wait()  // blocks until counter is zero
```

Three rules:
1. `Add` before `go`.
2. `Done` via `defer`.
3. `Wait` in the parent goroutine.

### Mutex

For shared state that can't fit cleanly into the channel model:

```go
var mu sync.Mutex
var counter int

func increment() {
    mu.Lock()
    defer mu.Unlock()
    counter++
}
```

Without the mutex, two goroutines could both read 5, both write 6, losing an increment (a **race condition**).

### RWMutex

Many readers, one writer:

```go
var mu sync.RWMutex

func get(key string) int {
    mu.RLock()
    defer mu.RUnlock()
    return data[key]
}

func set(key string, v int) {
    mu.Lock()
    defer mu.Unlock()
    data[key] = v
}
```

Use when reads vastly outnumber writes.

### Atomic

For simple counters, faster than a mutex:

```go
var counter int64
atomic.AddInt64(&counter, 1)
v := atomic.LoadInt64(&counter)
```

### Channels vs Mutexes

- **Channels:** when goroutines need to coordinate or pass data.
- **Mutexes:** when goroutines need to share state.

Or: *pass data with channels, protect data with mutexes.*

### The Race Detector

```bash
go run -race .
go test -race
```

Detects data races at runtime. Always use it when writing concurrent code.

---

## 16. The `context` Package

### Why Context Exists

Imagine an HTTP server. A request comes in; you start a goroutine to handle it. Then the user closes their browser. Your goroutine is still working — fetching data, calling APIs, wasting resources for a request that no one will read.

You need a way to tell the goroutine: "stop, this work isn't needed anymore."

That's the **cancellation problem**. Plus the related **deadline problem** ("give up if it takes more than 2 seconds").

### Why Channels Alone Aren't Enough

You can build cancellation with raw channels:

```go
stop := make(chan struct{})
go doWork(stop)
// later:
close(stop)
```

But it gets unwieldy:
1. Multiple goroutines, one cancellation — pass channel down many levels.
2. Multiple cancellation reasons (deadline, user, error).
3. Carrying request-scoped data (request ID, auth, trace span).
4. No standard convention — every library invents its own.

Context wraps the channel pattern in a standard interface for all four problems.

### The Interface

```go
type Context interface {
    Done() <-chan struct{}
    Err() error
    Deadline() (deadline time.Time, ok bool)
    Value(key any) any
}
```

Four methods:
- `Done()` — returns a channel that closes when canceled.
- `Err()` — returns why it was canceled (`context.Canceled` or `context.DeadlineExceeded`).
- `Deadline()` — when it'll auto-cancel, if any.
- `Value(key)` — retrieve attached data.

It's an **interface**, not a struct. You don't implement it yourself; the standard library provides all the implementations. Pass by value (just `ctx context.Context`) — interfaces are already 16-byte handles.

### Foundation Contexts

```go
context.Background()  // root, never cancels, no values
context.TODO()        // placeholder during refactoring
```

Useless on their own. Always derive children from them.

### The Tree Idea

Contexts form a tree. Cancellation flows **down** the tree.

```
                  context.Background()
                          │
             ┌────────────┼────────────┐
             │            │            │
       WithCancel    WithTimeout    WithValue
             │
        WithCancel
```

Cancel any node → all its descendants are canceled. Parent and other branches are unaffected.

### `context.WithCancel`

```go
ctx, cancel := context.WithCancel(parent)
defer cancel()  // ALWAYS

go func() {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            // work
        }
    }
}()

cancel()  // signals all goroutines watching ctx.Done()
```

### `context.WithTimeout` and `context.WithDeadline`

```go
ctx, cancel := context.WithTimeout(parent, 2*time.Second)
defer cancel()

ctx, cancel := context.WithDeadline(parent, deadlineTime)
defer cancel()
```

Auto-cancels after duration / at time. `defer cancel()` is still required even with timeouts (releases the timer if you finish early).

### Handling Context in Functions

```go
func slowOperation(ctx context.Context) (string, error) {
    select {
    case <-time.After(2 * time.Second):
        return "data", nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
```

This is the SHAPE every I/O operation in Go follows: race the work against `ctx.Done()`, return whichever comes first.

### Error Classification

```go
switch {
case errors.Is(err, context.Canceled):
    // cancel() was called
case errors.Is(err, context.DeadlineExceeded):
    // timeout/deadline fired
}
```

### Conventions

1. **`ctx context.Context` is the FIRST parameter.** Always.
2. **Pass by value** (`ctx context.Context`), never `*context.Context`.
3. **Always `defer cancel()`** when you call `WithCancel`, `WithTimeout`, or `WithDeadline`.
4. **Don't store context in structs** — pass it through function calls.

### `context.WithValue` (Use Sparingly)

```go
type contextKey int
const requestIDKey contextKey = 0

ctx = context.WithValue(ctx, requestIDKey, "req-12345")

// Later:
id, ok := ctx.Value(requestIDKey).(string)
```

Conventions:
- Use an **unexported custom type** as the key (prevents collisions across packages).
- Use only for cross-cutting concerns (request IDs, trace spans, auth) — NOT real arguments.

---

## 17. Concurrency Patterns

### Pattern 1: Pipeline

Chain stages connected by channels. Each stage is a goroutine.

```go
func produce() <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := 0; i < 5; i++ {
            out <- i
        }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for v := range in {
            out <- v * v
        }
    }()
    return out
}

// Usage:
nums := produce()
squares := square(nums)
for v := range squares {
    fmt.Println(v)
}
```

**Why pipelines work** (the three real wins):
1. **Memory** — only a few items in flight, never all loaded.
2. **Latency** — first result appears quickly.
3. **Throughput** — overlapping I/O and CPU keeps every stage busy.

The pipeline's total time is bounded by the **slowest stage**, not the sum of all stages.

### Pattern 2: Fan-Out

Distribute one stream of work across N parallel workers.

```go
func fanOut(in <-chan int, numWorkers int, work func(int) int) []<-chan int {
    outs := make([]<-chan int, numWorkers)
    for i := 0; i < numWorkers; i++ {
        out := make(chan int)
        outs[i] = out
        go func() {
            defer close(out)
            for v := range in {
                out <- work(v)
            }
        }()
    }
    return outs
}
```

The trick: **multiple goroutines can read from the same channel.** Each value goes to exactly one of them, chosen by the runtime — natural load balancing.

### Pattern 3: Fan-In

Merge N input channels into one output channel.

```go
func fanIn(inputs ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    for _, ch := range inputs {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(ch)
    }
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

Multiple goroutines write to one channel; values interleave. The closer goroutine waits for all forwarders, then closes.

### Pattern 4: Worker Pool (Bounded Concurrency)

```go
func WorkerPool(ctx context.Context, numWorkers int, jobs <-chan int, work func(int) int) <-chan int {
    results := make(chan int)
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for {
                select {
                case <-ctx.Done():
                    return
                case job, ok := <-jobs:
                    if !ok { return }
                    select {
                    case results <- work(job):
                    case <-ctx.Done():
                        return
                    }
                }
            }
        }()
    }
    go func() {
        wg.Wait()
        close(results)
    }()
    return results
}
```

Why bounded: rate-limiting access to finite resources (DB connections, API limits, file descriptors). Without bounds, 10,000 concurrent goroutines would overwhelm the resource.

### Pattern 5: Pub/Sub (Broadcast)

Channels alone don't broadcast — each value goes to exactly one receiver. To broadcast, you need a subscriber list and explicit forwarding:

```go
type PubSub struct {
    mu          sync.RWMutex
    subscribers []chan string
    closed      bool
}

func (p *PubSub) Subscribe() <-chan string {
    p.mu.Lock()
    defer p.mu.Unlock()
    ch := make(chan string, 8)
    p.subscribers = append(p.subscribers, ch)
    return ch
}

func (p *PubSub) Publish(msg string) {
    p.mu.RLock()
    defer p.mu.RUnlock()
    for _, ch := range p.subscribers {
        select {
        case ch <- msg:
        default:
            // drop on full
        }
    }
}

func (p *PubSub) Close() {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.closed = true
    for _, ch := range p.subscribers {
        close(ch)
    }
}
```

Three policies for slow subscribers:
1. **Drop on full** (above) — publisher never blocks.
2. **Block on full** — slow subscriber slows everyone.
3. **Disconnect slow subscribers** — drop the channel entirely.

For external systems, use a real broker (NATS, Redis pub/sub, Kafka). The Go pattern is for **in-process** broadcasting only.

### The Standard Combo

```
producer → fan-out → [N workers] → fan-in → consumer
```

This is the canonical "process a stream concurrently" shape. Speedup is roughly N× (bounded by CPU cores for CPU-bound work, much higher for I/O-bound work).

### Avoiding Deadlocks and Leaks

**Deadlocks** — every goroutine waiting for someone else.
- Forgetting to `close` a channel (range never exits)
- Closing a channel from a receiver (sender's next send panics)
- Buffered channel filling up while no one reads

**Goroutine leaks** — a goroutine blocked forever.

Fix for both: **always have a clear exit story for every goroutine.** Either it processes a finite stream, or it watches `ctx.Done()`.

Race detector catches data races but NOT deadlocks/leaks.

---

## 18. Channels Internals — Deep Dive

### What `make(chan int)` Actually Does

```go
ch := make(chan int)
```

Two things happen:
1. The runtime allocates a **channel struct on the heap**.
2. `ch` is set to a **pointer** to that struct.

So `ch` is internally an 8-byte pointer (`*hchan`). Not a struct. Just a pointer.

### Inside the Channel Struct

```go
type hchan struct {
    qcount   uint           // items currently in buffer
    dataqsiz uint           // buffer capacity
    buf      unsafe.Pointer // circular buffer
    elemsize uint16
    closed   uint32         // 0 if open, 1 if closed
    sendx    uint           // next send index
    recvx    uint           // next receive index
    recvq    waitq          // goroutines waiting to receive
    sendq    waitq          // goroutines waiting to send
    lock     mutex          // protects everything above
}
```

A channel is a small concurrent queue with two waitlists. The "blocking" you've been using is implemented by parking goroutines on these waitlists and waking them up when the other side arrives.

### What Happens on a Send

```
ch <- 5
```

1. Acquire the channel's lock.
2. Check if closed → if yes, panic.
3. Is there a goroutine in `recvq`? → hand value directly, wake them up. Done.
4. Otherwise, room in buffer? → copy into buffer. Done.
5. Otherwise → park current goroutine on `sendq`, release lock, sleep until woken.

Receive does the symmetric thing.

### Why Channels Are Pointers

If channels were value types (full structs in your variables), copying one would copy buffer, lock, waitlists — yielding two independent channels. Sending into one wouldn't be visible from the other. Useless.

Reference semantics are baked in. **There's no value-vs-pointer choice for channels** — they're already pointers under the hood.

### Q: Pass channels by value or pointer?

Just `chan T`, always. Never `*chan T`. The plain channel variable is already a pointer to the struct.

### Q: Can channels be returned from functions?

Yes — extremely common. The pattern:

```go
func produce() <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        // ... send into out ...
    }()
    return out
}
```

The function creates the channel, starts a background goroutine that fills it, and returns the channel handle. The caller receives.

### Q: Can the consumer close the channel?

Technically yes, but **don't**. Sending to a closed channel **panics**. If the consumer closes while the producer still has items to send, the next send crashes the program. Multi-sender case is even worse.

The convention: **only the sender closes**. Multi-sender setups use a coordinator (WaitGroup + closer goroutine).

To stop a producer early, signal through context, not by closing its output channel.

### Q: How does fan-out's slice reassignment work?

```go
for i := 0; i < numWorkers; i++ {
    out := make(chan int)  // fresh struct each iteration
    outs[i] = out          // copy pointer into slice slot
    go func() { ... }()
}
```

Each iteration creates a **new** channel struct. `out` is reassigned to the new pointer. The previous pointer value stays in `outs[i-1]` — reassigning `out` doesn't affect what the slice already holds.

After the loop, the slice has N pointers to N distinct channel structs.

### Pipeline Trace (Step by Step)

For `produce → square → consume`:

```
1.  main: calling produce()
2.  main: produce() returned (channel ready)
3.  main: calling square(nums)
4.  main: square() returned (channel ready)
5.  main: starting consumption loop
6.  produce: trying to send 0
7.  produce: sent 0, looping
8.  produce: trying to send 1
9.  square: got 0, computing 0*0=0
10. square: trying to send 0
11. square: sent 0, looping
12. square: got 1, computing 1*1=1
13. square: trying to send 1
14. produce: sent 1, looping
15. produce: trying to send 2
16. main: received 0, printing it
17. main: received 1, printing it
...
```

Note: `produce()` and `square()` return immediately — only the channels are returned, not data. The goroutines they launched continue running in the background. **Blocking only happens on channel send/receive operations**, not on function calls.

### The Close Cascade

When a producer's loop ends:
1. Goroutine reaches end of function.
2. Deferred `close(out)` fires.
3. Next stage's `range in` exits cleanly.
4. That stage's deferred `close(out)` fires.
5. Cascades all the way through.

This is why `defer close()` at the top of a producer goroutine is the key to clean termination.

### Q: Buffered or unbuffered channels in pipelines?

Both work. The choice depends on:

**Unbuffered:** strict synchronization, immediate backpressure, makes problems visible.

**Small buffered (1-16):** smooths transient mismatches in stage timing, often improves throughput in real-world pipelines.

**Large buffered:** usually a smell. Hides bottlenecks. Consumes memory.

Production rule of thumb: start unbuffered, add a small buffer when you measure a benefit. Avoid huge buffers.

### Q: Why does the pipeline pattern actually help?

It's not magic speedup — same total work happens. The win is **overlap**: while stage A reads from disk (slow I/O), stage B parses the previous chunk (CPU work), and stage C writes the chunk before that (slow I/O). All three resources are busy at once instead of sequentially.

The ceiling on speedup is determined by the slowest stage. To break that ceiling, fan-out the slow stage across multiple workers.

Real-world pipelines: log/data processing, image/video processing, web crawling, ETL, streaming responses, compression/encryption.

### Q: When NOT to use a pipeline?

- Few items, all fitting in memory — just process them.
- One stage dominates by 100x — fan-out the slow stage instead.
- Stages need each other's full output (sorting, dedup, aggregation).
- Work is trivially fast (microseconds) — channel overhead dominates.

---

## 19. Quick Reference Cards

### The "Always" Rules

```
✓ Always check err != nil
✓ Always defer cancel() after WithCancel/WithTimeout/WithDeadline
✓ Always defer Unlock() after Lock()
✓ Always defer wg.Done() inside goroutine
✓ Always defer close(out) for producer goroutines
✓ Always wg.Add(1) BEFORE the go statement
✓ Always pass ctx as the first parameter
✓ Always use errors.Is/As, not == for wrapped errors
✓ Always use range for runes, indexing for bytes
✓ Always run with -race during development
```

### The "Never" Rules

```
✗ Never close a channel from the receiver
✗ Never close a channel that has multiple senders (without coordination)
✗ Never send on a closed channel
✗ Never use *context.Context (just context.Context)
✗ Never store context in a struct
✗ Never reach for any/interface{} when generics fit
✗ Never panic for normal failures (use error)
✗ Never reach for goroutines without an exit story
✗ Never assume map iteration order is stable
```

### Channel Operation Cheat Sheet

```go
ch := make(chan int)           // unbuffered
ch := make(chan int, 5)        // buffered, capacity 5

ch <- v                        // send (blocks if no receiver / full)
v := <-ch                      // receive (blocks if no sender / empty)
v, ok := <-ch                  // ok=false if closed and drained

close(ch)                      // signal "no more values"
                               // sender closes; never receiver
                               // sending after close → PANIC

for v := range ch {            // receive until closed
    ...
}

select {                       // multiplex
case v := <-ch1:
case ch2 <- value:
case <-time.After(d):          // timeout
case <-ctx.Done():             // cancellation
default:                       // non-blocking
}
```

### Error Handling Patterns

```go
// Create
err := errors.New("simple")
err := fmt.Errorf("formatted: %s", x)
err := fmt.Errorf("with context: %w", innerErr)  // wrapping

// Check
if err != nil { return err }
if errors.Is(err, sentinelErr) { ... }
var typedErr *MyError
if errors.As(err, &typedErr) { ... }

// Custom type
type MyError struct { Field string }
func (e *MyError) Error() string { return "..." }

// Sentinel
var ErrNotFound = errors.New("not found")
```

### Goroutine + Context Skeleton

```go
func doWork(ctx context.Context, jobs <-chan Job) {
    for {
        select {
        case <-ctx.Done():
            return
        case job, ok := <-jobs:
            if !ok { return }
            // process job
        }
    }
}

// Caller:
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
go doWork(ctx, jobs)
```

### WaitGroup Skeleton

```go
var wg sync.WaitGroup
for _, item := range items {
    wg.Add(1)
    go func(item Item) {
        defer wg.Done()
        process(item)
    }(item)
}
wg.Wait()
```

### Mutex Skeleton

```go
type Counter struct {
    mu sync.Mutex
    n  int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.n++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.n
}
```

### Pipeline Skeleton

```go
func stage1() <-chan T {
    out := make(chan T)
    go func() {
        defer close(out)
        for ... {
            out <- value
        }
    }()
    return out
}

func stage2(in <-chan T) <-chan U {
    out := make(chan U)
    go func() {
        defer close(out)
        for v := range in {
            out <- transform(v)
        }
    }()
    return out
}

// Use:
ch := stage2(stage1())
for v := range ch { ... }
```

### Worker Pool Skeleton

```go
func WorkerPool(ctx context.Context, n int, jobs <-chan Job) <-chan Result {
    results := make(chan Result)
    var wg sync.WaitGroup
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for {
                select {
                case <-ctx.Done():
                    return
                case job, ok := <-jobs:
                    if !ok { return }
                    select {
                    case results <- process(job):
                    case <-ctx.Done():
                        return
                    }
                }
            }
        }()
    }
    go func() { wg.Wait(); close(results) }()
    return results
}
```

### Generic Function Skeleton

```go
import "cmp"

func Max[T cmp.Ordered](a, b T) T {
    if a > b { return a }
    return b
}

func MapSlice[T, U any](in []T, f func(T) U) []U {
    out := make([]U, len(in))
    for i, v := range in { out[i] = f(v) }
    return out
}

type Stack[T any] struct { items []T }
func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }
func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 { var zero T; return zero, false }
    v := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return v, true
}
```

### Mental Models in One-Liners

- **Strings** = byte sequences (UTF-8). Use `range` for runes, `len` for bytes.
- **Slices** = reference to a backing array. Header copied on pass; data shared.
- **Maps** = pointer to hash table. Reference semantics. Iteration order randomized.
- **Channels** = pointer to a coordination struct. Reference semantics. Sender closes.
- **Interfaces** = (type info, data pointer). Method dispatch via type info lookup.
- **Goroutines** = lightweight tasks (~2 KB stack). Multiplexed onto OS threads.
- **Context** = standard cancellation/deadline propagation. Tree structure. Cancel flows down.
- **Pointers** = address of a value. Take with `&`, follow with `*`. Auto-deref for `.` on structs.
- **Errors** = values, not exceptions. Return `(value, error)`. Always check.
- **Generics** = compile-time type substitution. Constraints determine allowed types and operations.

### File Organization

A typical Go project:

```
project/
├── go.mod                 module declaration
├── go.sum                 dependency hashes
├── main.go                package main
└── pkgname/
    └── pkgname.go         package pkgname (lowercase, single word)
```

Run: `go run .`
Build: `go build ./...`
Vet: `go vet ./...`
Test with race: `go test -race ./...`
Escape analysis: `go build -gcflags="-m" ./...`

### The 5 Big Concurrency Patterns

```
PIPELINE     → chain stages via channels (stage1 → stage2 → stage3)
FAN-OUT      → many workers reading one channel (load distribution)
FAN-IN       → many channels merged into one (combine results)
WORKER POOL  → bounded fan-out for queue-driven work
PUB/SUB      → broadcast (one publisher → all subscribers)
```

Composition: `producer → fan-out → [workers] → fan-in → consumer` is the canonical concurrent stream-processing shape.

---

## End

This document covers everything from "what is a type" through advanced concurrency patterns and the internal mechanics of channels and interfaces. Refer back as needed; the patterns become intuitive only with practice.

The reference packages built alongside these notes (in separate directories) contain runnable code for each topic — `customstrings`, `loopsref`, `funcsref`, `pointersref`, `interfacesref`, `genericsref`, `errorsref`, `concurrencyref`, `contextref`, `concpatterns`. Each compiles, vets cleanly, and was tested with the race detector where applicable.
