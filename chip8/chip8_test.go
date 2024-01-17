package chip8

import (
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
		chip8.Registers = make([]byte, 16)
		chip8.Registers[0x0] = 0
		chip8.PrimaryColor = rl.White
		chip8.SecondaryColor = rl.Black
		chip8.Screen = make([]color.RGBA, 8)
		for i := range chip8.Screen {
			chip8.Screen[i] = rl.Black
		}
		chip8.Width = 8
		chip8.Height = 1
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
		chip8.Memory = []byte{0, 0xff, 0x0f}
		chip8.Registers = make([]byte, 16)
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
		chip8.Memory = []byte{0, 0xff, 0x00}
		chip8.Registers = make([]byte, 16)
		chip8.Registers[0x1] = 0x4

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
	t.Run("instruction 0xd001 - if white on white should turn pixel black and set Vf to 1", func(t *testing.T) {
		chip8 := &Chip8{}
		chip8.Width = 64
		chip8.Height = 6
		chip8.PrimaryColor = rl.White
		chip8.SecondaryColor = rl.Black
		chip8.Screen = make([]color.RGBA, 64*6)
		for i := range chip8.Screen {
			chip8.Screen[i] = rl.White
		}
		chip8.Memory = []byte{0, 0xff, 0x00}
		chip8.Registers = make([]byte, 16)

		chip8.I = 0x1
		emulator := Emulator{EmulatorStore: chip8}

		emulator.Emulate(0xd0, 0x11)

		got := chip8.Screen
		want := []color.RGBA{
			rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.Black, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White,
			rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White,
			rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White,
			rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White,
			rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White,
			rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White, rl.White,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}

		gotVf := chip8.Registers[0xf]
		var wantedVf byte = 0x1
		AssertBytes(t, gotVf, wantedVf)
	})
}

func TestAddToRegister(t *testing.T) {
	t.Run("Instruction 7004 adds 4 to 0 register", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0] = 0
		emulator := Emulator{EmulatorStore: chip}

		emulator.Emulate(0x70, 0x04)
		got := chip.Registers[0]
		var want byte = 0x4

		if got != want {
			t.Errorf("got %x, want %x", got, want)
		}
	})
	t.Run("Instruction 0x76ff adds ff to register 6", func(t *testing.T) {
		chip := NewChip8()

		emulator := Emulator{EmulatorStore: chip}

		emulator.Emulate(0x76, 0xff)

		got := chip.Registers[0x6]
		var want byte = 0xff

		if got != want {
			t.Errorf("got %x, want %x", got, want)
		}
	})
}

func TestSkipNextInstruction(t *testing.T) {
	t.Run("instruction 0x3c00 increments program counter by 2", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0xc] = 0
		emulator := Emulator{EmulatorStore: chip}

		emulator.Emulate(0x3c, 0x00)

		got := chip.Pc
		var want uint16 = 0x202

		AssertAddress(t, got, want)
	})

	t.Run("instruction 0x3c00 doesn't increase the pc if value of register c != 00", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0xc] = 0x10
		emulator := Emulator{EmulatorStore: chip}

		emulator.Emulate(0x3c, 0x00)

		got := chip.Pc
		var want uint16 = 0x200

		AssertAddress(t, got, want)
	})
}

func TestJumpToLocationPlusV0(t *testing.T) {
	t.Run("instruction 0xb321 with v0 register at 0x4 sets pc to 0x325", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0] = 0x4
		emulator := Emulator{EmulatorStore: chip}

		emulator.Emulate(0xb3, 0x21)

		got := chip.Pc
		var want uint16 = 0x325

		AssertAddress(t, got, want)
	})
}

func TestCallAddress(t *testing.T) {
	t.Run("instruction 0x2342 increments stack pointer, puts current pc onto the stack and sets pc to new address", func(t *testing.T) {
		chip := NewChip8()
		emulator := Emulator{EmulatorStore: chip}

		emulator.Emulate(0x23, 0x42)

		got := len(chip.Stack)
		want := 1

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

		addressOnStack := chip.Stack[0]
		wantedAddress := 0x200

		AssertAddress(t, addressOnStack, uint16(wantedAddress))

		pc := chip.Pc
		wantedPc := 0x342

		AssertAddress(t, pc, uint16(wantedPc))
	})
}

func TestSkipIfNotEquals(t *testing.T) {
	t.Run("instruction 0x40ff skips increases pc if register 0 doesn't equal ff", func(t *testing.T) {
		chip := NewChip8()
		emulator := Emulator{EmulatorStore: chip}

		emulator.Emulate(0x40, 0xff)

		got := chip.Pc
		var want uint16 = 0x202

		AssertAddress(t, got, want)
	})
}

