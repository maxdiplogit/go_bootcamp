package main

import (
	"go_revision/utils"
	"go_revision/structs"
	"go_revision/arrays"
	"go_revision/slices"
	"go_revision/make_practice"
	"go_revision/maps"
	"go_revision/customStrings"

	"fmt"
	"time"
	"strings"
	"encoding/json"
)

// section prints a visual separator so each demo block is easy to find in
// the terminal output.
func section(title string) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("  " + title)
	fmt.Println(strings.Repeat("=", 60))
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

	var vehicle_0 *structs.Vehicle;
	vehicle_0 = &structs.Vehicle{
		Company: "BMW",
		FuelType: "Petrol",
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
		Tyres: 6,
		Pistons: 8,
	}
	truck_0.TruckJsonExample()

	fmt.Println()

	userIdentifier := &structs.Identifier{
		ID: 1,
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
    	fmt.Println(u.UserIdentifier)         // from Identifier
    	fmt.Println(u.CreatedAt)  // from Timestamps
    	fmt.Println(u.Version)    // from AuditTrail
    	fmt.Println(u.Email)      // User's own field

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
}
