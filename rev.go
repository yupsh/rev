package rev

import (
	"context"
	"fmt"
	"io"
	"strings"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"

	localopt "github.com/yupsh/rev/opt"
)

// Flags represents the configuration options for the rev command
type Flags = localopt.Flags

// Command implementation using StandardCommand abstraction
type command struct {
	yup.StandardCommand[Flags]
}

// Rev creates a new rev command with the given parameters
func Rev(parameters ...any) yup.Command {
	args := opt.Args[string, Flags](parameters...)
	return command{
		StandardCommand: yup.StandardCommand[Flags]{
			Positional: args.Positional,
			Flags:      args.Flags,
			Name:       "rev",
		},
	}
}

func (c command) Execute(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	return c.ProcessFiles(ctx, stdin, stdout, stderr,
		func(ctx context.Context, source yup.InputSource, output io.Writer) error {
			return c.processReader(ctx, source.Reader, output)
		},
	)
}

func (c command) processReader(ctx context.Context, reader io.Reader, output io.Writer) error {
	// Use ProcessLinesSimple to eliminate manual scanner management and context checking
	return yup.ProcessLinesSimple(ctx, reader, output,
		func(ctx context.Context, lineNum int, line string, output io.Writer) error {
			if bool(c.Flags.Separate) {
				// Reverse each word separately
				words := strings.Fields(line)
				reversedWords := make([]string, len(words))
				for i, word := range words {
					// Check for cancellation periodically when processing many words
					if i%100 == 0 {
						if err := yup.CheckContextCancellation(ctx); err != nil {
							return err
						}
					}
					reversedWords[i] = reverseStringWithContext(ctx, word)
				}
				fmt.Fprintln(output, strings.Join(reversedWords, " "))
			} else {
				// Reverse entire line
				fmt.Fprintln(output, reverseStringWithContext(ctx, line))
			}
			return nil
		},
	)
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func reverseStringWithContext(ctx context.Context, s string) string {
	runes := []rune(s)
	length := len(runes)

	// For very long strings, check for cancellation periodically
	for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
		// Check for cancellation every 1000 characters for efficiency
		if i%1000 == 0 {
			if err := yup.CheckContextCancellation(ctx); err != nil {
				// Return partial result on cancellation
				return string(runes)
			}
		}
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func (c command) String() string {
	return fmt.Sprintf("rev %v", c.Positional)
}