func TestSkipEqualRegisters(t *testing.T) {
	t.Run("instruction 0x5010 increases pc if register 0 and 1 have equal values", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 1
		chip.Registers[0x1] = 1
		emulator := Emulator{EmulatorStore: chip}

		emulator.Emulate(0x50, 0x10)

		got := chip.Pc
		var want uint16 = 0x202

		AssertAddress(t, got, want)
	})
}

func TestSkipNotEqualRegisters(t *testing.T) {
	t.Run("instruction 0x9010 increases pc if register 0 != register 1", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0
		chip.Registers[0x1] = 1
		emulator := Emulator{EmulatorStore: chip}

		emulator.Emulate(0x90, 0x10)

		got := chip.Pc
		var want uint16 = 0x202

		AssertAddress(t, got, want)

	})

	t.Run("instruction 0x9010 doesn't increase pc if values of registers are the same", func(t *testing.T) {
		chip := NewChip8()
		emulator := Emulator{EmulatorStore: chip}

		emulator.Emulate(0x90, 0x10)

		got := chip.Pc
		var want uint16 = 0x200

		AssertAddress(t, got, want)
	})
}

func TestReturn(t *testing.T) {
	t.Run("instruction 00EE pops value off the stack and puts it in pc", func(t *testing.T) {
		chip := NewChip8()
		chip.Stack = append(chip.Stack, 0x220)
		emulator := Emulator{EmulatorStore: chip}

		emulator.Emulate(0x00, 0xee)

		got := len(chip.Stack)
		want := 0

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

		pc := chip.Pc
		wantedPc := 0x220

		AssertAddress(t, pc, uint16(wantedPc))

	})
}

func TestVxGetsVy(t *testing.T) {
	t.Run("instruction 0x8010 stores value of register 1 in register 0", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x1] = 0xf

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0x80, 0x10)

		got := chip.Registers[0x0]
		var want byte = 0xf

		if got != want {
			t.Errorf("got %x, want %x", got, want)
		}
	})
}

func TestLoadRegistersFromMemory(t *testing.T) {
	t.Run("instruction 0xf165 loads registers 0 and 1 with values from memory starting at I", func(t *testing.T) {
		chip := NewChip8()
		chip.Memory[0x0] = 0xf
		chip.Memory[0x1] = 0xff
		chip.I = 0x0
		emulator := Emulator{EmulatorStore: chip}

		emulator.Emulate(0xf1, 0x65)

		got := chip.Registers[:2]
		want := []byte{0xf, 0xff}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestVxOrVy(t *testing.T) {
	t.Run("instruction 0x8011 stores result of Vx|Vy in Vx", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0xf  // 00001111
		chip.Registers[0x1] = 0xf1 // 11110001

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0x80, 0x11)

		got := chip.Registers[0x0]
		var want byte = 0xff // 00001111 | 11110000 should return 11111111

		if got != want {
			t.Errorf("got %x, want %x", got, want)
		}
	})
}

func TestVxAndVy(t *testing.T) {
	t.Run("instruction 0x8012 stores result of Vx&Vy in Vx", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0xf  // 00001111
		chip.Registers[0x1] = 0xf1 // 11110001

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0x80, 0x12)

		got := chip.Registers[0x0]
		var want byte = 0x1 // 00000001

		if got != want {
			t.Errorf("got %x, want %x", got, want)
		}
	})
}

func TestVxXorVy(t *testing.T) {
	t.Run("instruction 0x8013 stores result of Vx&Vy in Vx", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0xf  // 00001111
		chip.Registers[0x1] = 0xf1 // 11110001

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0x80, 0x13)

		got := chip.Registers[0x0]
		var want byte = 0xfe // 11111110

		if got != want {
			t.Errorf("got %x, want %x", got, want)
		}
	})
}

func TestVxAddVy(t *testing.T) {
	t.Run("instruction 0x8014 adds value in register 1 to value in register 0 and sets Vf", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x1] = 1
		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0x80, 0x14)

		got := chip.Registers[0x0]
		var want byte = 0x1

		AssertBytes(t, got, want)

		gotVf := chip.Registers[0xf]
		var wantedVf byte = 0x0

		AssertBytes(t, gotVf, wantedVf)
	})
	t.Run("instruction 0x8014 sets vf to 1 on overflow", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0xff
		chip.Registers[0x1] = 1
		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0x80, 0x14)

		got := chip.Registers[0x0]
		var want byte = 0x0

		AssertBytes(t, got, want)

		gotVf := chip.Registers[0xf]
		var wantedVf byte = 0x1

		AssertBytes(t, gotVf, wantedVf)
	})
	t.Run("instruction 0x8014 - set Vf to 0", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0x0
		chip.Registers[0x1] = 0xff
		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0x80, 0x14)

		got := chip.Registers[0x0]
		var want byte = 0xff

		AssertBytes(t, got, want)

		gotVf := chip.Registers[0xf]
		var wantedVf byte = 0x0

		AssertBytes(t, gotVf, wantedVf)
	})
}

