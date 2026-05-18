package main

import (
	"bufio"
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Clean all whitespace",
			input:    "  hello  world   ",
			expected: []string{"hello", "world"},
		},
		{
			name:     "Clean capitals",
			input:    "  hello  WORLD   ",
			expected: []string{"hello", "world"},
		},
		{
			name:     "Empty case",
			input:    "",
			expected: []string{},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if !reflect.DeepEqual(len(c.expected), len(actual)) {
			t.Errorf("%s: expected length: %v, got: %v", c.name, len(c.expected), len(actual))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if !reflect.DeepEqual(word, expectedWord) {
				t.Errorf("%s: expected word: %v, got: %v", c.name, expectedWord, word)
			}
		}
	}
}

func TestRunREPL(t *testing.T) {
	// Capture stdout
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	config := &Config{}
	input := strings.NewReader("foobar\n")
	scanner := bufio.NewScanner(input)

	runREPL(scanner, config)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)

	if !strings.Contains(buf.String(), "Unknown command") {
		t.Errorf("expected 'Unknown command', got %q", buf.String())
	}
}

// TestRunREPLHelpCommand

//     Feed "help\n" as scanner input
//     Capture stdout (same os.Pipe pattern already used in TestRunREPL)
//     Provide a non-empty command map (from buildCommands) so help output is meaningful
//     Assert output contains "Welcome to the Pokedex!"

// TestRunREPLEmptyInput

//     Feed "\n" (empty line) then EOF
//     Assert REPL exits without panicking (no assertion on output needed — just must not crash)
