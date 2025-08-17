package rev

import (
	"context"
	"strings"
	"testing"
	"time"

	yup "github.com/yupsh/framework"

	"github.com/yupsh/rev/opt"
)

func TestRev(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		flags    []any
		expected string
	}{
		{
			name:     "reverse single line",
			input:    "hello world",
			flags:    []any{},
			expected: "dlrow olleh\n",
		},
		{
			name:     "reverse multiple lines",
			input:    "hello\nworld\n",
			flags:    []any{},
			expected: "olleh\ndlrow\n",
		},
		{
			name:     "reverse empty line",
			input:    "",
			flags:    []any{},
			expected: "",
		},
		{
			name:     "reverse line with spaces",
			input:    "  hello world  ",
			flags:    []any{},
			expected: "  dlrow olleh  \n",
		},
		{
			name:     "reverse with unicode characters",
			input:    "hello 世界",
			flags:    []any{},
			expected: "界世 olleh\n",
		},
		{
			name:     "reverse each word separately",
			input:    "hello world test",
			flags:    []any{opt.Separate},
			expected: "olleh dlrow tset\n",
		},
		{
			name:     "reverse words with multiple spaces",
			input:    "hello   world",
			flags:    []any{opt.Separate},
			expected: "olleh dlrow\n",
		},
		{
			name:     "reverse empty string with separate flag",
			input:    "",
			flags:    []any{opt.Separate},
			expected: "",
		},
		{
			name:     "reverse single word with separate flag",
			input:    "hello",
			flags:    []any{opt.Separate},
			expected: "olleh\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Rev(tt.flags...)

			var output strings.Builder
			var stderr strings.Builder

			ctx := context.Background()
			err := cmd.Execute(ctx, strings.NewReader(tt.input), &output, &stderr)

			if err != nil {
				t.Fatalf("Execute failed: %v\nStderr: %s", err, stderr.String())
			}

			result := output.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestRevContextCancellation(t *testing.T) {
	// Create a very long string to test context cancellation
	longLine := strings.Repeat("a", 10000)
	input := strings.Repeat(longLine+"\n", 100)

	// Create a context that will be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	cmd := Rev()

	var output strings.Builder
	var stderr strings.Builder

	// Cancel context immediately
	cancel()

	err := cmd.Execute(ctx, strings.NewReader(input), &output, &stderr)

	// Should detect cancellation and return error
	if err == nil {
		t.Error("Expected context cancellation error, got nil")
	}

	if !strings.Contains(err.Error(), "context canceled") && !strings.Contains(err.Error(), "context cancelled") {
		t.Errorf("Expected context cancellation error, got: %v", err)
	}
}

func TestRevContextCancellationTimeout(t *testing.T) {
	// Create a very long string
	longLine := strings.Repeat("a", 50000)
	input := strings.Repeat(longLine+"\n", 50)

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	cmd := Rev()

	var output strings.Builder
	var stderr strings.Builder

	err := cmd.Execute(ctx, strings.NewReader(input), &output, &stderr)

	// Should timeout and return error
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestRevWithFiles(t *testing.T) {
	// Test with positional arguments (files)
	cmd := Rev("testfile1", "testfile2")

	// Since these files don't exist, we expect an error
	var output strings.Builder
	var stderr strings.Builder

	ctx := context.Background()
	err := cmd.Execute(ctx, nil, &output, &stderr)

	// Should get file not found error
	if err == nil {
		t.Error("Expected file not found error, got nil")
	}
}

func TestRevInterface(t *testing.T) {
	// Verify that Rev command implements yup.Command interface
	var _ yup.Command = Rev()
}

func BenchmarkRevSimple(b *testing.B) {
	input := "hello world test line\n"
	cmd := Rev()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var output strings.Builder
		var stderr strings.Builder
		cmd.Execute(ctx, strings.NewReader(input), &output, &stderr)
	}
}

func BenchmarkRevSeparate(b *testing.B) {
	input := "hello world test line\n"
	cmd := Rev(opt.Separate)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var output strings.Builder
		var stderr strings.Builder
		cmd.Execute(ctx, strings.NewReader(input), &output, &stderr)
	}
}

func BenchmarkRevLongLine(b *testing.B) {
	input := strings.Repeat("word ", 1000) + "\n"
	cmd := Rev()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var output strings.Builder
		var stderr strings.Builder
		cmd.Execute(ctx, strings.NewReader(input), &output, &stderr)
	}
}

// Example tests for documentation
func ExampleRev() {
	cmd := Rev()
	ctx := context.Background()

	input := strings.NewReader("hello world\n")
	cmd.Execute(ctx, input, &strings.Builder{}, &strings.Builder{})
	// Output would be: dlrow olleh
}

func ExampleRev_separate() {
	cmd := Rev(opt.Separate)
	ctx := context.Background()

	input := strings.NewReader("hello world\n")
	var output strings.Builder
	cmd.Execute(ctx, input, &output, &strings.Builder{})
	// Output would be: olleh dlrow
}
