// That single distinction is the source of ~90% of all string bugs in Go.
package customstrings
 
import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)
 
// =============================================================================
// SECTION 1: BASIC STRING FACTS
// =============================================================================
 
// StringHeader explains (via a returned string) what a Go string actually is
// at the implementation level. Useful as a mental anchor.
//
// Internally, a Go string is a tiny 16-byte header on 64-bit systems:
//   - a pointer to the underlying byte data
//   - a length (number of bytes)
//
// This is why passing strings around is CHEAP regardless of their length:
// you copy 16 bytes (the header), not the underlying data. The underlying
// byte array is shared and is immutable, which is what makes sharing safe.
func StringHeader() string {
	return "A Go string = (pointer to bytes, length). Immutable. Cheap to pass."
}
 
// LenIsBytes demonstrates that len(s) counts BYTES, not characters.
// This is the #1 gotcha in Go strings.
//
// For pure ASCII input, bytes == characters, so beginners don't notice.
// For anything with accents, CJK characters, or emoji, they differ.
func LenIsBytes(s string) (byteCount int, charCount int) {
	byteCount = len(s) // bytes
	// utf8.RuneCountInString decodes UTF-8 and counts actual code points.
	charCount = utf8.RuneCountInString(s)
	return
}
 
// IndexingReturnsBytes shows that s[i] yields a byte (uint8), not a character.
// It returns the byte at position i along with its decimal value.
//
// If i lands in the middle of a multi-byte UTF-8 sequence, you get a
// "continuation byte" that is NOT a valid character on its own.
func IndexingReturnsBytes(s string, i int) (b byte, asInt int, ok bool) {
	if i < 0 || i >= len(s) {
		return 0, 0, false
	}
	b = s[i]
	asInt = int(b)
	ok = true
	return
}
 
// StringsAreImmutable documents (but cannot demonstrate directly at runtime)
// the fact that strings cannot be mutated in place.
//
// The line `s[0] = 'H'` is a COMPILE error, not a runtime error.
// To "modify" a string you must build a new one, typically by going
// through []byte or []rune.
func StringsAreImmutable() string {
	return "Strings are immutable. To modify, convert to []byte or []rune, change, then convert back."
}
 
// =============================================================================
// SECTION 2: UTF-8, BYTES, AND RUNES
// =============================================================================
 
// ByteVsRune returns, for every character of s:
//   - the byte index where the character starts
//   - the rune (code point) at that position
//   - how many bytes that rune occupies in UTF-8
//
// This function is the clearest demonstration of why `range` over a string
// behaves differently from indexing. Notice how the returned indices are
// NOT sequential when the input has non-ASCII characters -- each index
// reflects the byte offset where that rune starts.
func ByteVsRune(s string) []struct {
	ByteIndex int
	Rune      rune
	ByteLen   int
} {
	// Pre-allocate with a sensible capacity hint: at most len(s) runes.
	out := make([]struct {
		ByteIndex int
		Rune      rune
		ByteLen   int
	}, 0, len(s))
 
	// `for i, r := range s` iterates RUNES:
	//   i -> byte index where this rune starts (NOT a sequential counter)
	//   r -> the decoded rune (int32 code point)
	for i, r := range s {
		out = append(out, struct {
			ByteIndex int
			Rune      rune
			ByteLen   int
		}{
			ByteIndex: i,
			Rune:      r,
			ByteLen:   utf8.RuneLen(r),
		})
	}
	return out
}
 
// RawBytes returns the raw byte view of a string without any UTF-8 decoding.
// Good for understanding WHY len(s) is what it is.
//
// Note: converting string -> []byte COPIES the bytes. This is required
// because strings are immutable and []byte is mutable; Go enforces the
// boundary by making both conversions O(n) copies.
func RawBytes(s string) []byte {
	return []byte(s)
}
 
