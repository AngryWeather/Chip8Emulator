package chip8

type Chip8 struct {
	Memory    []byte
	Registers []byte
	Timers    []byte
	Stack     []uint16
	Screen    []uint8
	Pc        uint16
	Sp        uint8
	I         uint16
}

func NewChip8() *Chip8 {
	chip := &Chip8{
		Memory:    make([]byte, 4096),
		Registers: make([]byte, 16),
		Timers:    make([]byte, 2),
		Stack:     make([]uint16, 16),
		Screen:    make([]uint8, 64*32),
		Pc:        0x200,
		Sp:        0,
		I:         0,
	}

	return chip
}

type EmulatorStore interface {
	ClearScreen()
}

type Emulator struct {
	EmulatorStore
}

func (e *Emulator) LoadRegister() {

}

func (c *Chip8) ClearScreen() {
	for i := range c.Screen {
		c.Screen[i] = 0
	}
}

func GetInstruction(firstByte, secondByte byte) uint16 {
	return uint16(firstByte)<<8 | uint16(secondByte)
}

func (e *Emulator) Emulate() {
	e.EmulatorStore.ClearScreen()
}
