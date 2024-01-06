package main

import (
	"chip8emulator/chip8"
	"fmt"
	"io"
	"os"
	"strings"
)

type NoFilenameError struct{}
type WrongFilenameExtension struct {
	filename string
}

func main() {
	filename, err := GetFilenameFromCommand(os.Args)

	if err != nil {
		panic(err)
	}

	program := readFileToBuffer(filename)

	chip := chip8.NewChip8()
	copy(chip.Memory[0x200:], program)
}

func readFileToBuffer(filename string) []byte {
	file, err := os.Open(filename)

	if err != nil {
		panic(fmt.Sprintf("Error reading file %q", filename))
	}

	defer file.Close()

	fileinfo, err := file.Stat()

	if err != nil {
		panic(err)
	}

	filesize := fileinfo.Size()

	b := make([]byte, filesize)

	for {
		_, err := file.Read(b)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
	}
	return b
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
	// chip8 programs need to have .ch8 extension
	if extension[1] != "ch8" {
		return "", WrongFilenameExtension{filename: args[1]}
	}

	return args[1], nil
}
