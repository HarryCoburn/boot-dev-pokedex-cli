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
	runREPL(scanner, map[string]cliCommand{}, config)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)

	if !strings.Contains(buf.String(), "Unknown command") {
		t.Errorf("expected 'Unknown command', got %q", buf.String())
	}
}
