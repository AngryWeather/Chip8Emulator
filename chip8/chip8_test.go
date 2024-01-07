package chip8

import "testing"

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

func assertInstruction(t testing.TB, got, want uint16) {
	t.Helper()

	if got != want {
		t.Errorf("got %X, want %X", got, want)
	}
}
