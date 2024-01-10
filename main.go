package main

import (
	"chip8emulator/chip8"
	"fmt"
	"image/color"
	"io"
	"os"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type NoFilenameError struct{}
type WrongFilenameExtension struct {
	filename string
}

func main() {
	// initialize chip8
	// chip8emulator := chip8.NewChip8()
	// emulator := chip8.Emulator{EmulatorStore: chip8emulator}

	width := int32(640)
	height := int32(320)
	textureWidth := int32(64)
	textureHeight := int32(32)

	rl.InitWindow(width, height, "Chip8")
	defer rl.CloseWindow()
	checked := rl.Image{Format: rl.UncompressedR8g8b8a8, Width: textureWidth, Height: textureHeight, Mipmaps: 1}

	t := rl.LoadTextureFromImage(&checked)
	rl.UnloadImage(&checked)

	rl.SetTargetFPS(60)
	var pixels = make([]color.RGBA, textureWidth*textureHeight)
	for i := range pixels {
		pixels[i] = rl.DarkPurple
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		rl.DrawTexturePro(t, rl.Rectangle{X: 0, Y: 0, Width: float32(textureWidth), Height: float32(textureHeight)}, rl.Rectangle{X: 0, Y: 0, Width: float32(width), Height: float32(height)}, rl.Vector2{X: 0, Y: 0}, 0, rl.White)
		rl.UpdateTexture(t, pixels)

		rl.EndDrawing()
	}

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
