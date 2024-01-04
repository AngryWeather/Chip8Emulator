package main

import (
	"errors"
	"os"
	"testing"
)

func TestExtractCommandLineArguments(t *testing.T) {
	t.Run("Extract filename from command", func(t *testing.T) {
		os.Args = []string{"go", "main.go", "file.ch8"}
		got, _ := GetFilenameFromCommand(os.Args)
		want := "file.ch8"

		assertFilename(t, got, want)
	})
	t.Run("Return 'test_file.ch8' from command", func(t *testing.T) {
		os.Args = []string{"go", "main.go", "test_file.ch8"}

		got, _ := GetFilenameFromCommand(os.Args)
		want := "test_file.ch8"

		assertFilename(t, got, want)
	})
	t.Run("Return error if no filename was given", func(t *testing.T) {
		os.Args = []string{"go", "main.go"}

		_, err := GetFilenameFromCommand(os.Args)

		if err == nil {
			t.Fatalf("expected an error")
		}

		var got NoFilenameError
		isNoFilenameError := errors.As(err, &got)
		want := NoFilenameError{}

		if !isNoFilenameError {
			t.Fatalf("was not a NoFilenameError got %T", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

	})

}

func assertFilename(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %q", got, want)
	}
}
