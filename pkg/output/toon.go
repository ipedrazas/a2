package output

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/language"
	"github.com/ipedrazas/a2/pkg/maturity"
	"github.com/ipedrazas/a2/pkg/runner"
)

// toonEncoder handles TOON format encoding.
type toonEncoder struct {
	indent     int
	indentSize int
	builder    strings.Builder
}

// newToonEncoder creates a new TOON encoder with default settings.
func newToonEncoder() *toonEncoder {
	return &toonEncoder{
		indent:     0,
		indentSize: 2,
	}
}

// TOON outputs the results in TOON format (Token-Oriented Object Notation).
// TOON is optimized for consumption by coding agents.
// Returns true if all checks passed, false otherwise, along with any output error.
func TOON(result runner.SuiteResult, detected language.DetectionResult, verbosity VerbosityLevel) (bool, error) {
	enc := newToonEncoder()

	// Convert languages to strings
	langs := make([]string, len(detected.Languages))
	for i, l := range detected.Languages {
		langs[i] = string(l)
	}

	// Calculate maturity estimation
	est := maturity.Estimate(result)

	// Build output
	// languages array
	enc.writeStringArray("languages", langs)

	// results array (tabular format for uniform objects)
	enc.writeResultsArray(result.Results, verbosity)

	// summary object
	enc.writeKey("summary")
	enc.indent++
	enc.writeKeyValue("total", enc.formatNumber(float64(result.ScoredChecks())))
	enc.writeKeyValue("passed", enc.formatNumber(float64(result.Passed)))
	enc.writeKeyValue("warnings", enc.formatNumber(float64(result.Warnings)))
	enc.writeKeyValue("failed", enc.formatNumber(float64(result.Failed)))
	enc.writeKeyValue("info", enc.formatNumber(float64(result.Info)))
	enc.writeKeyValue("score", enc.formatNumber(calculateScore(result)))
	enc.writeKeyValue("total_duration_ms", enc.formatNumber(float64(result.TotalDuration.Milliseconds())))
	enc.indent--

	// maturity object
	enc.writeKey("maturity")
	enc.indent++
	enc.writeKeyValue("level", enc.encodeString(est.Level.String()))
	enc.writeKeyValue("description", enc.encodeString(est.Level.Description()))
	enc.writeStringArray("suggestions", est.Suggestions)
	enc.indent--

	// boolean fields
	enc.writeKeyValue("aborted", enc.formatBool(result.Aborted))
	enc.writeKeyValue("success", enc.formatBool(result.Success()))

	// Output without trailing newline (per TOON spec)
	output := strings.TrimSuffix(enc.builder.String(), "\n")
	fmt.Print(output)

	return result.Success(), nil
}

// writeIndent writes the current indentation.
func (e *toonEncoder) writeIndent() {
	e.builder.WriteString(strings.Repeat(" ", e.indent*e.indentSize))
}

// writeKey writes a key followed by colon (for nested objects).
func (e *toonEncoder) writeKey(key string) {
	e.writeIndent()
	e.builder.WriteString(e.encodeKey(key))
	e.builder.WriteString(":\n")
}

// writeKeyValue writes a key-value pair.
func (e *toonEncoder) writeKeyValue(key, value string) {
	e.writeIndent()
	e.builder.WriteString(e.encodeKey(key))
	e.builder.WriteString(": ")
	e.builder.WriteString(value)
	e.builder.WriteByte('\n')
}

// writeStringArray writes an array of strings in inline format.
func (e *toonEncoder) writeStringArray(key string, values []string) {
	e.writeIndent()
	e.builder.WriteString(e.encodeKey(key))
	fmt.Fprintf(&e.builder, "[%d]:", len(values))
	if len(values) > 0 {
		e.builder.WriteByte(' ')
		encoded := make([]string, len(values))
		for i, v := range values {
			encoded[i] = e.encodeStringForArray(v, ',')
		}
		e.builder.WriteString(strings.Join(encoded, ","))
	}
	e.builder.WriteByte('\n')
}