// RuneSlice returns the rune (code-point) view of a string. This is what
// you want when you need to count, index, or reorder CHARACTERS, not bytes.
//
// The conversion string -> []rune:
//   - decodes UTF-8
//   - allocates a new []rune slice
//   - has cost O(n) in both time and memory
//
// So do NOT call []rune(s) in a hot loop just to index one character.
// Use utf8.DecodeRuneInString for that.
func RuneSlice(s string) []rune {
	return []rune(s)
}
 
// =============================================================================
// SECTION 3: CONVERSIONS BETWEEN string, []byte, []rune
// =============================================================================
 
// BytesRoundTrip demonstrates string <-> []byte conversion.
//
// Use []byte when:
//   - doing I/O (io.Reader/io.Writer work in bytes)
//   - handling binary data (hashes, encryption, images)
//   - mutating ASCII text
//   - interfacing with APIs that require []byte (json.Unmarshal, crypto, etc.)
//
// WARNING: []byte works on BYTES, so if you mutate bytes in the middle of
// a multi-byte UTF-8 character, you corrupt the string. For non-ASCII edits,
// go through []rune instead.
func BytesRoundTrip(s string) string {
	b := []byte(s) // copy #1: string -> []byte
	// Safe example: uppercase ASCII letters in place (only touches 1-byte chars).
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32 // ASCII: 'a' = 97, 'A' = 65, diff = 32
		}
	}
	return string(b) // copy #2: []byte -> string
}
 
// RunesRoundTrip demonstrates string <-> []rune conversion.
//
// Use []rune when:
//   - you need to work with CHARACTERS, not bytes
//   - the input may contain non-ASCII text
//   - you're reversing, reordering, or counting characters
func RunesRoundTrip(s string) string {
	r := []rune(s) // decodes UTF-8 into code points
	// Safe example: replace every character with its uppercase form.
	// unicode.ToUpper works correctly on any rune, not just ASCII.
	for i, c := range r {
		r[i] = unicode.ToUpper(c)
	}
	return string(r) // encodes runes back into UTF-8 bytes
}
 
// ReverseASCII reverses a string assuming it is pure ASCII.
// This is FAST but BROKEN for any input with multi-byte characters --
// it will split those characters into invalid byte sequences.
//
// Use this only when you can GUARANTEE the input is ASCII
// (e.g., a hex digest, a base64 string, a fixed protocol header).
func ReverseASCII(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}
 
// ReverseUnicode correctly reverses a string containing any Unicode text.
// It works character-by-character (rune-by-rune) rather than byte-by-byte,
// so multi-byte sequences are preserved as single units.
//
// Cost: O(n) for the []rune conversion and O(n) for the reversal.
// Memory: one []rune allocation of n runes, plus the final string.
func ReverseUnicode(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
 
// RuneToString shows the safe idiom for turning a single rune (int32)
// into a string. Always wrap integers in rune() for clarity -- plain
// string(65) still works but modern Go vets will flag it as confusing.
func RuneToString(r rune) string {
	return string(r)
}
 
// =============================================================================
// SECTION 4: COUNTING THINGS CORRECTLY
// =============================================================================
 
// CountBytes returns the number of bytes in s. This is just len(s) wrapped
// for symmetry with CountRunes. O(1).
func CountBytes(s string) int {
	return len(s)
}
 
// CountRunes returns the number of Unicode code points in s. O(n) because
// it has to decode UTF-8. Prefer utf8.RuneCountInString over
// len([]rune(s)) -- the former does no allocation, the latter allocates a
// throwaway []rune.
func CountRunes(s string) int {
	return utf8.RuneCountInString(s)
}
 
// CountWords returns a rough whitespace-separated word count, using
// Unicode-aware whitespace detection. "Rough" because it does not try to
// handle locale-specific word segmentation -- that's a much harder problem.
func CountWords(s string) int {
	count := 0
	inWord := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			inWord = false
		} else if !inWord {
			inWord = true
			count++
		}
	}
	return count
}
 
// CountLines returns the number of lines in s, treating "\n" as the
// separator. A string with no newline still counts as 1 line.
// Note: \n is a single byte and cannot appear inside any multi-byte UTF-8
// sequence, so byte-level counting here is safe for any UTF-8 text.
func CountLines(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}
 
