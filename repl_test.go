package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestAvailableCommands(t *testing.T) {
	t.Run("help", func(t *testing.T) {
		r, w, _ := os.Pipe()
		old := os.Stdout
		os.Stdout = w

		commandHelp(&config{})

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		got := buf.String()

		for _, want := range []string{
			"Welcome to the Pokedex!",
			"Usage:",
			"exit: Exit the Pokedex",
			"help: Display help information",
		} {
			if !strings.Contains(got, want) {
				t.Errorf("commandHelp() output missing %q\ngot:\n%s", want, got)
			}
		}
	})

	t.Run("exit", func(t *testing.T) {
		// commandExit calls os.Exit, so run it in a subprocess.
		if os.Getenv("TEST_COMMAND_EXIT") == "1" {
			commandExit(&config{})
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestAvailableCommands/exit")
		cmd.Env = append(os.Environ(), "TEST_COMMAND_EXIT=1")
		var buf bytes.Buffer
		cmd.Stdout = &buf
		_ = cmd.Run()

		got := buf.String()
		want := "Closing the Pokedex... Goodbye!\n"
		if got != want {
			t.Errorf("commandExit() output:\ngot:  %q\nwant: %q", got, want)
		}
	})
}

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "      ",
			expected: []string{},
		},
		{
			input:    "  hello  ",
			expected: []string{"hello"},
		},
		{
			input:    "   hello world   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(c.expected) {
			t.Errorf("Length of output from cleanInput (%d) does not match expected (%d)", len(actual), len(c.expected))
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Output from cleanInput (%s) does not match expected (%s)", word, expectedWord)
			}
		}
	}

}
