package output

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TOONTestSuite is the test suite for the TOON output package.
type TOONTestSuite struct {
	suite.Suite
	enc *toonEncoder
}

// SetupTest is called before each test method.
func (suite *TOONTestSuite) SetupTest() {
	suite.enc = newToonEncoder()
}

// TestFormatNumber_Zero tests that zero is formatted correctly.
func (suite *TOONTestSuite) TestFormatNumber_Zero() {
	result := suite.enc.formatNumber(0)
	suite.Equal("0", result)
}

// TestFormatNumber_NegativeZero tests that negative zero becomes zero.
func (suite *TOONTestSuite) TestFormatNumber_NegativeZero() {
	// Use math.Copysign to create actual negative zero
	negZero := -1.0 * 0.0
	result := suite.enc.formatNumber(negZero)
	suite.Equal("0", result)
}

// TestFormatNumber_Integer tests that integers are formatted without decimal.
func (suite *TOONTestSuite) TestFormatNumber_Integer() {
	result := suite.enc.formatNumber(42)
	suite.Equal("42", result)
}

// TestFormatNumber_NegativeInteger tests negative integers.
func (suite *TOONTestSuite) TestFormatNumber_NegativeInteger() {
	result := suite.enc.formatNumber(-42)
	suite.Equal("-42", result)
}

// TestFormatNumber_Float tests that floats are formatted with minimal decimals.
func (suite *TOONTestSuite) TestFormatNumber_Float() {
	result := suite.enc.formatNumber(3.14)
	suite.Equal("3.14", result)
}

// TestFormatNumber_FloatRounding tests that floats are rounded to 2 decimals.
func (suite *TOONTestSuite) TestFormatNumber_FloatRounding() {
	result := suite.enc.formatNumber(91.30434782608695)
	suite.Equal("91.3", result)
}

// TestFormatNumber_WholeFloat tests that 1.0 becomes 1.
func (suite *TOONTestSuite) TestFormatNumber_WholeFloat() {
	result := suite.enc.formatNumber(1.0)
	suite.Equal("1", result)
}

// TestFormatBool_True tests that true is formatted correctly.
func (suite *TOONTestSuite) TestFormatBool_True() {
	result := suite.enc.formatBool(true)
	suite.Equal("true", result)
}

// TestFormatBool_False tests that false is formatted correctly.
func (suite *TOONTestSuite) TestFormatBool_False() {
	result := suite.enc.formatBool(false)
	suite.Equal("false", result)
}

// TestIsValidUnquotedKey_Simple tests simple valid keys.
func (suite *TOONTestSuite) TestIsValidUnquotedKey_Simple() {
	suite.True(suite.enc.isValidUnquotedKey("name"))
	suite.True(suite.enc.isValidUnquotedKey("myKey"))
	suite.True(suite.enc.isValidUnquotedKey("_private"))
	suite.True(suite.enc.isValidUnquotedKey("key123"))
	suite.True(suite.enc.isValidUnquotedKey("dotted.key"))
}

// TestIsValidUnquotedKey_Invalid tests keys that need quoting.
func (suite *TOONTestSuite) TestIsValidUnquotedKey_Invalid() {
	suite.False(suite.enc.isValidUnquotedKey("123start"))
	suite.False(suite.enc.isValidUnquotedKey("key with space"))
	suite.False(suite.enc.isValidUnquotedKey("key:colon"))
	suite.False(suite.enc.isValidUnquotedKey(""))
	suite.False(suite.enc.isValidUnquotedKey("-hyphen"))
}

// TestNeedsQuoting_Empty tests that empty strings need quoting.
func (suite *TOONTestSuite) TestNeedsQuoting_Empty() {
	suite.True(suite.enc.needsQuoting("", ','))
}

// TestNeedsQuoting_Whitespace tests leading/trailing whitespace.
func (suite *TOONTestSuite) TestNeedsQuoting_Whitespace() {
	suite.True(suite.enc.needsQuoting(" leading", ','))
	suite.True(suite.enc.needsQuoting("trailing ", ','))
	suite.True(suite.enc.needsQuoting(" both ", ','))
}

// TestNeedsQuoting_ReservedWords tests reserved words need quoting.
func (suite *TOONTestSuite) TestNeedsQuoting_ReservedWords() {
	suite.True(suite.enc.needsQuoting("true", ','))
	suite.True(suite.enc.needsQuoting("false", ','))
	suite.True(suite.enc.needsQuoting("null", ','))
	// Case matters - these don't need quoting
	suite.False(suite.enc.needsQuoting("True", ','))
	suite.False(suite.enc.needsQuoting("FALSE", ','))
}

// TestNeedsQuoting_NumericLike tests numeric-like strings.
func (suite *TOONTestSuite) TestNeedsQuoting_NumericLike() {
	suite.True(suite.enc.needsQuoting("123", ','))
	suite.True(suite.enc.needsQuoting("-456", ','))
	suite.True(suite.enc.needsQuoting("3.14", ','))
	suite.True(suite.enc.needsQuoting("0123", ',')) // Leading zero
}

// TestNeedsQuoting_SpecialChars tests special characters.
func (suite *TOONTestSuite) TestNeedsQuoting_SpecialChars() {
	suite.True(suite.enc.needsQuoting("has:colon", ','))
	suite.True(suite.enc.needsQuoting("has\"quote", ','))
	suite.True(suite.enc.needsQuoting("has\\backslash", ','))
	suite.True(suite.enc.needsQuoting("[bracket", ','))
	suite.True(suite.enc.needsQuoting("{brace}", ','))
	suite.True(suite.enc.needsQuoting("has\nnewline", ','))
	suite.True(suite.enc.needsQuoting("has\ttab", ','))
}

