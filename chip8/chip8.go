package chip8

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Chip8 struct {
	Memory         []byte
	Registers      []byte
	Timers         []byte
	Stack          []uint16
	Screen         []color.RGBA
	Width          byte
	Height         byte
	Pc             uint16
	Sp             uint8
	I              uint16
	Texture        rl.Texture2D
	PrimaryColor   color.RGBA
	SecondaryColor color.RGBA
}

func NewChip8() *Chip8 {
	chip := &Chip8{
		Memory:         make([]byte, 4096),
		Registers:      make([]byte, 16),
		Timers:         make([]byte, 2),
		Stack:          make([]uint16, 0, 16),
		Screen:         make([]color.RGBA, 64*32),
		Width:          64,
		Height:         32,
		Pc:             0x200,
		Sp:             0,
		I:              0,
		PrimaryColor:   rl.White,
		SecondaryColor: rl.Black,
	}

	return chip
}

type EmulatorStore interface {
	ClearScreen()
	LoadRegister(firstByte, secondByte byte)
	LoadIndexRegister(firstByte, secondByte byte)
	JumpToInstruction(firstByte, secondByte byte)
	Draw(firstByte, secondByte byte)
	AddValueToRegister(firstByte, secondByte byte)
	SkipNextInstruction(firstByte, secondByte byte)
	JumpPlusRegister(firstByte, secondByte byte)
	CallAddress(firstByte, secondByte byte)
	SkipIfNotEquals(firstByte, secondByte byte)
	SkipEqualRegisters(firstByte, secondByte byte)
	SkipNotEqualRegisters(firstByte, secondByte byte)
	Return(firstByte, secondByte byte)
	VxGetsVy(firstByte, secondByte byte)
	LoadRegistersFromMemory(firstByte, secondByte byte)
	VxOrVy(firstByte, secondByte byte)
	VxAndVy(firstByte, secondByte byte)
	VxXorVy(firstByte, secondByte byte)
}

type Emulator struct {
	EmulatorStore
}

// ClearScreen clears the screen by setting all pixels to 0.
func (c *Chip8) ClearScreen() {
	for i := range c.Screen {
		c.Screen[i] = c.SecondaryColor
	}
}

func (c *Chip8) LoadRegistersFromMemory(firstByte, secondByte byte) {
	numOfRegisters := (firstByte & 0xf) + 1

	for i := 0; i < int(numOfRegisters); i++ {
		c.Registers[i] = c.Memory[c.I]
		c.I += 1
	}
}

// Return pops address from the stack and puts it in the pc.
func (c *Chip8) Return(firstByte, secondByte byte) {
	address := c.Stack[len(c.Stack)-1]
	c.Stack = c.Stack[:len(c.Stack)-1]
	c.Pc = address
}

func (c *Chip8) CallAddress(firstByte, secondByte byte) {
	c.Stack = append(c.Stack, c.Pc)
	c.Pc = get12BitValue(firstByte, secondByte)
}

// SkipEqualRegisters compares values of two registers and increases pc if they're equal.
func (c *Chip8) SkipEqualRegisters(firstByte, secondByte byte) {
	firstRegister := firstByte & 0xf
	secondRegister := secondByte >> 4

	if c.Registers[firstRegister] == c.Registers[secondRegister] {
		c.Pc += 2
	}
}

// SkipIfNotEquals increases the program counter by 2 if value in given register is different than secondByte.
func (c *Chip8) SkipIfNotEquals(firstByte, secondByte byte) {
	register := firstByte & 0xf
	registerValue := c.Registers[register]
	value := secondByte

	if registerValue != value {
		c.Pc += 2
	}
}

