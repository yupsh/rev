package command_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/rev"
)

// ==============================================================================
// Test Basic Functionality
// ==============================================================================

func TestRev_SingleLine(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("hello").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"olleh"})
}

func TestRev_MultipleLines(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("hello", "world", "test").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"olleh",
		"dlrow",
		"tset",
	})
}

func TestRev_EmptyInput(t *testing.T) {
	result := run.Quick(command.Rev())

	assertion.NoError(t, result.Err)
	assertion.Empty(t, result.Stdout)
}

func TestRev_EmptyLine(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{""})
}

func TestRev_MultipleEmptyLines(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("", "", "").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"", "", ""})
}

// ==============================================================================
// Test With Spaces and Special Characters
// ==============================================================================

func TestRev_WithSpaces(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("abc 123").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"321 cba"})
}

func TestRev_LeadingSpaces(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("   hello").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"olleh   "})
}

func TestRev_TrailingSpaces(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("hello   ").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"   olleh"})
}

func TestRev_OnlySpaces(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("     ").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"     "})
}

func TestRev_Tabs(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("a\tb\tc").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"c\tb\ta"})
}

func TestRev_MixedWhitespace(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("a  b\tc  d").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"d  c\tb  a"})
}

// ==============================================================================
// Test Unicode Support
// ==============================================================================

func TestRev_Unicode_Japanese(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("日本語").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"語本日"})
}

func TestRev_Unicode_Mixed(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("Hello世界123").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"321界世olleH"})
}

func TestRev_Unicode_Emoji(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("Hello👋World🌍").
		Run()

	assertion.NoError(t, result.Err)
	// Emojis should be reversed as single units
	assertion.Contains(t, result.Stdout, "🌍")
	assertion.Contains(t, result.Stdout, "👋")
}

func TestRev_Unicode_Arabic(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("مرحبا").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 1)
	assertion.True(t, len(result.Stdout[0]) > 0, "should have output")
}

func TestRev_Unicode_Greek(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("Ελληνικά").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"άκινηλλΕ"})
}

// ==============================================================================
// Test Special Characters
// ==============================================================================

func TestRev_Punctuation(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("Hello, World!").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"!dlroW ,olleH"})
}

func TestRev_SpecialCharacters(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("!@#$%^&*()").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{")(*&^%$#@!"})
}

func TestRev_Brackets(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("[{()}]").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"]})({["})
}

func TestRev_Quotes(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines(`"hello" 'world'`).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{`'dlrow' "olleh"`})
}

// ==============================================================================
// Test Edge Cases
// ==============================================================================

func TestRev_SingleCharacter(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("a").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"a"})
}

func TestRev_TwoCharacters(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("ab").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"ba"})
}

func TestRev_Palindrome(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("racecar", "level", "noon").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"racecar",
		"level",
		"noon",
	})
}

func TestRev_VeryLongLine(t *testing.T) {
	longLine := strings.Repeat("a", 10000)
	expected := strings.Repeat("a", 10000) // Palindrome

	result := run.Command(command.Rev()).
		WithStdinLines(longLine).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{expected})
}

func TestRev_VeryLongLine_NonPalindrome(t *testing.T) {
	longLine := "start" + strings.Repeat("middle", 2000) + "end"

	result := run.Command(command.Rev()).
		WithStdinLines(longLine).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 1)
	// Check it starts with "dne" and ends with "trats"
	assertion.True(t, strings.HasPrefix(result.Stdout[0], "dne"), "should start with 'dne'")
	assertion.True(t, strings.HasSuffix(result.Stdout[0], "trats"), "should end with 'trats'")
}

func TestRev_ManyLines(t *testing.T) {
	lines := make([]string, 1000)
	expected := make([]string, 1000)
	for i := range lines {
		lines[i] = "line"
		expected[i] = "enil"
	}

	result := run.Command(command.Rev()).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 1000)
	assertion.Lines(t, result.Stdout, expected)
}