// TestNeedsQuoting_Delimiter tests delimiter handling.
func (suite *TOONTestSuite) TestNeedsQuoting_Delimiter() {
	suite.True(suite.enc.needsQuoting("has,comma", ','))
	suite.True(suite.enc.needsQuoting("has|pipe", '|'))
	suite.True(suite.enc.needsQuoting("has\ttab", '\t'))
}

// TestNeedsQuoting_Hyphen tests hyphen handling.
func (suite *TOONTestSuite) TestNeedsQuoting_Hyphen() {
	suite.True(suite.enc.needsQuoting("-", ','))
	suite.True(suite.enc.needsQuoting("-start", ','))
}

// TestNeedsQuoting_Normal tests normal strings that don't need quoting.
func (suite *TOONTestSuite) TestNeedsQuoting_Normal() {
	suite.False(suite.enc.needsQuoting("hello", ','))
	suite.False(suite.enc.needsQuoting("Hello World", ','))
	suite.False(suite.enc.needsQuoting("snake_case", ','))
}

// TestQuoteString_Simple tests simple string quoting.
func (suite *TOONTestSuite) TestQuoteString_Simple() {
	result := suite.enc.quoteString("hello")
	suite.Equal(`"hello"`, result)
}

// TestQuoteString_Escapes tests escape sequences.
func (suite *TOONTestSuite) TestQuoteString_Escapes() {
	suite.Equal(`"has\\backslash"`, suite.enc.quoteString("has\\backslash"))
	suite.Equal(`"has\"quote"`, suite.enc.quoteString(`has"quote`))
	suite.Equal(`"line\nbreak"`, suite.enc.quoteString("line\nbreak"))
	suite.Equal(`"carriage\rreturn"`, suite.enc.quoteString("carriage\rreturn"))
	suite.Equal(`"has\ttab"`, suite.enc.quoteString("has\ttab"))
}

// TestEncodeKey_Simple tests simple key encoding.
func (suite *TOONTestSuite) TestEncodeKey_Simple() {
	suite.Equal("name", suite.enc.encodeKey("name"))
	suite.Equal("myKey", suite.enc.encodeKey("myKey"))
}

// TestEncodeKey_NeedsQuoting tests keys that need quoting.
func (suite *TOONTestSuite) TestEncodeKey_NeedsQuoting() {
	suite.Equal(`"123"`, suite.enc.encodeKey("123"))
	suite.Equal(`"key:value"`, suite.enc.encodeKey("key:value"))
}

// TestEncodeString_Normal tests normal string encoding.
func (suite *TOONTestSuite) TestEncodeString_Normal() {
	suite.Equal("hello", suite.enc.encodeString("hello"))
	suite.Equal("Hello World", suite.enc.encodeString("Hello World"))
}

// TestEncodeString_NeedsQuoting tests strings that need quoting.
func (suite *TOONTestSuite) TestEncodeString_NeedsQuoting() {
	suite.Equal(`""`, suite.enc.encodeString(""))
	suite.Equal(`"true"`, suite.enc.encodeString("true"))
	suite.Equal(`"123"`, suite.enc.encodeString("123"))
}

// TestIsNumericLike tests numeric-like detection.
func (suite *TOONTestSuite) TestIsNumericLike() {
	suite.True(suite.enc.isNumericLike("123"))
	suite.True(suite.enc.isNumericLike("-456"))
	suite.True(suite.enc.isNumericLike("3.14"))
	suite.True(suite.enc.isNumericLike("1e10"))
	suite.True(suite.enc.isNumericLike("0123")) // Leading zeros
	suite.False(suite.enc.isNumericLike("abc"))
	suite.False(suite.enc.isNumericLike("12abc"))
}

// TestWriteKeyValue tests key-value writing.
func (suite *TOONTestSuite) TestWriteKeyValue() {
	suite.enc.writeKeyValue("name", "Alice")
	suite.Equal("name: Alice\n", suite.enc.builder.String())
}

// TestWriteKeyValue_WithIndent tests key-value writing with indentation.
func (suite *TOONTestSuite) TestWriteKeyValue_WithIndent() {
	suite.enc.indent = 1
	suite.enc.writeKeyValue("name", "Alice")
	suite.Equal("  name: Alice\n", suite.enc.builder.String())
}

// TestWriteStringArray tests string array writing.
func (suite *TOONTestSuite) TestWriteStringArray() {
	suite.enc.writeStringArray("tags", []string{"go", "python", "rust"})
	suite.Equal("tags[3]: go,python,rust\n", suite.enc.builder.String())
}

// TestWriteStringArray_Empty tests empty array writing.
func (suite *TOONTestSuite) TestWriteStringArray_Empty() {
	suite.enc.writeStringArray("tags", []string{})
	suite.Equal("tags[0]:\n", suite.enc.builder.String())
}

// TestWriteStringArray_NeedsQuoting tests array with values needing quotes.
func (suite *TOONTestSuite) TestWriteStringArray_NeedsQuoting() {
	suite.enc.writeStringArray("values", []string{"a,b", "normal", "true"})
	suite.Equal(`values[3]: "a,b",normal,"true"`+"\n", suite.enc.builder.String())
}

// TestTOONTestSuite runs all the tests in the suite.
func TestTOONTestSuite(t *testing.T) {
	suite.Run(t, new(TOONTestSuite))
}