// =============================================================================
// SECTION 5: EFFICIENT STRING BUILDING
// =============================================================================
 
// ConcatWithPlus is the OBVIOUS but SLOW way to build a string in a loop.
// It's quadratic (O(n^2)) because every `+=` allocates a new string and
// copies both the old contents and the new piece.
//
// Use this only for a small, fixed number of concatenations. Never in a
// loop of unknown size.
func ConcatWithPlus(parts []string) string {
	s := ""
	for _, p := range parts {
		s += p // allocates + copies every iteration
	}
	return s
}
 
// ConcatWithBuilder is the idiomatic, FAST way to build a string from many
// pieces. strings.Builder manages a growing internal []byte and only
// materializes the final string once (with a zero-copy conversion).
//
// Calling Grow(n) up-front when you can estimate the final size avoids
// intermediate reallocations entirely.
func ConcatWithBuilder(parts []string) string {
	// Estimate total size to pre-grow the buffer.
	total := 0
	for _, p := range parts {
		total += len(p)
	}
 
	var b strings.Builder
	b.Grow(total) // single allocation up-front
 
	for _, p := range parts {
		b.WriteString(p) // appends to internal buffer, no intermediate string
	}
	return b.String() // single final conversion
}
 
// JoinWithSeparator uses strings.Join, which is internally equivalent to a
// pre-sized strings.Builder. Whenever you're joining a []string with a
// fixed separator, this is the cleanest choice.
func JoinWithSeparator(parts []string, sep string) string {
	return strings.Join(parts, sep)
}
 
// =============================================================================
// SECTION 6: COMMON STRING OPERATIONS (strings package)
// =============================================================================
 
// Transformations is a showcase of the most common operations from the
// `strings` package, collected in one place so you can grep for them later.
//
// Note that Index/IndexByte/IndexRune all return a BYTE index, not a
// character index. If you need a character offset, you have to iterate.
func Transformations(s string) map[string]any {
	return map[string]any{
		"Contains_hello":  strings.Contains(s, "hello"),
		"HasPrefix_Hi":    strings.HasPrefix(s, "Hi"),
		"HasSuffix_dot":   strings.HasSuffix(s, "."),
		"Index_of_space":  strings.Index(s, " "), // BYTE index; -1 if not found
		"Count_of_l":      strings.Count(s, "l"),
		"Replaced":        strings.ReplaceAll(s, "l", "L"),
		"Split_on_space":  strings.Split(s, " "),
		"Fields":          strings.Fields(s), // splits on ANY whitespace, collapses runs
		"ToLower":         strings.ToLower(s),
		"ToUpper":         strings.ToUpper(s),
		"TrimSpace":       strings.TrimSpace(s),
		"TrimPrefix":      strings.TrimPrefix(s, "Hello "),
		"Repeat_3":        strings.Repeat(s, 3),
	}
}
 
// =============================================================================
// SECTION 7: RUNE CLASSIFICATION (unicode package)
// =============================================================================
 
// ClassifyRune reports the Unicode category of a single rune. The key
// point: these functions work correctly on ANY rune, not just ASCII, so
// they're safe for internationalized text.
//
// Example: unicode.IsDigit('٥') returns true because that's the Arabic-Indic
// digit five -- a real digit, just not in the ASCII range.
func ClassifyRune(r rune) map[string]bool {
	return map[string]bool{
		"IsLetter":  unicode.IsLetter(r),
		"IsDigit":   unicode.IsDigit(r),
		"IsSpace":   unicode.IsSpace(r),
		"IsUpper":   unicode.IsUpper(r),
		"IsLower":   unicode.IsLower(r),
		"IsPunct":   unicode.IsPunct(r),
		"IsSymbol":  unicode.IsSymbol(r),
		"IsControl": unicode.IsControl(r),
	}
}
 
// =============================================================================
// SECTION 8: UTF-8 VALIDATION AND LOW-LEVEL DECODING
// =============================================================================
 
