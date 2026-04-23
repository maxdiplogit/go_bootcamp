// Package loopsref is a personal reference package for loops in Go.
//
// Go has only ONE looping keyword -- `for` -- and it has FOUR shapes that
// together cover every looping pattern other languages split across `while`,
// `do-while`, `for`, `foreach`, `for...of`, etc.
//
// THE FOUR SHAPES (memorize these):
//
//  1. Three-part:       for init; cond; post { ... }   // classic counted loop
//  2. Condition-only:   for cond { ... }               // Go's `while`
//  3. Infinite:         for { ... }                    // Go's `while true`
//  4. Range-based:      for i, v := range collection   // foreach style
//
// Plus two control statements: `break` (exit innermost loop) and `continue`
// (skip to next iteration). Both can target a labeled loop to escape nested
// loops cleanly.
//
// This package collects small, runnable examples of each pattern, with
// explanatory comments. Read top-to-bottom; each function is independent
// so you can reorder or pick-and-choose freely from main.go.
package loopsref

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// =============================================================================
// SECTION 1: THE FOUR SHAPES OF `for`
// =============================================================================

// CountUp uses the classic three-part for. This is the workhorse for any
// loop where you know how many iterations you need up front.
//
// Mental execution:
//  1. `i := 0` runs ONCE before the loop starts.
//  2. Check `i < n`. If true, run body. If false, exit.
//  3. After the body, run `i++`. Go back to step 2.
//
// Note: `i` is scoped to the loop. It does not exist outside the for block.
func CountUp(n int) []int {
	out := make([]int, 0, n) // pre-allocate to avoid append reallocations
	for i := 0; i < n; i++ {
		out = append(out, i)
	}
	return out
}

// CountDown shows that the three parts are flexible -- you can decrement,
// you can step by any amount, and you can use any comparison. There is no
// "right way" to write the parts beyond "they must make sense together."
func CountDown(n int) []int {
	out := make([]int, 0, n)
	for i := n; i > 0; i-- {
		out = append(out, i)
	}
	return out
}

// StepBy demonstrates a custom step. The post-statement isn't limited to
// `i++`; it can be any statement that updates the loop variable.
func StepBy(start, end, step int) []int {
	var out []int
	for i := start; i < end; i += step {
		out = append(out, i)
	}
	return out
}

// WhileLessThan uses the CONDITION-ONLY form -- Go's equivalent of `while`.
// Reach for this when you don't know up front how many iterations you need:
// you're waiting for some condition to become true (or false).
//
// Here we keep doubling until we cross a threshold. The number of
// iterations isn't predictable from the inputs without doing the math.
func WhileLessThan(start, limit int) int {
	n := start
	for n < limit {
		n *= 2
	}
	return n
}

// InfiniteWithBreak uses the BARE `for {}` -- Go's equivalent of
// `while true`. The loop has no built-in exit; you must `break`, `return`,
// or otherwise escape from inside the body.
//
// This shape is everywhere in real Go code: servers, event loops, and
// retry loops where the exit logic is too complex to fit in a condition.
func InfiniteWithBreak(target int) int {
	i := 0
	for {
		if i >= target {
			break // exit the innermost enclosing loop
		}
		i++
	}
	return i
}

// =============================================================================
// SECTION 2: BREAK AND CONTINUE
// =============================================================================

// SkipAndStop shows both control statements in one loop:
//   - `continue` skips the REST of the current iteration and goes to the next.
//   - `break` exits the loop entirely.
//
// Walk through it with n=10, skip=3, stop=7:
//
//	i=0,1,2 -> appended
//	i=3     -> continue, NOT appended
//	i=4,5,6 -> appended
//	i=7     -> break, loop exits
//
// Result: [0 1 2 4 5 6]
func SkipAndStop(n, skip, stop int) []int {
	var out []int
	for i := 0; i < n; i++ {
		if i == skip {
			continue
		}
		if i == stop {
			break
		}
		out = append(out, i)
	}
	return out
}

// FilterEven uses `continue` as a clean way to skip unwanted items.
// This style reads more naturally than wrapping the real work inside an
// `if v%2 == 0 { ... }`, especially when the filtering logic gets complex.
func FilterEven(nums []int) []int {
	var evens []int
	for _, v := range nums {
		if v%2 != 0 {
			continue // not even -- skip
		}
		evens = append(evens, v)
	}
	return evens
}

// FindFirst uses `break` to stop as soon as a match is found.
// Returns (index, true) on success or (-1, false) if no match.
//
// This is the canonical "search" pattern: scan, stop early, report.
func FindFirst(nums []int, target int) (int, bool) {
	for i, v := range nums {
		if v == target {
			return i, true // an early return is a clean alternative to break + flag
		}
	}
	return -1, false
}

// =============================================================================
// SECTION 3: RANGE OVER SLICES
// =============================================================================

