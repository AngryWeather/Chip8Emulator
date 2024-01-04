package main

import (
	"fmt"
	"strings"
)

type NoFilenameError struct{}
type WrongFilenameExtension struct {
	filename string
}

func main() {
}

func (n NoFilenameError) Error() string {
	return "no filename was given"
}

func (w WrongFilenameExtension) Error() string {
	return fmt.Sprintf("filename %s doesn't have .ch8 extension", w.filename)
}

func GetFilenameFromCommand(args []string) (string, error) {
	if len(args) < 2 {
		return "", NoFilenameError{}
	}

	extension := strings.Split(args[1], ".")
	if extension[1] != "ch8" {
		return "", WrongFilenameExtension{filename: args[1]}
	}

	return args[1], nil
}