func TestRev_MixedContent(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines(
			"normal",
			"",
			"with spaces",
			"日本語",
			"123",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"lamron",
		"",
		"secaps htiw",
		"語本日",
		"321",
	})
}

// ==============================================================================
// Test Numbers and Alphanumeric
// ==============================================================================

func TestRev_Numbers(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("12345").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"54321"})
}

func TestRev_Alphanumeric(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("abc123def456").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"654fed321cba"})
}

func TestRev_HexString(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("0xDEADBEEF").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"FEEBDAEDx0"})
}

// ==============================================================================
// Test Error Handling
// ==============================================================================

func TestRev_InputError(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinError(errors.New("read failed")).
		Run()

	assertion.ErrorContains(t, result.Err, "read failed")
}

func TestRev_OutputError(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("test").
		WithStdoutError(errors.New("write failed")).
		Run()

	assertion.ErrorContains(t, result.Err, "write failed")
}

// ==============================================================================
// Table-Driven Tests
// ==============================================================================

func TestRev_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "single word",
			input:    []string{"hello"},
			expected: []string{"olleh"},
		},
		{
			name:     "multiple words",
			input:    []string{"hello", "world"},
			expected: []string{"olleh", "dlrow"},
		},
		{
			name:     "with spaces",
			input:    []string{"hello world"},
			expected: []string{"dlrow olleh"},
		},
		{
			name:     "empty line",
			input:    []string{""},
			expected: []string{""},
		},
		{
			name:     "numbers",
			input:    []string{"123"},
			expected: []string{"321"},
		},
		{
			name:     "unicode",
			input:    []string{"日本語"},
			expected: []string{"語本日"},
		},
		{
			name:     "palindrome",
			input:    []string{"racecar"},
			expected: []string{"racecar"},
		},
		{
			name:     "special chars",
			input:    []string{"!@#$"},
			expected: []string{"$#@!"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := run.Command(command.Rev()).
				WithStdinLines(tt.input...).
				Run()

			assertion.NoError(t, result.Err)
			assertion.Lines(t, result.Stdout, tt.expected)
		})
	}
}

// ==============================================================================
// Test Line Endings and Whitespace Preservation
// ==============================================================================

func TestRev_PreservesLineStructure(t *testing.T) {
	// Each line should be reversed independently
	result := run.Command(command.Rev()).
		WithStdinLines("abc", "def", "ghi").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 3)
	assertion.Equal(t, result.Stdout[0], "cba", "line 1")
	assertion.Equal(t, result.Stdout[1], "fed", "line 2")
	assertion.Equal(t, result.Stdout[2], "ihg", "line 3")
}

func TestRev_EmptyLinesBetween(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("abc", "", "def").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"cba",
		"",
		"fed",
	})
}

// ==============================================================================
// Test Real-World Scenarios
// ==============================================================================

func TestRev_FilePath(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("/path/to/file.txt").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"txt.elif/ot/htap/"})
}

func TestRev_URL(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("https://example.com/path").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"htap/moc.elpmaxe//:sptth"})
}

func TestRev_Email(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("user@example.com").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"moc.elpmaxe@resu"})
}

func TestRev_SentenceWithPunctuation(t *testing.T) {
	result := run.Command(command.Rev()).
		WithStdinLines("Hello, how are you?").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"?uoy era woh ,olleH"})
}

// ==============================================================================
// Test Flags (for coverage)
// ==============================================================================

func TestRev_WithSeparateFlag(t *testing.T) {
	// The Separate flag is defined but not currently used in the implementation
	// This test ensures the flag can be set without errors
	result := run.Command(command.Rev(command.Separate)).
		WithStdinLines("hello").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"olleh"})
}

func TestRev_WithNoSeparateFlag(t *testing.T) {
	// The NoSeparate flag is defined but not currently used in the implementation
	// This test ensures the flag can be set without errors
	result := run.Command(command.Rev(command.NoSeparate)).
		WithStdinLines("hello").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"olleh"})
}

