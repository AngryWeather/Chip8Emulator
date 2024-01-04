package main

import (
	"os"
	"testing"
)

func TestExtractCommandLineArguments(t *testing.T) {
	t.Run("Extract filename from command", func(t *testing.T) {
		os.Args = []string{"go", "main.go", "file.ch8"}
		got := GetFilenameFromCommand(os.Args)
		want := "file.ch8"

		if got != want {
			t.Errorf("got %v, want %q", got, want)
		}
	})
	t.Run("Return 'test_file.ch8' from command", func(t *testing.T) {
		os.Args = []string{"go", "main.go", "test_file.ch8"}

		got := GetFilenameFromCommand(os.Args)
		want := "test_file.ch8"

		if got != want {
			t.Errorf("got %v, want %q", got, want)
		}
	})

}