// IsValidUTF8 reports whether every byte in s forms part of a valid UTF-8
// sequence. Go does NOT enforce UTF-8 validity on string values; a string
// can legally contain arbitrary bytes, including invalid encodings.
// So if you're ingesting data from the outside world and you need to be
// sure it's valid UTF-8, check explicitly.
func IsValidUTF8(s string) bool {
	return utf8.ValidString(s)
}
 
// FirstRune returns the first rune in s and how many bytes it consumed.
// When the input is empty or malformed, DecodeRuneInString returns
// utf8.RuneError (U+FFFD) and size 0 or 1 -- check these if you care.
//
// This is the lowest-overhead way to peek at one character without
// allocating a []rune for the whole string.
func FirstRune(s string) (r rune, size int) {
	return utf8.DecodeRuneInString(s)
}
 
// =============================================================================
// SECTION 9: STRING LITERALS -- INTERPRETED VS RAW
// =============================================================================
 
// LiteralExamples documents (via comments) the two forms of string literal
// in Go, because you'll use both constantly.
//
//	Interpreted (double-quoted): escapes are processed.
//	    "line1\nline2\t\"quoted\""
//	Raw (backtick-quoted): no escapes, literal newlines allowed.
//	    `line1
//	    line2 "quoted"`
//
// Use raw strings for:
//   - regular expressions: `\d+\s*\w+`  (no double-escaping backslashes)
//   - multi-line SQL, HTML, or JSON fixtures
//   - Windows paths: `C:\Users\Asha`
//
// The only thing a raw string CAN'T contain is a backtick.
func LiteralExamples() (interpreted, raw string) {
	interpreted = "line1\nline2\t\"quoted\""
	raw = `line1
line2	"quoted"`
	return
}
 
// =============================================================================
// SECTION 10: A COMPLETE ANALYZER (ties it all together)
// =============================================================================
 
// TextStats is the result of analyzing a piece of text at multiple levels.
// This struct uses the same field-name discipline you learned earlier:
// exported fields, PascalCase.
type TextStats struct {
	Bytes    int  // len(s)
	Runes    int  // number of code points
	Words    int  // whitespace-separated chunks
	Lines    int  // newline-separated chunks
	ASCII    bool // true if every byte is < 128
	ValidUTF bool // true if the bytes form valid UTF-8
}
 
// Analyze returns a TextStats for the given string. It's the example from
// the tutorial, packaged as a reusable function.
//
// This is the function to look at if you ever forget "what's the right
// way to inspect a string in Go?" -- it uses the correct tool at each
// level: len() for bytes, utf8.RuneCountInString for characters,
// range+IsSpace for words, strings.Count for newlines.
func Analyze(s string) TextStats {
	stats := TextStats{
		Bytes:    len(s),
		Runes:    utf8.RuneCountInString(s),
		Lines:    CountLines(s),
		ValidUTF: utf8.ValidString(s),
		ASCII:    true, // assume true, invalidate on first non-ASCII byte
	}
 
	// Word counting + ASCII check in one pass over runes.
	inWord := false
	for _, r := range s {
		if r > 127 {
			stats.ASCII = false
		}
		if unicode.IsSpace(r) {
			inWord = false
		} else if !inWord {
			inWord = true
			stats.Words++
		}
	}
	return stats
}
 
// PrettyPrint returns a human-readable multi-line summary of a TextStats.
// Handy for quick debugging / REPL-style exploration from main.go.
func (t TextStats) PrettyPrint() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Bytes:    %d\n", t.Bytes)
	fmt.Fprintf(&b, "Runes:    %d\n", t.Runes)
	fmt.Fprintf(&b, "Words:    %d\n", t.Words)
	fmt.Fprintf(&b, "Lines:    %d\n", t.Lines)
	fmt.Fprintf(&b, "ASCII:    %v\n", t.ASCII)
	fmt.Fprintf(&b, "ValidUTF: %v\n", t.ValidUTF)
	return b.String()
}

