package main

import (
	"chip8emulator/chip8"
	"fmt"
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

	filename, err := GetFilenameFromCommand(os.Args)

	if err != nil {
		panic(err)
	}

	program := readFileToBuffer(filename)

	//initialize chip8
	chip := chip8.NewChip8()
	for i := range chip.Screen {
		chip.Screen[i] = rl.Black
	}
	copy(chip.Memory[0x200:], program)
	emulator := chip8.Emulator{EmulatorStore: chip}

	width := int32(1280)
	height := int32(720)
	textureWidth := int32(64)
	textureHeight := int32(32)

	rl.InitWindow(width, height, "Chip8")
	defer rl.CloseWindow()
	checked := rl.Image{Format: rl.UncompressedR8g8b8a8, Width: textureWidth, Height: textureHeight, Mipmaps: 1}

	t := rl.LoadTextureFromImage(&checked)
	chip.Texture = t
	rl.UnloadImage(&checked)
	rl.SetTargetFPS(60)

	target := rl.LoadRenderTexture(width, height)

	for chip.Pc < uint16(len(program)+0x200) && !rl.WindowShouldClose() {
		rl.BeginDrawing()
		if rl.WindowShouldClose() {
			rl.CloseWindow()
		}
		firstByte := chip.Memory[chip.Pc]
		secondByte := chip.Memory[chip.Pc+1]
		fmt.Printf("%x%x\n", firstByte, secondByte)
		emulator.Emulate(firstByte, secondByte)

		if (firstByte == 0x00 && secondByte == 0xe0) || (firstByte>>4 == 0xd) {
			rl.BeginTextureMode(target)
			rl.DrawTexturePro(t, rl.Rectangle{X: 0, Y: 0, Width: float32(textureWidth), Height: float32(textureHeight)}, rl.Rectangle{X: 0, Y: 0, Width: float32(width), Height: float32(height)}, rl.Vector2{X: 0, Y: 0}, 0, rl.White)
			rl.UpdateTexture(chip.Texture, chip.Screen)
			rl.EndTextureMode()
		}
		rl.DrawTexturePro(target.Texture, rl.NewRectangle(0, 0, float32(target.Texture.Width), float32(-target.Texture.Height)), rl.NewRectangle(0, 0, float32(width), float32(height)), rl.NewVector2(0, 0), 0, rl.White)
		chip.Pc += 2

		rl.EndDrawing()
	}
	rl.UnloadTexture(t)
	rl.UnloadRenderTexture(target)
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
