package main

import (
	"chip8emulator/chip8"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const dropText = "Drop .ch8 file"
const dropTextFontSize = 120

type NoFilenameError struct{}
type WrongFilenameExtension struct {
	filename string
}

func main() {

	//initialize chip8
	chip := chip8.NewChip8()
	for i := range chip.Screen {
		chip.Screen[i] = rl.Black
	}

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
	height := int32(770)
	textureWidth := int32(64)
	textureHeight := int32(32)
	colorUIHeight := int32(100)
	topUIHeight := int32(50)

	rl.InitWindow(width, height+colorUIHeight, "Chip8")

	// set gui style
	gui.LoadStyle("assets/style_terminal.rgs")
	uiColor := rl.NewColor(0x16, 0x13, 0x13, 0xff)
	uiTextColor := rl.NewColor(0x38, 0xf6, 0x20, 0xff)

	defer rl.CloseWindow()
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()
	sound := rl.LoadSound("assets/beep.wav")
	defer rl.UnloadSound(sound)

	// create image for chip texture
	checked := rl.Image{Format: rl.UncompressedR8g8b8a8, Width: textureWidth, Height: textureHeight, Mipmaps: 1}

	primaryColors := [10]rl.Rectangle{}
	// center colors
	centerPos := width/2 - (20*9+60*10)/2

	// create rectangles for primary colors
	for i := 0; i < len(primaryColors); i++ {
		primaryColors[i].X = float32(i*80) + float32(centerPos)
		// center vertically
		primaryColors[i].Y = float32(colorUIHeight)/2 - 30
		primaryColors[i].Width = 60
		primaryColors[i].Height = 60
	}

	t := rl.LoadTextureFromImage(&checked)

	rl.SetTextureFilter(t, rl.TextureFilterNearest)
	colors := [10]rl.Color{rl.Gold, rl.White, rl.Red, rl.Blue, rl.Green, rl.Yellow, uiTextColor, rl.Orange,
		rl.Purple, rl.Pink}

	chip.Texture = t
	rl.UnloadImage(&checked)
	rl.SetTargetFPS(60)

	target := rl.LoadRenderTexture(width, height)
	uiTarget := rl.LoadRenderTexture(width, colorUIHeight)
	topUITarget := rl.LoadRenderTexture(width, topUIHeight)

	var colorTint rl.Color = rl.White

	// load pixel font
	pixelFont := rl.LoadFont("assets/pixelplay.png")
	defer rl.UnloadFont(pixelFont)

	// calculate center position on the screen of the drop file text
	centerDropText := rl.MeasureTextEx(pixelFont, dropText, dropTextFontSize, 4)
	centerDropTextX := centerDropText.X / 2
	centerDropTextY := centerDropText.Y / 2

	// wait for player to drop file
	for !rl.IsFileDropped() && !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.DrawTextEx(pixelFont, dropText, rl.Vector2{
			X: float32(width/2 - int32(centerDropTextX)),
			Y: float32((height+colorUIHeight)/2 - int32(centerDropTextY))}, dropTextFontSize, 4, rl.White)
		rl.EndDrawing()

	}

	var program []byte
	if rl.IsFileDropped() {
		program = readFileToBuffer(rl.LoadDroppedFiles()[0])
		copy(chip.Memory[0x200:], program)

	}

	var tickrateSpinner int32 = 10

	for chip.Pc < uint16(len(program)+0x200) && !rl.WindowShouldClose() {

		rl.BeginDrawing()

		// render topUI to buffer
		rl.BeginTextureMode(topUITarget)
		rl.ClearBackground(uiColor)
		rl.SetMouseOffset(0, 0)
        // create spinenr for changing tickrate
		tickrateSpinner = gui.Spinner(rl.NewRectangle(0, 0, 100, 50), "tickrate", &tickrateSpinner, 1, 1000, true)

		rl.EndTextureMode()

		rl.BeginTextureMode(uiTarget)
		rl.SetMouseOffset(0, -int(height))
		mousePos := rl.GetMousePosition()

		rl.ClearBackground(uiColor)

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
		for i := 0; i < int(tickrateSpinner); i++ {
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

		for t := range chip.Timers {
			if chip.Timers[t] > 0 {
				// chip.Timers[1] is a sound timer so if it's greater than 0 play sound
				if t == 1 {
					if !rl.IsSoundPlaying(sound) {
						rl.PlaySound(sound)
					}
				}
				chip.Timers[t] -= 1
			} else {
				// deactivate sound timer and stop sound
				if t == 1 {
					rl.StopSound(sound)
				}
				chip.Timers[t] = 0

			}
		}

		// render topUI
		rl.DrawTexturePro(topUITarget.Texture,
			rl.NewRectangle(0, 0, float32(topUITarget.Texture.Width),
				float32(-topUITarget.Texture.Height)),
			rl.NewRectangle(0, 0, float32(width), float32(topUIHeight)),
			rl.NewVector2(0, 0), 0, rl.White)

		// render texture target
		rl.DrawTexturePro(target.Texture, rl.NewRectangle(0, 0,
			float32(target.Texture.Width), float32(-target.Texture.Height)),
			rl.NewRectangle(0, float32(topUIHeight), float32(width), float32(height-topUITarget.Texture.Height)),
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

func readFileToBuffer(filepath string) []byte {
	file, err := os.Open(filepath)

	if err != nil {
		panic(fmt.Sprintf("Error opening file %q", filepath))
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

func GetFilenameFromGUI(path string) string {
	path = filepath.ToSlash(path)
	splits := strings.Split(path, "/")
	return splits[len(splits)-1]
}