// writeResultsArray writes the results array in tabular format.
func (e *toonEncoder) writeResultsArray(results []checker.Result, verbosity VerbosityLevel) {
	e.writeIndent()
	// Use tabular format: results[N]{fields}:
	// Include raw_output field when verbosity > 0
	if verbosity > VerbosityNormal {
		fmt.Fprintf(&e.builder, "results[%d]{name,id,passed,status,message,language,duration_ms,raw_output}:\n", len(results))
	} else {
		fmt.Fprintf(&e.builder, "results[%d]{name,id,passed,status,message,language,duration_ms}:\n", len(results))
	}
	e.indent++
	for _, r := range results {
		e.writeIndent()
		// Each row: name,id,passed,status,message,language,duration_ms[,raw_output]
		row := []string{
			e.encodeStringForArray(r.Name, ','),
			e.encodeStringForArray(r.ID, ','),
			e.formatBool(r.Passed),
			e.encodeStringForArray(statusToString(r.Status), ','),
			e.encodeStringForArray(r.Message, ','),
			e.encodeStringForArray(string(r.Language), ','),
			e.formatNumber(float64(r.Duration.Milliseconds())),
		}

		// Include raw output based on verbosity level
		if verbosity > VerbosityNormal {
			rawOutput := ""
			if r.RawOutput != "" {
				shouldInclude := verbosity == VerbosityAll ||
					(verbosity == VerbosityFailures && (r.Status == checker.Fail || r.Status == checker.Warn))
				if shouldInclude {
					rawOutput = r.RawOutput
				}
			}
			row = append(row, e.encodeStringForArray(rawOutput, ','))
		}

		e.builder.WriteString(strings.Join(row, ","))
		e.builder.WriteByte('\n')
	}
	e.indent--
}

// encodeKey encodes a key, quoting if necessary.
// Keys may be unquoted only if matching: ^[A-Za-z_][A-Za-z0-9_.]*$
func (e *toonEncoder) encodeKey(key string) string {
	if e.isValidUnquotedKey(key) {
		return key
	}
	return e.quoteString(key)
}

// isValidUnquotedKey checks if a key can be unquoted.
var validKeyRegex = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_.]*$`)

func (e *toonEncoder) isValidUnquotedKey(key string) bool {
	return validKeyRegex.MatchString(key)
}

// encodeString encodes a string value, quoting if necessary.
func (e *toonEncoder) encodeString(s string) string {
	if e.needsQuoting(s, ',') {
		return e.quoteString(s)
	}
	return s
}

// encodeStringForArray encodes a string for use in an array with given delimiter.
func (e *toonEncoder) encodeStringForArray(s string, delim rune) string {
	if e.needsQuotingWithDelim(s, delim) {
		return e.quoteString(s)
	}
	return s
}

// needsQuoting checks if a string needs quoting (using comma as default delimiter).
func (e *toonEncoder) needsQuoting(s string, delim rune) bool {
	return e.needsQuotingWithDelim(s, delim)
}

// needsQuotingWithDelim checks if a string needs quoting with given delimiter.
// A string MUST be quoted when:
// - Empty
// - Has leading/trailing whitespace
// - Equals true, false, or null (case-sensitive)
// - Numeric-like
// - Contains colon, quote, backslash, brackets, braces
// - Contains control characters
// - Contains the active delimiter
// - Equals "-" or starts with hyphen
func (e *toonEncoder) needsQuotingWithDelim(s string, delim rune) bool {
	// Empty string
	if s == "" {
		return true
	}

	// Leading/trailing whitespace
	if s != strings.TrimSpace(s) {
		return true
	}

	// Reserved words
	if s == "true" || s == "false" || s == "null" {
		return true
	}

	// Numeric-like (matches JSON number pattern or leading zeros)
	if e.isNumericLike(s) {
		return true
	}

	// Starts with hyphen or equals "-"
	if strings.HasPrefix(s, "-") {
		return true
	}

	// Contains special characters
	for _, r := range s {
		switch r {
		case ':', '"', '\\', '[', ']', '{', '}':
			return true
		case '\n', '\r', '\t':
			return true
		}
		if r == delim {
			return true
		}
	}

	return false
}

// isNumericLike checks if string looks like a number.
var numericRegex = regexp.MustCompile(`^-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?$`)
var leadingZeroRegex = regexp.MustCompile(`^0\d+$`)

func (e *toonEncoder) isNumericLike(s string) bool {
	return numericRegex.MatchString(s) || leadingZeroRegex.MatchString(s)
}

// quoteString quotes and escapes a string.
// Only valid escapes: \\, \", \n, \r, \t
func (e *toonEncoder) quoteString(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString("\\\\")
		case '"':
			b.WriteString("\\\"")
		case '\n':
			b.WriteString("\\n")
		case '\r':
			b.WriteString("\\r")
		case '\t':
			b.WriteString("\\t")
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}

// formatNumber formats a number in canonical TOON format.
// - No exponent notation
// - No leading zeros except single "0"
// - No trailing fractional zeros
// - Integer when fractional part is zero
func (e *toonEncoder) formatNumber(n float64) string {
	// Handle special case of negative zero
	if n == 0 {
		return "0"
	}

	// Check if it's an integer
	if n == float64(int64(n)) {
		return strconv.FormatInt(int64(n), 10)
	}

	// Format as float with reasonable precision, then remove trailing zeros
	s := strconv.FormatFloat(n, 'f', 2, 64)
	// Remove trailing zeros after decimal point
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}

// formatBool formats a boolean.
func (e *toonEncoder) formatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
