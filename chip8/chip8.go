package chip8

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Chip8 struct {
	Memory    []byte
	Registers []byte
	Timers    []byte
	Stack     []uint16
	Screen    []color.RGBA
	Width     byte
	Height    byte
	Pc        uint16
	Sp        uint8
	I         uint16
	Texture   rl.Texture2D
}

func NewChip8() *Chip8 {
	chip := &Chip8{
		Memory:    make([]byte, 4096),
		Registers: make([]byte, 16),
		Timers:    make([]byte, 2),
		Stack:     make([]uint16, 16),
		Screen:    make([]color.RGBA, 64*32),
		Width:     64,
		Height:    32,
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
	JumpToInstruction(firstByte, secondByte byte)
	Draw(firstByte, secondByte byte)
}

type Emulator struct {
	EmulatorStore
}

// ClearScreen clears the screen by setting all pixels to 0.
func (c *Chip8) ClearScreen() {
	for i := range c.Screen {
		c.Screen[i] = rl.Black
	}
}

// LoadRegister loads secondByte into register.
func (c *Chip8) LoadRegister(firstByte, secondByte byte) {
	// register name is represented by the last 4 bits of the first byte.
	register := firstByte & 0x0f
	c.Registers[register] = secondByte
}

func get12BitValue(firstByte, secondByte byte) uint16 {
	// value is the last 4 bits of the first byte and the last byte.
	// shift first byte to fit the second byte and connect them together.
	return uint16(firstByte&0xf)<<8 | uint16(secondByte)
}

// LoadIndexRegister loads 12 bits into index register.
func (c *Chip8) LoadIndexRegister(firstByte, secondByte byte) {
	c.I = get12BitValue(firstByte, secondByte)
}

// JumpToInstruction sets the program counter to the new value.
func (c *Chip8) JumpToInstruction(firstByte, secondByte byte) {
	c.Pc = get12BitValue(firstByte, secondByte)
}

func (c *Chip8) Draw(firstByte, secondByte byte) {
	bytesToRead := secondByte & 0xf
	x := c.Registers[firstByte&0xf]
	y := c.Registers[secondByte>>4]

	for i := c.I; i < (uint16(bytesToRead) + c.I); i++ {
		var currentByte byte = c.Memory[i]
		var color color.RGBA

		// check each bit in the current byte
		for j := 0; j < 8; j++ {
			pixel := currentByte >> 7 & 0x1
			if pixel == 1 {
				color = rl.White
			} else {
				color = rl.Black
			}

			// position in 1D array is based on x, y and width
			var position int = int(x) + (int(y) * int(c.Width))
			c.Screen[position] = color

			// shift byte to access next bit from left
			currentByte = currentByte << 1
			// increase x to draw in the next x coordinate
			x += 1
		}
		// increase y to move down
		y += 1
		// reset x
		x = c.Registers[firstByte&0xf]

	}
}

func (e *Emulator) Emulate(firstByte, secondByte byte) {
	switch firstByte >> 4 {
	case 0x0:
		switch secondByte {
		case 0xe0:
			e.ClearScreen()
		}
	case 0x1:
		e.JumpToInstruction(firstByte, secondByte)
	case 0x6:
		e.LoadRegister(firstByte, secondByte)
	case 0xa:
		e.LoadIndexRegister(firstByte, secondByte)
	case 0xd:
		e.Draw(firstByte, secondByte)
	default:
		fmt.Printf("Instruction %x not implemented", uint16(firstByte)<<8|uint16(secondByte))
	}

}
