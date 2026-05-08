package main

import (
	"go_revision/arrays"
	"go_revision/contextref"
	customstrings "go_revision/customStrings"
	"go_revision/errorsref"
	"go_revision/funcsref"
	genericsref "go_revision/generics"
	"go_revision/goroutines"
	interfacesref "go_revision/interfaces"
	"go_revision/loopsref"
	"go_revision/make_practice"
	"go_revision/maps"
	"go_revision/pointersref"
	"go_revision/slices"
	"go_revision/structs"
	"go_revision/utils"

	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// section prints a visual separator so each demo block is easy to find in
// the terminal output.
func section(title string) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("  " + title)
	fmt.Println(strings.Repeat("=", 60))
}

// toIntMap converts map[string][]string to map[string]int (using length)
// just so we can reuse loopsref.MapKeys for stable display ordering above.
// This is a tiny helper for the demo, not part of the reference.
func toIntMap(m map[string][]string) map[string]int {
	out := make(map[string]int, len(m))
	for k, v := range m {
		out[k] = len(v)
	}
	return out
}

func main() {
	var x, y int

	fmt.Print("Enter x: ")
	fmt.Scanf("%d", &x)

	fmt.Print("Enter y: ")
	fmt.Scanf("%d", &y)

	var res int = utils.AddTwoNumbers(x, y)

	fmt.Printf("Result: %d\n", res)

	utils.PrintStars()
	utils.PrintDay()

	var a, b int = utils.ReturnMultipleValues()
	fmt.Printf("a = %d; b = %d\n", a, b)

	var vehicle_0 *structs.Vehicle
	vehicle_0 = &structs.Vehicle{
		Company:      "BMW",
		FuelType:     "Petrol",
		SerialNumber: 12,
	}

	vehicle_0.PrintVehicleData()
	vehicle_0.UpdateVehicleData("Audi", "Diesel", 65535)
	vehicle_0.PrintVehicleData()

	arrays.ArrayPractice()

	slices.SlicesPractice()

	make_practice.MakeExample()

	maps.MapsExample()

	truck_0 := &structs.Truck{
		EngineModel: "XA31",
		Tyres:       6,
		Pistons:     8,
	}
	truck_0.TruckJsonExample()

	fmt.Println()

	userIdentifier := &structs.Identifier{
		ID:   1,
		UUID: "user-abc",
	}

	u := structs.User{
		UserIdentifier: userIdentifier,
		Timestamps: structs.Timestamps{
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
		},
		AuditTrail: structs.AuditTrail{
			CreatedBy: 99,
			UpdatedBy: 99,
			Version:   2,
		},
		Email:    "asha@example.com",
		FullName: "Asha Rao",
	}

	// Promoted access — clean and ergonomic:
	fmt.Println(u.UserIdentifier) // from Identifier
	fmt.Println(u.CreatedAt)      // from Timestamps
	fmt.Println(u.Version)        // from AuditTrail
	fmt.Println(u.Email)          // User's own field

	// JSON output is flat — looks like one coherent object:
	data, _ := json.MarshalIndent(u, "", "  ")
	fmt.Println(string(data))

	fmt.Println()

	// A short ASCII string and a string with accents + CJK + emoji.
	// Keep both around -- most demos want to show the contrast.
	asciiStr := "hello"
	mixedStr := "héllo 世界 🙂"

	// -----------------------------------------------------------------
	section("1. What a string is")
	fmt.Println(customstrings.StringHeader())
	fmt.Println(customstrings.StringsAreImmutable())

	// -----------------------------------------------------------------
	section("2. len() counts bytes, not characters")
	for _, s := range []string{asciiStr, mixedStr} {
		b, r := customstrings.LenIsBytes(s)
		fmt.Printf("%-20q bytes=%d runes=%d\n", s, b, r)
	}

	// -----------------------------------------------------------------
	section("3. Indexing gives a byte")
	// For ASCII, s[0] == 'h' == 104.
	// For the mixed string, s[1] is a raw continuation byte of 'é' --
	// it is NOT a usable character on its own.
	for _, s := range []string{asciiStr, mixedStr} {
		b, n, _ := customstrings.IndexingReturnsBytes(s, 1)
		fmt.Printf("%-20q s[1] = byte %d (0x%X)\n", s, n, b)
	}

	// -----------------------------------------------------------------
	section("4. range decodes UTF-8 into runes")
	// Notice that ByteIndex JUMPS for multi-byte characters.
	// 'é' at byte 1 occupies 2 bytes, so the next entry is at byte 3.
	for _, entry := range customstrings.ByteVsRune(mixedStr) {
		fmt.Printf("  byte %2d -> rune %-8U %q (len %d bytes)\n",
			entry.ByteIndex, entry.Rune, entry.Rune, entry.ByteLen)
	}

	// -----------------------------------------------------------------
	section("5. Raw bytes vs rune slice")
	fmt.Printf("bytes of %q -> %v\n", mixedStr, customstrings.RawBytes(mixedStr))
	fmt.Printf("runes of %q -> %v\n", mixedStr, customstrings.RuneSlice(mixedStr))

	// -----------------------------------------------------------------
	section("6. Round-trips through []byte and []rune")
	fmt.Println("BytesRoundTrip (ASCII upper): ", customstrings.BytesRoundTrip("Hello 世界"))
	// ^ only ASCII letters are upper-cased; CJK is left alone -- but
	//   because we're mutating bytes, this is safe only when we limit
	//   ourselves to 1-byte characters, which we do above.

	fmt.Println("RunesRoundTrip (full upper):  ", customstrings.RunesRoundTrip("héllo world"))
	// ^ unicode.ToUpper handles accented characters correctly.

	// -----------------------------------------------------------------
	section("7. Reversing a string the right (and wrong) way")
	fmt.Printf("ReverseASCII(%q)   = %q\n", "hello", customstrings.ReverseASCII("hello"))
	// The next line intentionally shows corruption if you uncomment it:
	// fmt.Printf("ReverseASCII(%q)   = %q\n", mixedStr, customstrings.ReverseASCII(mixedStr))
	fmt.Printf("ReverseUnicode(%q) = %q\n", mixedStr, customstrings.ReverseUnicode(mixedStr))

	// Convert a rune to a 1-character string (idiomatically).
	fmt.Println("RuneToString('A') =", customstrings.RuneToString('A'))
	fmt.Println("RuneToString('世') =", customstrings.RuneToString('世'))

	// -----------------------------------------------------------------
	section("8. Counting correctly")
	fmt.Println("CountBytes:", customstrings.CountBytes(mixedStr))
	fmt.Println("CountRunes:", customstrings.CountRunes(mixedStr))
	fmt.Println("CountWords:", customstrings.CountWords("Hello there, 世界! How are you?"))
	fmt.Println("CountLines:", customstrings.CountLines("line 1\nline 2\nline 3"))

	// -----------------------------------------------------------------
	section("9. Efficient string building")
	parts := []string{"one", "two", "three", "four", "five"}
	fmt.Println("ConcatWithPlus:    ", customstrings.ConcatWithPlus(parts))
	fmt.Println("ConcatWithBuilder: ", customstrings.ConcatWithBuilder(parts))
	fmt.Println("JoinWithSeparator: ", customstrings.JoinWithSeparator(parts, ", "))
	// For 5 items the difference in speed is invisible, but in a tight
	// loop over thousands of items, Builder / Join can be orders of
	// magnitude faster than `+=`.

	// -----------------------------------------------------------------
	section("10. Standard strings operations")
	for k, v := range customstrings.Transformations("Hello world. hello again.") {
		fmt.Printf("  %-18s -> %v\n", k, v)
	}

	// -----------------------------------------------------------------
	section("11. Rune classification (unicode package)")
	for _, r := range []rune{'A', '5', ' ', '!', '世', '٥'} {
		fmt.Printf("  %c (%U):\n", r, r)
		for k, v := range customstrings.ClassifyRune(r) {
			if v {
				fmt.Printf("      %s\n", k)
			}
		}
	}

	// -----------------------------------------------------------------
	section("12. UTF-8 validation and decoding one rune")
	fmt.Println("IsValidUTF8 of mixedStr:", customstrings.IsValidUTF8(mixedStr))
	// A deliberately-broken byte sequence: 0xFF is never valid UTF-8.
	fmt.Println("IsValidUTF8 of 0xFF:    ", customstrings.IsValidUTF8("\xff"))
	r, size := customstrings.FirstRune(mixedStr)
	fmt.Printf("FirstRune of %q = %c (size %d bytes)\n", mixedStr, r, size)

	// -----------------------------------------------------------------
	section("13. Literal forms")
	interp, raw := customstrings.LiteralExamples()
	fmt.Println("Interpreted:")
	fmt.Println(interp)
	fmt.Println("Raw:")
	fmt.Println(raw)

	// -----------------------------------------------------------------
	section("14. Full analyzer")
	stats := customstrings.Analyze("Hello, 世界!\nThis is Go. 🙂")
	fmt.Print(stats.PrettyPrint())

	// -----------------------------------------------------------------
	section("1. The four shapes of `for`")
	fmt.Println("CountUp(5):       ", loopsref.CountUp(5))
	fmt.Println("CountDown(5):     ", loopsref.CountDown(5))
	fmt.Println("StepBy(0,20,3):   ", loopsref.StepBy(0, 20, 3))
	fmt.Println("WhileLessThan:    ", loopsref.WhileLessThan(1, 100))
	fmt.Println("InfiniteWithBreak:", loopsref.InfiniteWithBreak(7))

	// -----------------------------------------------------------------
	section("2. break and continue")
	// SkipAndStop(10, skip=3, stop=7):
	//   visits 0..9, skips 3, breaks at 7
	//   -> 0 1 2 4 5 6
	fmt.Println("SkipAndStop:      ", loopsref.SkipAndStop(10, 3, 7))
	fmt.Println("FilterEven:       ", loopsref.FilterEven([]int{1, 2, 3, 4, 5, 6, 7, 8}))
	idx, ok := loopsref.FindFirst([]int{10, 20, 30, 40}, 30)
	fmt.Println("FindFirst(30):    ", idx, ok)
	idx, ok = loopsref.FindFirst([]int{10, 20, 30, 40}, 99)
	fmt.Println("FindFirst(99):    ", idx, ok)

	// -----------------------------------------------------------------
	section("3. Range over slices")
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	fmt.Println("SumSlice:         ", loopsref.SumSlice(nums))

	// Make a copy so we don't mutate the original we're using elsewhere.
	cp := append([]int{}, nums...)
	loopsref.DoubleInPlace(cp)
	fmt.Println("DoubleInPlace:    ", cp) // notice: cp was modified by the function

	max, ok := loopsref.MaxValue(nums)
	fmt.Println("MaxValue:         ", max, ok)
	emptyMax, ok := loopsref.MaxValue(nil)
	fmt.Println("MaxValue(nil):    ", emptyMax, ok)

	// -----------------------------------------------------------------
	section("4. Range over strings: bytes vs runes")
	mixed := "héllo"
	// IterateBytes returns 6 entries (5-char string + 1 extra byte for é).
	fmt.Printf("IterateBytes(%q): %v (len=%d)\n",
		mixed, loopsref.IterateBytes(mixed), len(loopsref.IterateBytes(mixed)))
	// IterateRunes returns 5 entries -- one per character.
	fmt.Printf("IterateRunes(%q): %v (len=%d)\n",
		mixed, loopsref.IterateRunes(mixed), len(loopsref.IterateRunes(mixed)))
	fmt.Println("CountVowels('héllo World'):", loopsref.CountVowels("héllo World"))

	// -----------------------------------------------------------------
	section("5. Range over maps")
	prices := map[string]int{
		"apple":  50,
		"banana": 20,
		"cherry": 80,
		"date":   30,
	}
	fmt.Println("SumMapValues:    ", loopsref.SumMapValues(prices))
	// MapKeys order is random! Run it twice to see the difference.
	fmt.Println("MapKeys (run 1): ", loopsref.MapKeys(prices))
	fmt.Println("MapKeys (run 2): ", loopsref.MapKeys(prices))
	// SortedMapEntries always gives the same, alphabetical order.
	fmt.Println("SortedMapEntries:", loopsref.SortedMapEntries(prices))

	// -----------------------------------------------------------------
	section("6. Nested loops")
	grid := loopsref.BuildGrid(4, 4)
	fmt.Println("BuildGrid(4,4):")
	for _, row := range grid {
		fmt.Printf("  %v\n", row)
	}
	fmt.Println("AllPairs([1 2 3 4]):", loopsref.AllPairs([]int{1, 2, 3, 4}))

	// -----------------------------------------------------------------
	section("7. Labels: breaking and continuing outer loops")
	row, col := loopsref.FindInGrid(grid, 6)
	fmt.Printf("FindInGrid(grid, 6): row=%d col=%d\n", row, col)
	row, col = loopsref.FindInGrid(grid, 999)
	fmt.Printf("FindInGrid(grid, 999): row=%d col=%d (not found)\n", row, col)
	// ContinueOuter(4): for each i, only emits j values 0..i (lower triangle).
	fmt.Println("ContinueOuter(4):", loopsref.ContinueOuter(4))

	// -----------------------------------------------------------------
	section("8. Common iteration patterns")
	// Filter: keep only positive numbers.
	pos := loopsref.Filter([]int{-3, 0, 5, -2, 8}, func(v int) bool { return v > 0 })
	fmt.Println("Filter positives:", pos)

	// MapInts: square every element.
	squares := loopsref.MapInts([]int{1, 2, 3, 4, 5}, func(v int) int { return v * v })
	fmt.Println("MapInts squares: ", squares)

	// Reduce: product of [1..5] = 120.
	prod := loopsref.Reduce([]int{1, 2, 3, 4, 5}, 1, func(acc, v int) int { return acc * v })
	fmt.Println("Reduce product:  ", prod)

	// GroupBy: group words by their first letter.
	groups := loopsref.GroupBy(
		[]string{"apple", "banana", "blueberry", "avocado", "cherry"},
		func(w string) string { return w[:1] },
	)
	fmt.Println("GroupBy first letter:")
	// (We re-sort keys here only for stable, predictable output.)
	keys := loopsref.MapKeys(toIntMap(groups))
	for _, k := range keys {
		fmt.Printf("  %s -> %v\n", k, groups[k])
	}

	// -----------------------------------------------------------------
	section("9. Realistic example: WordFrequency")
	text := `Go is fun. Go is fast. Go is simple.
Sometimes Go is verbose, but Go is consistent.`
	freq := loopsref.WordFrequency(text)
	for _, wc := range freq {
		fmt.Printf("  %-12s %d\n", wc.Word, wc.Count)
	}

	// -----------------------------------------------------------------
	section("1. Basics")
	fmt.Println("Add(3, 4):       ", funcsref.Add(3, 4))
	fmt.Println("Greet:           ", funcsref.Greet("Asha"))
	fmt.Println("NoArgsNoReturn:  ", funcsref.NoArgsNoReturn())

	// -----------------------------------------------------------------
	section("2. Call by value vs call by pointer")

	// Plain int by value -- the original is untouched.
	x_func := 42
	insideValue := funcsref.FailedReset(x_func)
	fmt.Printf("FailedReset:        x outside=%d, function returned=%d\n", x, insideValue)
	// ^ Notice: the function returned 0, but x is still 42 here.

	// Plain int by pointer -- we pass &x so the function can mutate it.
	funcsref.SuccessfulReset(&x)
	fmt.Printf("SuccessfulReset:    x outside=%d (now zero!)\n", x)

	// Struct by value vs by pointer.
	c := funcsref.CounterStruct{Count: 100}
	funcsref.IncrementValue(c)
	fmt.Printf("IncrementValue:     c.Count=%d (unchanged - struct was copied)\n", c.Count)

	funcsref.IncrementPointer(&c)
	fmt.Printf("IncrementPointer:   c.Count=%d (incremented - same struct)\n", c.Count)

	// Slice -- the surprising case. Header was copied, but the data backs
	// onto the same array, so the function's writes are visible to us.
	func_nums := []int{1, 2, 3, 4}
	funcsref.SliceMutation(func_nums)
	fmt.Printf("SliceMutation:      nums=%v (changed even though we passed by value)\n", nums)

	// -----------------------------------------------------------------
	section("3. Multiple return values")
	q, rem := funcsref.Divide(17, 5)
	fmt.Printf("Divide(17, 5):      quotient=%d remainder=%d\n", q, rem)

	// Use NEW names so := works cleanly:
	sq, sok := funcsref.SafeDivide(20, 4)
	fmt.Printf("SafeDivide(20, 4):  result=%d ok=%v\n", sq, sok)
	sq, sok = funcsref.SafeDivide(20, 0) // plain = now, both already exist
	fmt.Printf("SafeDivide(20, 0):  result=%d ok=%v\n", sq, sok)

	// Discarding values with the blank identifier `_`.
	q2, _ := funcsref.Divide(100, 3) // we don't care about the remainder
	fmt.Printf("Divide w/ _ :       quotient only=%d\n", q2)

	mn, mx, ok := funcsref.MinMax([]int{3, 1, 4, 1, 5, 9, 2, 6})
	fmt.Printf("MinMax:             min=%d max=%d ok=%v\n", mn, mx, ok)
	_, _, ok = funcsref.MinMax(nil)
	fmt.Printf("MinMax(nil):        ok=%v (false because empty input)\n", ok)

	// -----------------------------------------------------------------
	section("4. Named return values")
	first, last := funcsref.SplitName("Asha Rao")
	fmt.Printf("SplitName:          first=%q last=%q\n", first, last)
	first, last = funcsref.SplitName("Cher")
	fmt.Printf("SplitName('Cher'):  first=%q last=%q (last is the zero value \"\")\n", first, last)

	start, end, ok := funcsref.ParseRange("10-20")
	fmt.Printf("ParseRange('10-20'):  start=%d end=%d ok=%v\n", start, end, ok)
	start, end, ok = funcsref.ParseRange("oops")
	fmt.Printf("ParseRange('oops'):   start=%d end=%d ok=%v\n", start, end, ok)

	// -----------------------------------------------------------------
	section("5. Variadic functions")
	fmt.Println("SumAll():            ", funcsref.SumAll())           // 0
	fmt.Println("SumAll(7):           ", funcsref.SumAll(7))          // 7
	fmt.Println("SumAll(1, 2, 3, 4):  ", funcsref.SumAll(1, 2, 3, 4)) // 10

	// Spreading an existing slice into a variadic call:
	xs := []int{10, 20, 30}
	fmt.Println("SumAll(xs...):       ", funcsref.SumAll(xs...)) // 60

	fmt.Println(funcsref.LogPrefixed("INFO", "user", "logged", "in"))
	fmt.Println(funcsref.LogPrefixed("ERROR", "no", "such", "file"))

	// -----------------------------------------------------------------
	section("6. Functions as values")
	// Apply takes a function as a parameter.
	fmt.Println("Apply(10, 3, AddOp):", funcsref.Apply(10, 3, funcsref.AddOp))
	fmt.Println("Apply(10, 3, SubOp):", funcsref.Apply(10, 3, funcsref.SubOp))

	// Inline anonymous function passed at the call site:
	mul := func(a, b int) int { return a * b }
	fmt.Println("Apply(10, 3, mul):  ", funcsref.Apply(10, 3, mul))

	// Dispatch table: pick the operation by name.
	ops := funcsref.OpsByName()
	for _, name := range []string{"+", "-", "*", "/"} {
		fmt.Printf("OpsByName[%q](20, 6) = %d\n", name, ops[name](20, 6))
	}

	// -----------------------------------------------------------------
	section("7. Anonymous functions")
	fmt.Println("AnonymousImmediate:", funcsref.AnonymousImmediate()) // 5*5+7*7 = 74

	// FilterStrings with an inline predicate -- keeps non-empty strings.
	words := []string{"hi", "", "world", "", "go"}
	kept := funcsref.FilterStrings(words, func(s string) bool { return s != "" })
	fmt.Println("FilterStrings (non-empty):", kept)

	// Same function, different predicate -- keeps short strings.
	short := funcsref.FilterStrings(words, func(s string) bool { return len(s) > 0 && len(s) <= 2 })
	fmt.Println("FilterStrings (len<=2):   ", short)

	// -----------------------------------------------------------------
	section("8. Closures")
	// Each call to MakeCounter creates an INDEPENDENT counter.
	c1 := funcsref.MakeCounter()
	c2 := funcsref.MakeCounter()
	fmt.Println("c1():", c1(), c1(), c1()) // 1 2 3
	fmt.Println("c2():", c2(), c2())       // 1 2 -- independent
	fmt.Println("c1():", c1())             // 4 -- continues from where it was

	// Configuration captured at creation:
	double := funcsref.MakeMultiplier(2)
	triple := funcsref.MakeMultiplier(3)
	fmt.Println("double(5):", double(5)) // 10
	fmt.Println("triple(5):", triple(5)) // 15

	// Encapsulated state via paired closures:
	deposit, balance := funcsref.MakeAccount(100.0)
	deposit(50)
	deposit(25)
	deposit(10.50)
	fmt.Printf("Account balance: %.2f\n", balance()) // 185.50

	// -----------------------------------------------------------------
	section("9. Recursion")
	for i := 0; i <= 6; i++ {
		fmt.Printf("Factorial(%d) = %d\n", i, funcsref.Factorial(i))
	}
	fmt.Print("Fibonacci 0..10: ")
	for i := 0; i <= 10; i++ {
		fmt.Print(funcsref.Fibonacci(i), " ")
	}
	fmt.Println()

	// -----------------------------------------------------------------
	section("10. Realistic example: Calculator")
	calc := funcsref.NewCalculator(10)

	// Apply named operations one at a time.
	calc.Apply("double", func(x int) int { return x * 2 })
	calc.Apply("plus5", func(x int) int { return x + 5 })

	// ApplyMany takes any number of NamedOp values via variadic.
	calc.ApplyMany(
		funcsref.NamedOp{Name: "square", Op: func(x int) int { return x * x }},
		funcsref.NamedOp{Name: "minus100", Op: func(x int) int { return x - 100 }},
	)

	fmt.Println(calc.Summary())
	fmt.Println("History:")
	for _, h := range calc.History {
		fmt.Println("  ", h)
	}

	// PreloadedCalculator returns BOTH a calculator and a map of ops.
	calc2, presets := funcsref.PreloadedCalculator(5, 3)
	calc2.Apply("add(3)", presets["add"])
	calc2.Apply("multiply(3)", presets["multiply"])
	calc2.Apply("square", presets["square"])
	fmt.Println()
	fmt.Println("Preloaded:", calc2.Summary())
	for _, h := range calc2.History {
		fmt.Println("  ", h)
	}

	// Pointers
	// -----------------------------------------------------------------
	section("1. Pointer basics: & and *")
	fmt.Println(pointersref.ShowAddress())
	direct, viaPointer := pointersref.ReadThrough()
	fmt.Printf("ReadThrough:  direct=%d via pointer=%d (same value, two paths)\n", direct, viaPointer)
	fmt.Println("WriteThrough:", pointersref.WriteThrough(), "(x was 42, now 100 via *p = 100)")
	p, q := pointersref.Aliasing()
	fmt.Printf("Aliasing:     *p=%d *q=%d (both see the same change)\n", p, q)

	// -----------------------------------------------------------------
	section("2. Nil pointers")
	fmt.Println("Zero value of *int is nil?", pointersref.ZeroValue())

	// SafeDeref: nil case
	v, ok := pointersref.SafeDeref(nil)
	fmt.Printf("SafeDeref(nil):  val=%d ok=%v\n", v, ok)

	// SafeDeref: valid case
	x_0 := 123
	v, ok = pointersref.SafeDeref(&x_0)
	fmt.Printf("SafeDeref(&x):   val=%d ok=%v\n", v, ok)

	// -----------------------------------------------------------------
	section("3. Creating pointers: &, new, constructors")
	p1 := pointersref.WithAddress()
	fmt.Printf("WithAddress():    *p = %d\n", *p1)

	p2 := pointersref.WithNew()
	fmt.Printf("WithNew():        *p = %d\n", *p2)

	cfg := pointersref.NewConfig("example.com")
	fmt.Printf("NewConfig('example.com'): %+v\n", cfg)

	// -----------------------------------------------------------------
	section("4. Pointers with structs")
	pt := pointersref.Point{X: 0, Y: 0}

	// Value: returns a new modified Point; the caller's pt is untouched
	// unless they reassign.
	moved := pointersref.MoveValue(pt, 3, 4)
	fmt.Printf("MoveValue:    original=%+v returned=%+v\n", pt, moved)

	// Pointer: modifies the caller's Point directly.
	pointersref.MovePointer(&pt, 10, 20)
	fmt.Printf("MovePointer:  original=%+v (modified in place)\n", pt)

	// first, second := pointersref.ShareStruct()
	// fmt.Printf("ShareStruct:  via p=%+v via q=%+v (same underlying Point)\n", first, second)

	local := pointersref.ReturnLocalPointer()
	fmt.Printf("ReturnLocalPointer: %+v (safe in Go thanks to escape analysis)\n", *local)

	// -----------------------------------------------------------------
	section("5. Pointers with slices")
	nums_0 := []int{1, 2, 3, 4}
	pointersref.DoubleSliceContents(nums_0)
	fmt.Println("DoubleSliceContents:", nums, "(mutations visible -- slice header was copied, but data is shared)")

	// Demonstrate the append pitfall.
	nums2 := []int{1, 2, 3}
	pointersref.AppendWithoutPointer(nums2, 99)
	fmt.Println("AppendWithoutPointer:", nums2, "(unchanged! - the function's header was a copy)")

	// Correct: reassign the return value.
	nums2 = pointersref.AppendReturning(nums2, 99)
	fmt.Println("AppendReturning:    ", nums2, "(idiomatic Go: return new slice, caller reassigns)")

	// Alternative: pass *[]int.
	pointersref.AppendViaPointer(&nums2, 111)
	fmt.Println("AppendViaPointer:   ", nums2, "(less common but sometimes useful)")

	// Slice of VALUES: range gives copies, mutation in loop doesn't stick.
	items := []pointersref.Item{
		{Name: "a", Count: 1},
		{Name: "b", Count: 2},
	}
	pointersref.IncrementAllValues(items)
	fmt.Println("IncrementAllValues:       ", items, "(BUG: nothing was modified!)")

	pointersref.IncrementAllValuesFixed(items)
	fmt.Println("IncrementAllValuesFixed:  ", items, "(index-based mutation works)")

	// Slice of POINTERS: range gives pointer copies, but they point to the same data.
	itemPtrs := []*pointersref.Item{
		{Name: "a", Count: 1},
		{Name: "b", Count: 2},
	}
	pointersref.IncrementAllPointers(itemPtrs)
	fmt.Print("IncrementAllPointers:      ")
	for _, p := range itemPtrs {
		fmt.Printf("%+v ", *p)
	}
	fmt.Println("(works without index tricks)")

	// -----------------------------------------------------------------
	section("6. Pointers with maps")
	scores := map[string]int{"Asha": 100}
	pointersref.AddToMap(scores, "Bo", 90)
	fmt.Println("AddToMap:", scores, "(passing a map already lets the function mutate it)")

	// Value in a map: can't take address, so we do copy-modify-write.
	inventory := map[string]pointersref.Item{
		"apple": {Name: "apple", Count: 5},
	}
	pointersref.IncrementMapValue(inventory, "apple")
	fmt.Printf("IncrementMapValue:    %+v (workaround 1: read, modify, write back)\n", inventory)

	// Pointer in a map: can modify in place through the pointer.
	inventoryPtr := map[string]*pointersref.Item{
		"apple": {Name: "apple", Count: 5},
	}
	pointersref.IncrementMapPointer(inventoryPtr, "apple")
	fmt.Printf("IncrementMapPointer:  apple=%+v (workaround 2: store pointers directly)\n", *inventoryPtr["apple"])

	// -----------------------------------------------------------------
	section("7. Nil-pointer panic (recovered for display)")
	fmt.Println(pointersref.DerefNilRecovered())
	fmt.Println("(program didn't crash because we used recover(); don't rely on this in real code)")

	// -----------------------------------------------------------------
	section("8. Memory: stack vs heap")
	heapP := pointersref.MakeOnHeap()
	fmt.Printf("MakeOnHeap:   %+v (escape analysis placed the Point on the heap)\n", *heapP)

	stackP := pointersref.NoEscape()
	fmt.Printf("NoEscape:     %+v (returned by value; stays on the stack)\n", stackP)

	fmt.Println()
	fmt.Println("You can verify these decisions with:")
	fmt.Println("    go build -gcflags=\"-m\" ./...")
	fmt.Println("and look for 'moved to heap' / 'does not escape' lines.")

	// Show memory stats before, during, and after some work.
	heapBefore, totalBefore, gcBefore := pointersref.MemStats()
	fmt.Printf("\nBefore work:  heap=%dKB  cumulative=%dKB  gc_cycles=%d\n",
		heapBefore/1024, totalBefore/1024, gcBefore)

	// -----------------------------------------------------------------
	section("9. GC-friendly vs GC-hostile patterns")
	// Do some allocation work.
	const iters = 10_000
	_ = pointersref.AllocateInLoop(iters)
	heapMid, totalMid, gcMid := pointersref.MemStats()
	fmt.Printf("After AllocateInLoop(%d):  heap=%dKB  cumulative=%dKB  gc_cycles=%d\n",
		iters, heapMid/1024, totalMid/1024, gcMid)

	_ = pointersref.ReuseBuffer(iters)
	heapAfter, totalAfter, gcAfter := pointersref.MemStats()
	fmt.Printf("After ReuseBuffer(%d):     heap=%dKB  cumulative=%dKB  gc_cycles=%d\n",
		iters, heapAfter/1024, totalAfter/1024, gcAfter)

	fmt.Printf("\nNote cumulative-alloc difference:\n")
	fmt.Printf("  AllocateInLoop added ~%dKB of alloc pressure\n", (totalMid-totalBefore)/1024)
	fmt.Printf("  ReuseBuffer added    ~%dKB (essentially nothing)\n", (totalAfter-totalMid)/1024)

	// Force a GC to observe the counter increment.
	pointersref.ForceGC()
	_, _, gcForced := pointersref.MemStats()
	fmt.Printf("\nAfter ForceGC(): gc_cycles=%d (should have incremented)\n", gcForced)

	// Preallocation demo.
	sq_0 := pointersref.PreallocateSlice(10)
	fmt.Println("PreallocateSlice(10):", sq_0, "(make([]T, 0, n) avoids regrow-and-copy)")

	// -----------------------------------------------------------------
	section("10. Linked list -- the canonical pointer structure")
	list := pointersref.NewLinkedList()
	list.Append(1)
	list.Append(2)
	list.Append(3)
	list.Prepend(0) // front
	list.Append(4)  // back
	fmt.Println("After append/prepend:", list.ToSlice())
	fmt.Println("Length:", list.Length())

	removed := list.RemoveFirst(2)
	fmt.Printf("RemoveFirst(2) -> %v, list now: %v\n", removed, list.ToSlice())

	removed = list.RemoveFirst(999)
	fmt.Printf("RemoveFirst(999) -> %v (nothing to remove), list unchanged: %v\n", removed, list.ToSlice())

	// The removed node had its memory reclaimed automatically by the GC.
	// In C, we'd have to manually free it; in Go, unreachable = collectible.
	fmt.Println("\n(The removed node is now unreachable; the GC will reclaim its memory.")
	fmt.Println(" You did no malloc and no free. That's the whole point.)")

	// Interfaces
	// -----------------------------------------------------------------
	section("1. Basics: defining and satisfying an interface")
	// Pass two completely different concrete types to the same function.
	// SayHi accepts any Greeter, so both work.
	fmt.Println(interfacesref.SayHi(interfacesref.EnglishSpeaker{Name: "Asha"}))
	fmt.Println(interfacesref.SayHi(interfacesref.SpanishSpeaker{Name: "Bo"}))

	// -----------------------------------------------------------------
	section("2. Why interfaces matter -- one function, many types")
	// A heterogeneous slice of shapes -- something you couldn't express
	// without an interface.
	shapes := []interfacesref.Shape{
		interfacesref.Rectangle{Width: 3, Height: 4},
		interfacesref.Circle{Radius: 5},
		interfacesref.Triangle{Base: 6, Height: 8},
	}
	fmt.Printf("Total area = %.2f\n", interfacesref.TotalArea(shapes))

	// -----------------------------------------------------------------
	section("3. Empty interface -- `any`")
	// AnyDescription accepts anything. The %T verb shows the dynamic type
	// the interface value is carrying.
	fmt.Println(interfacesref.AnyDescription(42))
	fmt.Println(interfacesref.AnyDescription("hello"))
	fmt.Println(interfacesref.AnyDescription(true))
	fmt.Println(interfacesref.AnyDescription(interfacesref.EnglishSpeaker{Name: "Cleo"}))

	// -----------------------------------------------------------------
	section("4. Type assertions -- recover the concrete type")
	// Successful assertion: ok=true, value is what we asked for.
	s, ok := interfacesref.AssertString("hello world")
	fmt.Printf("AssertString(\"hello world\"): %q ok=%v\n", s, ok)

	// Failed assertion: ok=false, value is the zero value (no panic).
	s, ok = interfacesref.AssertString(42)
	fmt.Printf("AssertString(42):            %q ok=%v\n", s, ok)

	// Asserting for an INTERFACE rather than a concrete type:
	// "is this thing capable of greeting?"
	_, ok = interfacesref.AssertGreeter(interfacesref.EnglishSpeaker{Name: "Drew"})
	fmt.Printf("AssertGreeter(EnglishSpeaker): ok=%v\n", ok)

	_, ok = interfacesref.AssertGreeter(42)
	fmt.Printf("AssertGreeter(42):             ok=%v\n", ok)

	// -----------------------------------------------------------------
	section("5. Type switch -- handle multiple possibilities")
	for _, v := range []any{
		nil,
		7,
		"hello",
		true,
		interfacesref.SpanishSpeaker{Name: "Eli"}, // matches the Greeter case
		3.14, // no case matches -> default
	} {
		fmt.Printf("  %v -> %s\n", v, interfacesref.Describe(v))
	}

	// -----------------------------------------------------------------
	section("6. Embedding interfaces")
	store := &interfacesref.MemoryStore{}
	// store satisfies Reader, Writer, and ReadWriter all at once -- because
	// it has both Read and Write methods.
	result := interfacesref.UseReadWriter(store)
	fmt.Println("UseReadWriter result:", result)

	// -----------------------------------------------------------------
	section("7. Pointer receivers and interface satisfaction")
	c_0 := interfacesref.NewIntCounter()
	final := interfacesref.IncrementMany(c_0, 5)
	fmt.Printf("After 5 increments, counter = %d\n", final)

	// If we tried `interfacesref.IncrementMany(*c, 5)` instead, we'd get
	// a compile error: "IntCounter does not implement Counter (Increment
	// method has pointer receiver)". Useful to know -- the compiler tells
	// you exactly what's wrong.

	// -----------------------------------------------------------------
	section("8. Realistic example: pluggable log sinks")
	lines := []string{"server started", "user logged in", "shutting down"}

	// Sink #1: writes to an in-memory string builder.
	var sb strings.Builder
	console := &interfacesref.ConsoleSink{Buffer: &sb}
	interfacesref.LogAll(console, lines)
	fmt.Println("Console sink output:")
	fmt.Print(sb.String())

	// Sink #2: stores lines in a slice. Good for tests.
	memory := &interfacesref.MemorySink{}
	interfacesref.LogAll(memory, lines)
	fmt.Printf("Memory sink captured %d lines: %v\n", len(memory.Lines), memory.Lines)

	// Sink #3: a PrefixSink that wraps another sink. Decorator pattern --
	// because PrefixSink itself satisfies LogSink, it can stand in
	// anywhere a LogSink is expected, including being wrapped again.
	memory2 := &interfacesref.MemorySink{}
	prefixed := &interfacesref.PrefixSink{Prefix: "[INFO]", Inner: memory2}
	interfacesref.LogAll(prefixed, lines)
	fmt.Printf("Prefixed sink captured: %v\n", memory2.Lines)

	// Generics

	// -----------------------------------------------------------------
	section("1. First generic function -- Max")
	// Same function, three different types. Type inference figures out T.
	fmt.Println("Max(3, 7):           ", genericsref.Max(3, 7))
	fmt.Println("Max(3.14, 2.71):     ", genericsref.Max(3.14, 2.71))
	fmt.Println("Max(\"apple\", \"pear\"):", genericsref.Max("apple", "pear"))

	fmt.Println("Min(3, 7):           ", genericsref.Min(3, 7))

	// -----------------------------------------------------------------
	section("2. Custom constraint -- Number, Sum")
	fmt.Println("Sum([]int{1,2,3,4,5}):       ", genericsref.Sum([]int{1, 2, 3, 4, 5}))
	fmt.Println("Sum([]float64{1.5, 2.5, 3}): ", genericsref.Sum([]float64{1.5, 2.5, 3}))
	// Sum([]string{...}) would NOT compile -- string isn't in `Number`.

	// -----------------------------------------------------------------
	section("3. The `comparable` constraint -- Contains, IndexOf")
	fmt.Println("Contains([1 2 3], 2):       ", genericsref.Contains([]int{1, 2, 3}, 2))
	fmt.Println("Contains([\"a\" \"b\"], \"x\"):  ", genericsref.Contains([]string{"a", "b"}, "x"))
	fmt.Println("IndexOf([10 20 30], 20):    ", genericsref.IndexOf([]int{10, 20, 30}, 20))
	fmt.Println("IndexOf([10 20 30], 99):    ", genericsref.IndexOf([]int{10, 20, 30}, 99))

	// -----------------------------------------------------------------
	section("4. Type inference -- when you can omit [T]")
	for _, line := range genericsref.CallExamples() {
		fmt.Println(" ", line)
	}

	// -----------------------------------------------------------------
	section("5. Generic higher-order functions")
	nums_1 := []int{1, 2, 3, 4, 5}

	// Map: int -> int (square each)
	squared := genericsref.MapSlice(nums_1, func(n int) int { return n * n })
	fmt.Println("MapSlice (square):     ", squared)

	// Map: int -> string (different types -- T=int, U=string)
	asStrings := genericsref.MapSlice(nums, func(n int) string {
		return fmt.Sprintf("n%d", n)
	})
	fmt.Println("MapSlice (int->string):", asStrings)

	// Filter: keep even numbers
	evens := genericsref.FilterSlice(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("FilterSlice (evens):   ", evens)

	// Reduce: sum (T=int, U=int)
	total := genericsref.Reduce(nums, 0, func(acc, v int) int { return acc + v })
	fmt.Println("Reduce (sum):          ", total)

	// Reduce: build a string (T=int, U=string -- different types!)
	joined := genericsref.Reduce(nums, "", func(acc string, v int) string {
		if acc == "" {
			return fmt.Sprintf("%d", v)
		}
		return acc + "," + fmt.Sprintf("%d", v)
	})
	fmt.Println("Reduce (join):         ", joined)

	// -----------------------------------------------------------------
	section("6. Generic type -- Stack[T]")
	// Two stacks of completely different types from the same definition.
	intStack := genericsref.NewStack[int]()
	intStack.Push(1)
	intStack.Push(2)
	intStack.Push(3)
	fmt.Printf("intStack: len=%d\n", intStack.Len())
	for intStack.Len() > 0 {
		v, _ := intStack.Pop()
		fmt.Printf("  popped %d\n", v)
	}

	stringStack := genericsref.NewStack[string]()
	stringStack.Push("first")
	stringStack.Push("second")
	v_1, _ := stringStack.Peek()
	fmt.Printf("stringStack peek: %q (len still %d)\n", v_1, stringStack.Len())

	// Note: intStack and stringStack are DIFFERENT TYPES.
	// You can't assign one to the other; the compiler treats Stack[int]
	// and Stack[string] as completely separate.

	// -----------------------------------------------------------------
	section("7. Generic type with two parameters -- KVPair, KeyValueStore")
	pair := genericsref.KVPair[string, int]{Key: "score", Value: 42}
	fmt.Println("KVPair:", pair) // uses the String() method

	store_0 := genericsref.NewKVStore[string, int]()
	store_0.Set("apple", 50)
	store_0.Set("banana", 20)
	store_0.Set("cherry", 80)
	apple, ok := store_0.Get("apple")
	fmt.Printf("store[\"apple\"] = %d (ok=%v)\n", apple, ok)
	missing, ok := store_0.Get("grape")
	fmt.Printf("store[\"grape\"] = %d (ok=%v -- not in store)\n", missing, ok)

	// -----------------------------------------------------------------
	section("8. Generic interface -- Container[T] + SumContainer")
	list_0 := genericsref.NewSimpleList[int]()
	list_0.Add(10)
	list_0.Add(20)
	list_0.Add(30)
	// SimpleList[int] satisfies Container[int].
	// SumContainer accepts any Container[T] where T is a Number.
	fmt.Println("SumContainer(list):", genericsref.SumContainer[int](list_0))

	// Error Hadnling in Go
	// -----------------------------------------------------------------
	section("1. Basic error handling -- (result, error) pattern")
	// Successful case: use the result, ignore (or check nil) the error.
	q_2, err := errorsref.Divide(10, 2)
	fmt.Printf("Divide(10, 2): result=%g err=%v\n", q_2, err)

	// Failure case: result is meaningless, error tells you what happened.
	q_2, err = errorsref.Divide(10, 0)
	fmt.Printf("Divide(10, 0): result=%g err=%v\n", q_2, err)

	// fmt.Errorf with formatted message.
	scores_0 := map[string]int{"asha": 100}
	_, err = errorsref.Lookup(scores_0, "bo")
	fmt.Printf("Lookup missing key: err=%v\n", err)

	// -----------------------------------------------------------------
	section("2. Sentinel errors -- exported error values")
	// Successful FindUser.
	name, err := errorsref.FindUser(1)
	fmt.Printf("FindUser(1): name=%q err=%v\n", name, err)

	// Failed FindUser -- returns the sentinel directly (not wrapped).
	_, err = errorsref.FindUser(99)
	fmt.Printf("FindUser(99): err=%v\n", err)

	// Plain == works HERE because the error wasn't wrapped.
	fmt.Printf("err == ErrUserNotFound? %v\n", err == errorsref.ErrUserNotFound)
	// In real code, prefer errors.Is even for unwrapped sentinels --
	// it survives later wrapping if you ever change the code.
	fmt.Printf("errors.Is(err, ErrUserNotFound)? %v\n",
		errors.Is(err, errorsref.ErrUserNotFound))

	// -----------------------------------------------------------------
	section("3. Custom error types -- structured errors")
	err = errorsref.ValidateEmail("")
	fmt.Printf("ValidateEmail(\"\"):       err=%v\n", err)

	err = errorsref.ValidateEmail("nopeAtNothing")
	fmt.Printf("ValidateEmail(\"nope...\"): err=%v\n", err)

	err = errorsref.ValidateEmail("ok@example.com")
	fmt.Printf("ValidateEmail(ok):       err=%v\n", err)

	// -----------------------------------------------------------------
	section("4. Wrapping errors with %w")
	// StartApp wraps LoadUserConfig wraps FindUser. Each layer added
	// context, but the original sentinel is still inside.
	err = errorsref.StartApp()
	fmt.Println("Full chain message:", err)

	// Look at the chain by repeatedly calling errors.Unwrap.
	fmt.Println("\nWalking the chain with errors.Unwrap:")
	for e := err; e != nil; e = errors.Unwrap(e) {
		fmt.Printf("  - %v\n", e)
	}

	// -----------------------------------------------------------------
	section("5. Unwrapping with errors.Is and errors.As")
	// The error from StartApp is wrapped multiple times, but errors.Is
	// finds ErrUserNotFound inside.
	fmt.Println("Is wrapped err equal to ErrUserNotFound?",
		errors.Is(err, errorsref.ErrUserNotFound))
	fmt.Println("Is wrapped err equal to ErrInvalidEmail?",
		errors.Is(err, errorsref.ErrInvalidEmail))

	// errors.As: extract a *ValidationError from a chain.
	emailErr := errorsref.ValidateEmail("badformat")
	field, msg, ok := errorsref.ExtractValidationError(emailErr)
	fmt.Printf("Extracted from ValidationError: field=%q msg=%q ok=%v\n",
		field, msg, ok)

	// errors.As against an error that ISN'T a *ValidationError.
	field, msg, ok = errorsref.ExtractValidationError(errorsref.ErrUserNotFound)
	fmt.Printf("Extracted from sentinel:        field=%q msg=%q ok=%v\n",
		field, msg, ok)

	// -----------------------------------------------------------------
	section("6. defer -- cleanup and return interaction")
	fmt.Println("DeferOrder:       returns =", errorsref.DeferOrder())
	// Note: DeferOrder returns "body". The defers DID run, but they
	// modified a local builder that we already extracted before they
	// fired. The returned string was set BEFORE any defer ran.

	fmt.Println("DeferOrderNamed:  returns =", errorsref.DeferOrderNamed())
	// Now we use a NAMED return. The defers append to it. The caller
	// sees the final, post-defer value.

	snap, latest := errorsref.DeferArgsCapture()
	fmt.Printf("DeferArgsCapture: snapshot(captured)=%d latest(closure)=%d\n",
		snap, latest)
	// snapshot is 1: the value of x AT defer registration time.
	// latest is 99: the closure read x by reference, after we changed it.

	// -----------------------------------------------------------------
	section("7. panic and recover")
	// Successful path: no panic.
	r_0, err := errorsref.SafeDivide(10, 2)
	fmt.Printf("SafeDivide(10, 2): result=%d err=%v\n", r_0, err)

	// Failing path: MustDivide panics, SafeDivide's deferred recover
	// catches it and returns a normal error. The program does NOT crash.
	r_0, err = errorsref.SafeDivide(10, 0)
	fmt.Printf("SafeDivide(10, 0): result=%d err=%v\n", r_0, err)
	fmt.Println("(notice: the program is still running -- recover prevented the crash)")

	// -----------------------------------------------------------------
	section("8. Putting it together -- a registration flow")
	store_1 := errorsref.NewUserStore()

	// Validation failure path.
	fmt.Println(store_1.HandleRegister(1, ""))
	// Validation failure (no @).
	fmt.Println(store_1.HandleRegister(1, "noatsign"))
	// Success.
	fmt.Println(store_1.HandleRegister(1, "asha@example.com"))
	// Conflict (id already used).
	fmt.Println(store_1.HandleRegister(1, "another@example.com"))

	// Goroutines
	goroutines.ExecuteExample()

	fmt.Println("ChannelExample: ", goroutines.ChannelExample())

	goroutines.BufferedChannelExample()

	goroutines.RangeOverChannelExample()

	fmt.Println("SelectChannelExample: ", goroutines.SelectChannelExample())

	fmt.Println("TimeoutSelectChannelExample: ", goroutines.TimeoutSelectChannelExample())

	fmt.Println("DefaultSelectChannelExample: ", goroutines.DefaultSelectChannelExample())

	goroutines.WaitgroupExample()

	fmt.Println("MutexExample Result: ", goroutines.MutexExample())

	goroutines.WorkerPoolExample()

	// ContextRefExample
	// -----------------------------------------------------------------
	section("1. Why: raw channel vs context for cancellation")
	fmt.Println("Raw channel:", contextref.CancelWithRawChannel())
	fmt.Println("Context:    ", contextref.CancelWithContext())
	fmt.Println()
	fmt.Println("Both achieve the same thing. Context is the standard,")
	fmt.Println("composable way; raw channels work but don't scale.")

	// -----------------------------------------------------------------
	section("2. WithCancel: manual cancellation")
	count := contextref.RunWithManualCancel()
	fmt.Printf("Background job ran %d iterations before cancel\n", count)

	// -----------------------------------------------------------------
	section("3. WithTimeout: automatic cancellation after a duration")
	// Slow operation: timeout fires before completion
	err_1 := contextref.RunFetchTooSlow()
	fmt.Printf("Slow fetch: err = %v\n", err_1)

	// Fast operation: completes before timeout
	res_3, err := contextref.RunFetchInTime()
	fmt.Printf("Fast fetch: result = %q, err = %v\n", res_3, err)

	// -----------------------------------------------------------------
	section("4. Classifying context errors")
	// Build a context, cancel it, and observe the classification.
	ctx1, cancel1 := context.WithCancel(context.Background())
	cancel1()
	fmt.Println("After manual cancel:    ", contextref.ClassifyContextError(ctx1.Err()))

	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel2()
	time.Sleep(5 * time.Millisecond)
	fmt.Println("After deadline expired: ", contextref.ClassifyContextError(ctx2.Err()))

	// -----------------------------------------------------------------
	section("5. Tree propagation: parent cancel propagates to child")
	parentErr, childErr := contextref.ParentChildPropagation()
	fmt.Println("Parent error:", parentErr)
	fmt.Println("Child error: ", childErr)
	fmt.Println("Both got canceled when we canceled only the parent.")

	fmt.Println()
	fmt.Println("Fan-out: 4 workers all share one context with a 30ms deadline")
	results := contextref.FanOutWithCancel(4)
	for i, c := range results {
		fmt.Printf("  worker %d: %d iterations\n", i, c)
	}

	// -----------------------------------------------------------------
	section("6. WithValue: request-scoped data")
	fmt.Println(contextref.SimulateRequestHandler())
	fmt.Println()
	fmt.Println("Use sparingly! WithValue is for cross-cutting concerns")
	fmt.Println("(request IDs, trace spans, auth info) -- NOT real arguments.")

	// -----------------------------------------------------------------
	section("7. Realistic mini-handler")
	// Successful case: parent context never cancels, query completes in time
	result_3, err := contextref.HandleRequest(context.Background(), "req-001")
	fmt.Printf("Normal request: result=%q err=%v\n", result_3, err)

	// Parent context with very short deadline -- handler will time out
	tightCtx, tightCancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer tightCancel()
	result, err = contextref.HandleRequest(tightCtx, "req-002")
	fmt.Printf("Tight deadline: result=%q err=%v\n", result, err)

}
