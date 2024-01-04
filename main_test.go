package main

import (
	"errors"
	"os"
	"testing"
)

func TestExtractCommandLineArguments(t *testing.T) {
	t.Run("Extract filename from command", func(t *testing.T) {
		os.Args = []string{"go run main.go", "file.ch8"}
		got, err := GetFilenameFromCommand(os.Args)
		want := "file.ch8"

		if err != nil {
			t.Fatalf("didn't expect an error, got %T", err)
		}

		assertFilename(t, got, want)
	})
	t.Run("Return 'test_file.ch8' from command", func(t *testing.T) {
		os.Args = []string{"go run main.go", "test_file.ch8"}

		got, err := GetFilenameFromCommand(os.Args)
		want := "test_file.ch8"

		if err != nil {
			t.Fatalf("didn't expect an error, got %T", err)
		}

		assertFilename(t, got, want)
	})
	t.Run("Return error if no filename was given", func(t *testing.T) {
		os.Args = []string{"go run main.go"}

		_, err := GetFilenameFromCommand(os.Args)

		assertErrorExpected(t, err)

		var got NoFilenameError
		isNoFilenameError := errors.As(err, &got)
		want := NoFilenameError{}

		assertIsError(t, isNoFilenameError, want, err)

		assertError(t, got, want)
	})

	t.Run("Return error if file extension is not '.ch8'", func(t *testing.T) {
		os.Args = []string{"go run main.go", "test_file.txt"}

		_, err := GetFilenameFromCommand(os.Args)

		assertErrorExpected(t, err)

		var got WrongFilenameExtension
		isWrongFilenameExtension := errors.As(err, &got)
		want := WrongFilenameExtension{filename: os.Args[1]}

		assertIsError(t, isWrongFilenameExtension, want, err)

		assertError(t, got, want)
	})
}

func assertIsError(t testing.TB, b bool, want, err error) {
	if !b {
		t.Fatalf("was not a %T, got %T", want, err)
	}
}

func assertError(t testing.TB, got, want error) {
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertErrorExpected(t testing.TB, err error) {
	if err == nil {
		t.Fatalf("expected an error")
	}
}

func assertFilename(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %q", got, want)
	}
}