// SkipNotEqualRegisters increases pc if values of registers are different.
func (c *Chip8) SkipNotEqualRegisters(firstByte, secondByte byte) {
	firstRegister := firstByte & 0xf
	firstValue := c.Registers[firstRegister]
	secondRegister := secondByte >> 4
	secondValue := c.Registers[secondRegister]

	if firstValue != secondValue {
		c.Pc += 2
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

// VxGetsVy stores value of y register in x register.
func (c *Chip8) VxGetsVy(firstByte, secondByte byte) {
	registerX := firstByte & 0xf
	registerY := secondByte >> 4
	registerYValue := c.Registers[registerY]
	c.Registers[registerX] = registerYValue
}

// JumpToInstruction sets the program counter to the new value.
func (c *Chip8) JumpToInstruction(firstByte, secondByte byte) {
	c.Pc = get12BitValue(firstByte, secondByte)
}

// SkipNextInstruction increases the program counter if value of the register is equal to secondByte.
func (c *Chip8) SkipNextInstruction(firstByte, secondByte byte) {
	registerValue := c.Registers[firstByte&0xf]
	value := secondByte

	if registerValue == value {
		c.Pc += 2
	}
}

func (c *Chip8) JumpPlusRegister(firstByte, secondByte byte) {
	register := c.Registers[0x0]
	c.Pc = get12BitValue(firstByte, secondByte) + uint16(register)
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
				color = c.PrimaryColor
			} else {
				color = c.SecondaryColor
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

func (c *Chip8) AddValueToRegister(firstByte, secondByte byte) {
	register := firstByte & 0xf
	value := secondByte

	c.Registers[register] += value
}

func (c *Chip8) VxOrVy(firstByte, secondByte byte) {
	registerX := firstByte & 0xf

	value := c.Registers[registerX] | c.Registers[secondByte>>4]
	c.Registers[registerX] = value
}

// VxAndVy calculates result of Vx&Vy and stores the result in Vx.
func (c *Chip8) VxAndVy(firstByte, secondByte byte) {
	registerX := firstByte & 0xf

	value := c.Registers[registerX] & c.Registers[secondByte>>4]
	c.Registers[registerX] = value
}

// VxXorVy calculates result of Vx^Vy and stores the result in Vx.
func (c *Chip8) VxXorVy(firstByte, secondByte byte) {
	registerX := firstByte & 0xf

	value := c.Registers[registerX] ^ c.Registers[secondByte>>4]
	c.Registers[registerX] = value
}

func (e *Emulator) Emulate(firstByte, secondByte byte) {
	switch firstByte >> 4 {
	case 0x0:
		switch secondByte {
		case 0xe0:
			e.ClearScreen()
		case 0xee:
			e.Return(firstByte, secondByte)
		}

	case 0x1:
		e.JumpToInstruction(firstByte, secondByte)
	case 0x2:
		e.CallAddress(firstByte, secondByte)
	case 0x3:
		e.SkipNextInstruction(firstByte, secondByte)
	case 0x4:
		e.SkipIfNotEquals(firstByte, secondByte)
	case 0x5:
		e.SkipEqualRegisters(firstByte, secondByte)
	case 0x6:
		e.LoadRegister(firstByte, secondByte)
	case 0x7:
		e.AddValueToRegister(firstByte, secondByte)
	case 0x8:
		switch secondByte & 0xf {
		case 0x0:
			e.VxGetsVy(firstByte, secondByte)
		case 0x1:
			e.VxOrVy(firstByte, secondByte)
		case 0x2:
			e.VxAndVy(firstByte, secondByte)
		case 0x3:
			e.VxXorVy(firstByte, secondByte)
		default:
			panic(fmt.Sprintf("Instruction %x not implemented", uint16(firstByte)<<8|uint16(secondByte)))
		}
	case 0x9:
		e.SkipNotEqualRegisters(firstByte, secondByte)
	case 0xa:
		e.LoadIndexRegister(firstByte, secondByte)
	case 0xb:
		e.JumpPlusRegister(firstByte, secondByte)
	case 0xd:
		e.Draw(firstByte, secondByte)
	case 0xf:
		switch secondByte {
		case 0x65:
			e.LoadRegistersFromMemory(firstByte, secondByte)
		}
	default:
		fmt.Printf("Instruction %x not implemented\n", uint16(firstByte)<<8|uint16(secondByte))
		panic(fmt.Sprintf("Instruction %x not implemented", uint16(firstByte)<<8|uint16(secondByte)))
	}

}