func TestVxSubVy(t *testing.T) {
	t.Run("instruction 8015  sets Vx to Vx - Vy and Vf to 1 if V0 > V1", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0xff
		chip.Registers[0x1] = 0xf

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0x80, 0x15)

		got := chip.Registers[0x0]
		var want byte = 0xf0

		AssertBytes(t, got, want)

		gotVf := chip.Registers[0xf]
		var wantedVf byte = 0x1

		AssertBytes(t, gotVf, wantedVf)
	})

	t.Run("instruction 0x8015 - set Vx to 0", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0xf
		chip.Registers[0x1] = 0xff

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0x80, 0x15)

		got := chip.Registers[0x0]
		var want byte = 0x10

		AssertBytes(t, got, want)

		gotVf := chip.Registers[0xf]
		var wantedVf byte = 0x0

		AssertBytes(t, gotVf, wantedVf)
	})
	t.Run("instruction 0x8f15 Vx ix Vf", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0xf] = 0x14
		chip.Registers[0x1] = 0xf

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0x8f, 0x15)

		gotVf := chip.Registers[0xf]
		var wantedVf byte = 0x1

		AssertBytes(t, gotVf, wantedVf)
	})
}

func TestVySubVx(t *testing.T) {
	t.Run("instruction 0x8017 sets Vf to 1 if V1 > V0, sets Vx to Vy - Vx", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0xf
		chip.Registers[0x1] = 0xff

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0x80, 0x17)

		got := chip.Registers[0x0]
		var want byte = 0xf0

		AssertBytes(t, got, want)

		gotVf := chip.Registers[0xf]
		var wantedVf byte = 0x1

		AssertBytes(t, gotVf, wantedVf)
	})
}

func TestVxRightShift(t *testing.T) {
	t.Run("instruction 0x8016 sets Vf to 1 for Vx = 0x3 and divides Vx by 2", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0x3

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0x80, 0x16)

		got := chip.Registers[0x0]
		var want byte = 0x1

		AssertBytes(t, got, want)

		gotVf := chip.Registers[0xf]
		var wantedVf byte = 0x1

		AssertBytes(t, gotVf, wantedVf)
	})

	t.Run("instruction 0x8016 sets Vf to 0 for Vx = 0x2 and divides Vx by 2", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0x2

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0x80, 0x16)

		got := chip.Registers[0x0]
		var want byte = 0x1

		AssertBytes(t, got, want)

		gotVf := chip.Registers[0xf]
		var wantedVf byte = 0x0

		AssertBytes(t, gotVf, wantedVf)
	})
}

func TestVxLeftShift(t *testing.T) {
	t.Run("instruction 0x801e sets Vf to 1 for Vx = 0xf0 and multiplies Vx by 2", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0x80

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0x80, 0x1e)

		got := chip.Registers[0x0]
		var want byte = 0x0

		AssertBytes(t, got, want)

		gotVf := chip.Registers[0xf]
		var wantedVf byte = 0x1

		AssertBytes(t, gotVf, wantedVf)
	})
}

func TestLoadRegistersToMemory(t *testing.T) {
	t.Run("instruction 0xf155 loads values of registers 0 and 1 to memory at I", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers = []byte{0xff, 0xf}
		chip.Memory = []byte{0x0, 0x0}

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0xf1, 0x55)

		got := chip.Memory
		want := []byte{0xff, 0xf}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestStoreBCDRepresentationInMemory(t *testing.T) {
	t.Run("instruction 0xf033 with V0=0xf0 stores 2 at I, 4 at I+1 and 0 at I+2", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0xf0 // 240 in decimal
		chip.Memory = []byte{0x0, 0x0, 0x0}

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0xf0, 0x33)

		got := chip.Memory
		want := []byte{0x2, 0x4, 0x0}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestStoreValueOfVxPlusIInI(t *testing.T) {
	t.Run("instruction f01e with V0=0x1 and I=0x3 stores V0+I in I", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0x1
		chip.I = 0x3

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0xf0, 0x1e)

		got := chip.I
		var want uint16 = 0x4

		AssertAddress(t, got, want)
	})
}

func TestSetDelayTimer(t *testing.T) {
	t.Run("instruction 0xf015 sets delay timer to 0x10 for V0=0x10", func(t *testing.T) {
		chip := NewChip8()
		chip.Registers[0x0] = 0x10

		emulator := Emulator{EmulatorStore: chip}
		emulator.Emulate(0xf0, 0x15)

		got := chip.Timers[0]
		var want byte = 0x10

		AssertBytes(t, got, want)
	})
}

func AssertBytes(t testing.TB, got, want byte) {
	t.Helper()
	if got != want {
		t.Errorf("got %x, want %x", got, want)
	}
}

func AssertAddress(t testing.TB, got, want uint16) {
	t.Helper()
	if got != want {
		t.Errorf("got %x, want %x", got, want)
	}
}