// SumSlice ranges over a slice and accumulates a total. Use the blank
// identifier `_` to discard the index when you don't need it.
//
// `v` is a COPY of each element. For ints (and other small values) this
// is essentially free. For large structs you'd index instead to avoid
// copying.
func SumSlice(nums []int) int {
	total := 0
	for _, v := range nums {
		total += v
	}
	return total
}

// DoubleInPlace demonstrates the ONE THING you cannot do with a value
// from `range`: mutate the original collection. The value `v` is a copy;
// changing it doesn't change the slice.
//
// To actually mutate the slice, range over INDICES and write back through
// the index. That's what we do here.
func DoubleInPlace(nums []int) {
	for i := range nums {
		nums[i] *= 2
	}
	// Note: we don't need to return -- slices are reference-like, so the
	// caller sees our changes. (More on this when we cover slices formally.)
}

// MaxValue scans for the largest value. Pattern: initialize from the first
// element, then range over the REST and update if anything beats it.
//
// We use nums[1:] (a slice of everything from index 1 onward) so we don't
// re-compare the first element to itself. (Slices come up next in the roadmap.)
func MaxValue(nums []int) (int, bool) {
	if len(nums) == 0 {
		return 0, false // no max for an empty slice
	}
	max := nums[0]
	for _, v := range nums[1:] {
		if v > max {
			max = v
		}
	}
	return max, true
}

// =============================================================================
// SECTION 4: RANGE OVER STRINGS (BYTES vs RUNES)
// =============================================================================

// IterateBytes walks a string by BYTE INDEX -- s[i] is a byte (uint8).
// For pure ASCII, one byte == one character. For anything else, you'll
// see "broken" bytes when iterating this way -- this is the wrong approach
// for Unicode-aware code.
//
// Compare this with IterateRunes below to see the difference clearly.
func IterateBytes(s string) []byte {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		out = append(out, s[i]) // s[i] is type byte
	}
	return out
}

// IterateRunes walks a string with `range`, which DECODES UTF-8 for you.
// Each iteration yields:
//   - i: the BYTE OFFSET where this rune starts (NOT a sequential counter)
//   - r: the rune (int32 code point)
//
// The byte offset jumps by the size of each rune in UTF-8 (1, 2, 3, or 4).
// This is the safe, default way to walk a string when characters matter.
func IterateRunes(s string) []rune {
	out := make([]rune, 0, len(s)) // upper-bound capacity
	for _, r := range s {
		out = append(out, r)
	}
	return out
}

// CountVowels uses range + unicode.ToLower to count vowels in a Unicode-safe
// way. Note the `switch` inside the loop -- a clean way to test against
// multiple values without writing a long `||` chain.
func CountVowels(s string) int {
	count := 0
	for _, r := range s {
		switch unicode.ToLower(r) {
		case 'a', 'e', 'i', 'o', 'u':
			count++
		}
	}
	return count
}

// =============================================================================
// SECTION 5: RANGE OVER MAPS
// =============================================================================

// SumMapValues iterates a map with `range`. For maps, `range` yields
// (key, value) pairs.
//
// IMPORTANT: map iteration order is RANDOMIZED. Each pass over the same
// map may visit keys in a different order. This is intentional: it
// prevents code from accidentally depending on order, which would create
// brittle programs that break on Go runtime upgrades.
func SumMapValues(m map[string]int) int {
	total := 0
	for _, v := range m {
		total += v
	}
	return total
}

// MapKeys returns a slice of the map's keys. Useful when you want to do
// something with all the keys, like sort them.
//
// We pre-allocate `make([]string, 0, len(m))` because we know exactly
// how many we'll append -- one per map entry.
func MapKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m { // single variable form -- only the key
		keys = append(keys, k)
	}
	return keys
}

// SortedMapEntries shows the canonical "iterate map in sorted order" idiom:
//  1. Collect the keys into a slice.
//  2. Sort the slice.
//  3. Iterate the sorted slice, looking up each value as you go.
//
// You'll write this pattern many times in real Go code, because map iteration
// order is random and a great many problems need a stable, ordered view.
func SortedMapEntries(m map[string]int) []string {
	keys := MapKeys(m)
	sort.Strings(keys)

	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, fmt.Sprintf("%s=%d", k, m[k]))
	}
	return out
}

// =============================================================================
// SECTION 6: NESTED LOOPS
// =============================================================================

// BuildGrid returns rows*cols cells as a 2D slice, where cell[i][j] = i*j.
// Demonstrates a basic nested loop -- the inner loop runs to completion
// for every iteration of the outer loop.
//
// Pre-allocation matters more for nested loops: we make the outer slice
// with len=rows up front, and each inner slice with len=cols. No appends
// needed, no reallocations.
func BuildGrid(rows, cols int) [][]int {
	grid := make([][]int, rows)
	for i := 0; i < rows; i++ {
		grid[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			grid[i][j] = i * j
		}
	}
	return grid
}

