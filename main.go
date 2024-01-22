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

	font := []byte{
		0xf0, 0x90, 0x90, 0x90, 0xf0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xf0, 0x10, 0xf0, 0x80, 0xf0, // 2
		0xf0, 0x10, 0xf0, 0x10, 0xf0, // 3
		0x90, 0x90, 0xf0, 0x10, 0x10, // 4
		0xf0, 0x80, 0xf0, 0x10, 0xf0, // 5
		0xf0, 0x80, 0xf0, 0x90, 0xf0, // 6
		0xf0, 0x10, 0x20, 0x40, 0x40, // 7
		0xf0, 0x90, 0xf0, 0x90, 0xf0, // 8
		0xf0, 0x90, 0xf0, 0x10, 0xf0, // 9
		0xf0, 0x90, 0xf0, 0x90, 0x90, // A
		0xe0, 0x90, 0xe0, 0x90, 0xe0, // B
		0xf0, 0x80, 0x80, 0x80, 0xf0, // C
		0xe0, 0x90, 0x90, 0x90, 0xe0, // D
		0xf0, 0x80, 0xf0, 0x80, 0xf0, // E
		0xf0, 0x80, 0xf0, 0x80, 0x80, // F
	}

	copy(chip.Memory[0x00:len(font)], font)

	emulator := chip8.Emulator{EmulatorStore: chip}

	width := int32(1280)
	height := int32(720)
	textureWidth := int32(64)
	textureHeight := int32(32)
	colorUIHeight := int32(100)

	rl.InitWindow(width, height+colorUIHeight, "Chip8")
	defer rl.CloseWindow()
	checked := rl.Image{Format: rl.UncompressedR8g8b8a8, Width: textureWidth, Height: textureHeight, Mipmaps: 1}

	primaryColors := [10]rl.Rectangle{}

	centerPos := width/2 - (20*9+60*10)/2

	// create rectangles for primary colors
	for i := 0; i < len(primaryColors); i++ {
		primaryColors[i].X = float32(i*80) + float32(centerPos)
		primaryColors[i].Y = 10
		primaryColors[i].Width = 60
		primaryColors[i].Height = 60
	}

	t := rl.LoadTextureFromImage(&checked)

	rl.SetTextureFilter(t, rl.TextureFilterNearest)
	colors := [10]rl.Color{rl.DarkBlue, rl.White, rl.Red, rl.Blue, rl.Green, rl.Yellow, rl.Brown, rl.Gray,
		rl.Purple, rl.Pink}

	chip.Texture = t
	rl.UnloadImage(&checked)
	rl.SetTargetFPS(60)

	target := rl.LoadRenderTexture(width, height)
	uiTarget := rl.LoadRenderTexture(width, colorUIHeight)

	rl.SetMouseOffset(0, -int(height))
	var colorTint rl.Color = rl.White

	for chip.Pc < uint16(len(program)+0x200) && !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.BeginTextureMode(uiTarget)
		mousePos := rl.GetMousePosition()
		rl.ClearBackground(rl.LightGray)

		// draw rectangles of primaryColor possibilities
		for i := 0; i < len(primaryColors); i++ {
			rl.DrawRectangleRec(primaryColors[i], colors[i])
		}

		// check if mouse is inside rectangle
		for i := 0; i < len(primaryColors); i++ {
			if rl.CheckCollisionPointRec(mousePos, primaryColors[i]) &&
				rl.IsMouseButtonDown(rl.MouseButtonLeft) {
				colorTint = colors[i]
			}
		}

		rl.EndTextureMode()

		// run 10 instructions per frame
		for i := 0; i < 10; i++ {
			if rl.WindowShouldClose() {
				rl.CloseWindow()
			}
			firstByte := chip.Memory[chip.Pc]
			secondByte := chip.Memory[chip.Pc+1]
			emulator.Emulate(firstByte, secondByte)

			if (firstByte == 0x00 && secondByte == 0xe0) || (firstByte>>4 == 0xd) {
				rl.BeginTextureMode(target)
				rl.DrawTexturePro(t, rl.Rectangle{X: 0, Y: 0, Width: float32(textureWidth), Height: float32(textureHeight)}, rl.Rectangle{X: 0, Y: 0, Width: float32(width), Height: float32(height)}, rl.Vector2{X: 0, Y: 0}, 0, rl.White)
				rl.UpdateTexture(chip.Texture, chip.Screen)
				rl.EndTextureMode()
			}

			// these instructions should not increase pc
			if firstByte>>4 != 0x1 && firstByte>>4 != 0x2 {
				chip.Pc += 2
			}
		}

		if chip.Timers[0] > 0 {
			chip.Timers[0] -= 1
		} else {
			chip.Timers[0] = 0
		}

		rl.DrawTexturePro(target.Texture, rl.NewRectangle(0, 0,
			float32(target.Texture.Width), float32(-target.Texture.Height)),
			rl.NewRectangle(0, 0, float32(width), float32(height)),
			rl.NewVector2(0, 0), 0, colorTint)

		// render colors ui
		rl.DrawTexturePro(uiTarget.Texture,
			rl.NewRectangle(0, 0, float32(uiTarget.Texture.Width),
				float32(-uiTarget.Texture.Height)),
			rl.NewRectangle(0, float32(height), float32(width), float32(colorUIHeight)),
			rl.NewVector2(0, 0),
			0,
			rl.White,
		)

		rl.EndDrawing()
	}
	rl.UnloadTexture(t)
	rl.UnloadRenderTexture(target)
	rl.UnloadRenderTexture(uiTarget)
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
