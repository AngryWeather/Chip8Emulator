package chip8

import (
	"fmt"
	"image/color"
	"reflect"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Chip8StubStore struct {
	Screen []byte
}

func TestClearScreen(t *testing.T) {
	t.Run("Clears the screen", func(t *testing.T) {
		chip8 := &Chip8{}
		emulator := Emulator{EmulatorStore: chip8}
		chip8.PrimaryColor = rl.White
		chip8.SecondaryColor = rl.Black

		chip8.Screen = []color.RGBA{
			rl.Black,
			rl.Black,
			rl.White,
			rl.White,
		}

		emulator.Emulate(0x00, 0xe0)

		got := chip8.Screen
		want := []color.RGBA{
			rl.Black,
			rl.Black,
			rl.Black,
			rl.Black,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

}

func TestLoadRegister(t *testing.T) {
	t.Run("Load register 0x1 with value 0xff", func(t *testing.T) {
		chip8 := &Chip8{}
		emulator := Emulator{EmulatorStore: chip8}
		chip8.Registers = []byte{0, 0, 0}

		emulator.Emulate(0x61, 0xff)

		got := chip8.Registers
		want := []byte{0, 0xff, 0}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestLoadIndexRegister(t *testing.T) {
	t.Run("Load index register with 0x231", func(t *testing.T) {
		chip8 := &Chip8{}
		emulator := Emulator{EmulatorStore: chip8}

		emulator.Emulate(0xa2, 0x31)

		got := chip8.I
		var want uint16 = 0x231

		AssertAddress(t, got, want)
	})
}

func TestJumpInstruction(t *testing.T) {
	t.Run("Instruction with bytes 0x13 and 0x45 sets program counter to 0x345", func(t *testing.T) {
		chip8 := &Chip8{}
		emulator := Emulator{EmulatorStore: chip8}

		emulator.Emulate(0x13, 0x45)

		got := chip8.Pc
		var want uint16 = 0x345

		AssertAddress(t, got, want)
	})
}

func TestDrawInstruction(t *testing.T) {
	t.Run("Instruction 0xd001 changes screen", func(t *testing.T) {
		chip8 := &Chip8{}
		chip8.Memory = []byte{0, 0xff}
		chip8.Registers = []byte{0, 0, 0}
		chip8.Registers[0x0] = 0
		chip8.PrimaryColor = rl.White
		chip8.SecondaryColor = rl.Black
		chip8.Screen = make([]color.RGBA, 8)
		chip8.Width = 8
		chip8.I = 0x1
		emulator := Emulator{EmulatorStore: chip8}

		emulator.Emulate(0xd0, 0x01)

		got := chip8.Screen
		want := []color.RGBA{
			rl.White,
			rl.White,
			rl.White,
			rl.White,
			rl.White,
			rl.White,
			rl.White,
			rl.White,
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}

	})

	t.Run("Instruction 0xd002 Draws two bytes", func(t *testing.T) {
		chip8 := &Chip8{}
		chip8.Width = 12
		chip8.Height = 2
		chip8.PrimaryColor = rl.White
		chip8.SecondaryColor = rl.Black
		chip8.Screen = make([]color.RGBA, 12*2)
		for i := range chip8.Screen {
			chip8.Screen[i] = rl.Black
		}
		fmt.Printf("%v", rl.Black)
		chip8.Memory = []byte{0, 0xff, 0x0f}
		chip8.Registers = []byte{0, 0, 0}
		chip8.Registers[0x0] = 0
		chip8.I = 0x1
		emulator := Emulator{EmulatorStore: chip8}

		emulator.Emulate(0xd0, 0x02)

		got := chip8.Screen
		want := []color.RGBA{
			// first row
			rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.Black, rl.Black, rl.Black, rl.Black,
			// second row
			rl.Black, rl.Black, rl.Black, rl.Black, rl.White, rl.White, rl.White, rl.White, rl.Black, rl.Black, rl.Black, rl.Black,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Instruction 0xd011 Draws one byte - test for overflow position", func(t *testing.T) {
		chip8 := &Chip8{}
		chip8.Width = 64
		chip8.Height = 6
		chip8.PrimaryColor = rl.White
		chip8.SecondaryColor = rl.Black
		chip8.Screen = make([]color.RGBA, 64*6)
		for i := range chip8.Screen {
			chip8.Screen[i] = rl.Black
		}
		fmt.Printf("%v", rl.Black)
		chip8.Memory = []byte{0, 0xff, 0x00}
		chip8.Registers = []byte{0, 0x4, 0}

		chip8.I = 0x1
		emulator := Emulator{EmulatorStore: chip8}

		emulator.Emulate(0xd0, 0x11)

		got := chip8.Screen
		want := []color.RGBA{
			rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black,
			rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black,
			rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black,
			rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black,
			rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black,
			rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func AssertAddress(t testing.TB, got, want uint16) {
	t.Helper()
	if got != want {
		t.Errorf("got %x, want %x", got, want)
	}
}
