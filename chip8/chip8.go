package chip8

import "fmt"

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
	LoadRegister(firstByte, secondByte byte)
	LoadIndexRegister(firstByte, secondByte byte)
}

type Emulator struct {
	EmulatorStore
}

func (c *Chip8) ClearScreen() {
	for i := range c.Screen {
		c.Screen[i] = 0
	}
}

func (c *Chip8) LoadRegister(firstByte, secondByte byte) {
	register := firstByte & 0x0f
	c.Registers[register] = secondByte
}

func (c *Chip8) LoadIndexRegister(firstByte, secondByte byte) {
	value := uint16(firstByte&0xf)<<8 | uint16(secondByte)
	c.I = value
}

func (e *Emulator) Emulate(firstByte, secondByte byte) {
	switch firstByte >> 4 {
	case 0x0:
		switch secondByte {
		case 0xe0:
			e.ClearScreen()
		}
	case 0x6:
		e.LoadRegister(firstByte, secondByte)
	case 0xa:
		e.LoadIndexRegister(firstByte, secondByte)
	default:
		fmt.Printf("Instruction %x not implemented", uint16(firstByte)<<8|uint16(secondByte))
	}

}