// AllPairs returns every (a, b) pair from the input slice where a != b.
// Useful for "compare every element to every other" problems -- O(n^2).
//
// Note we use j := i+1 (not 0) to avoid both:
//   - comparing an element with itself
//   - producing each pair twice (e.g. (1,2) and (2,1))
func AllPairs(nums []int) [][2]int {
	var pairs [][2]int
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			pairs = append(pairs, [2]int{nums[i], nums[j]})
		}
	}
	return pairs
}

// =============================================================================
// SECTION 7: LABELS -- BREAKING OUT OF NESTED LOOPS
// =============================================================================

// FindInGrid searches a 2D grid for a target value. Returns the (row, col)
// of the first occurrence, or (-1, -1) if not found.
//
// THE KEY POINT: a plain `break` inside the inner loop only exits the
// inner loop. The outer loop would keep going. To exit BOTH loops at once,
// we attach a LABEL to the outer loop and use `break outer`.
//
// (In this specific function we could use `return` instead, since we know
// the answer the moment we find it. But the label form is what you'd use
// when there's more work to do after the search.)
func FindInGrid(grid [][]int, target int) (int, int) {
	row, col := -1, -1

outer: // label -- attached to the for statement that immediately follows
	for i, r := range grid {
		for j, v := range r {
			if v == target {
				row, col = i, j
				break outer // exits BOTH loops, jumping past the labeled one
			}
		}
	}
	return row, col
}

// ContinueOuter shows `continue` targeting a label.
// `continue outer` jumps to the next iteration of the OUTER loop, skipping
// the rest of the current outer iteration AND any remaining inner work.
//
// Result for n=4: only pairs where j <= i are produced.
func ContinueOuter(n int) [][2]int {
	var out [][2]int

outer:
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > i {
				continue outer // skip to next i, abandon this inner loop
			}
			out = append(out, [2]int{i, j})
		}
	}
	return out
}

// =============================================================================
// SECTION 8: COMMON ITERATION PATTERNS (the ones you'll write 1000 times)
// =============================================================================

// Filter returns a new slice containing only the elements for which keep(v)
// is true. The classic "predicate filter" pattern.
//
// We pre-allocate with len=0 and cap=len(in), since the result can't be
// bigger than the input. This avoids append reallocations entirely.
func Filter(in []int, keep func(int) bool) []int {
	out := make([]int, 0, len(in))
	for _, v := range in {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

// MapInts applies a transformation function to every element. Output has
// the same length as input, so we can pre-size it exactly and assign by
// index instead of appending -- a tiny but real efficiency win.
func MapInts(in []int, f func(int) int) []int {
	out := make([]int, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}

// Reduce collapses a slice into a single value using an accumulator function.
// Many problems are reduces in disguise: sum, product, max, min, concat.
func Reduce(in []int, initial int, f func(acc, v int) int) int {
	acc := initial
	for _, v := range in {
		acc = f(acc, v)
	}
	return acc
}

// GroupBy groups items by a key function. Returns a map from each key to
// the list of items that produced it.
//
// Notice the natural use of `append` to a map value -- Go's append handles
// the nil-slice case correctly, so we don't need to special-case "first
// time we've seen this key".
func GroupBy(words []string, key func(string) string) map[string][]string {
	groups := map[string][]string{}
	for _, w := range words {
		k := key(w)
		groups[k] = append(groups[k], w)
	}
	return groups
}

// =============================================================================
// SECTION 9: A FULL, REALISTIC EXAMPLE
// =============================================================================

// WordFrequency counts how many times each unique word appears in the
// input text. This is a tiny but realistic program that uses:
//   - strings.Fields (splits on any whitespace)
//   - a `range` loop with the index discarded
//   - a map for counting
//   - the lowercase-key trick for case-insensitive grouping
//   - the "collect keys, sort, iterate" idiom for stable output
//
// This is the kind of utility you might actually write at work. Real Go
// code is full of small loops + maps doing exactly this kind of work.
func WordFrequency(text string) []WordCount {
	counts := map[string]int{}

	for _, raw := range strings.Fields(text) {
		// Normalize: lowercase + strip surrounding punctuation so "Hello,"
		// and "hello" count as the same word.
		clean := strings.ToLower(strings.TrimFunc(raw, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r)
		}))
		if clean == "" {
			continue
		}
		counts[clean]++
	}

	// Sort: highest count first, then alphabetical for ties.
	out := make([]WordCount, 0, len(counts))
	for w, c := range counts {
		out = append(out, WordCount{Word: w, Count: c})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Word < out[j].Word
	})
	return out
}

// WordCount pairs a word with its frequency. Exported for use in main.
type WordCount struct {
	Word  string
	Count int
}
