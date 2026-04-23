package main

import (
	"go_revision/arrays"
	customstrings "go_revision/customStrings"
	"go_revision/loopsref"
	"go_revision/make_practice"
	"go_revision/maps"
	"go_revision/slices"
	"go_revision/structs"
	"go_revision/utils"

	"encoding/json"
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

}
