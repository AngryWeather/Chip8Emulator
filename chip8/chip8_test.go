package chip8

import (
	"reflect"
	"testing"
)

type Chip8StubStore struct {
	Screen []byte
}

func (c *Chip8StubStore) ClearScreen() {
	for i := range c.Screen {
		c.Screen[i] = 0
	}
}

func TestGetInstruction(t *testing.T) {
	t.Run("with bytes 00 E0 get instruction 00E0", func(t *testing.T) {
		var firstByte byte = 0x00
		var secondByte byte = 0xE0
		instruction := GetInstruction(firstByte, secondByte)
		got := instruction
		var want uint16 = 0x00E0

		assertInstruction(t, got, want)
	})

	t.Run("with bytes 67 34 get instruction 6734", func(t *testing.T) {
		var firstByte byte = 0x67
		var secondByte byte = 0x34

		instruction := GetInstruction(firstByte, secondByte)
		got := instruction
		var want uint16 = 0x6734

		assertInstruction(t, got, want)
	})
}

func TestClearScreen(t *testing.T) {
	chip8 := &Chip8StubStore{
		Screen: make([]byte, 6),
	}
	emulator := Emulator{chip8}

	t.Run("Clears the screen", func(t *testing.T) {
		emulator.Emulate()

		got := chip8.Screen
		want := []byte{0, 0, 0, 0, 0, 0}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func assertInstruction(t testing.TB, got, want uint16) {
	t.Helper()

	if got != want {
		t.Errorf("got %X, want %X", got, want)
	}
}
