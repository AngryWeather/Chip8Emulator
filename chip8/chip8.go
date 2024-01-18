package chip8

import (
	"fmt"
	"image/color"
	"math/rand"

	"os"

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
		Stack:          make([]uint16, 0, 48),
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
	LoadRegistersToMemory(firstByte, secondByte byte)
	VxOrVy(firstByte, secondByte byte)
	VxAndVy(firstByte, secondByte byte)
	VxXorVy(firstByte, secondByte byte)
	VxAddVy(firstByte, secondByte byte)
	VxSubVy(firstByte, secondByte byte)
	VySubVx(firstByte, secondByte byte)
	VxRightShift(firstByte, secondByte byte)
	VxLeftShift(firstByte, secondByte byte)
	StoreBCDRepresentationInMemory(firstByte, secondByte byte)
	StoreValueOfVxPlusIInI(firstByte, secondByte byte)
	SetDelayTimer(firstByte byte)
	SkipKeyNotPressed(firstByte byte)
	SkipKeyPressed(firstByte byte)
	PutTimerInRegister(firstByte byte)
	WaitForKeyPress(firstByte byte)
	SetRandomNumber(firstByte, secondByte byte)
}

type Emulator struct {
	EmulatorStore
}

var keymap = map[byte]int32{
	0x1: rl.KeyOne,
	0x2: rl.KeyTwo,
	0x3: rl.KeyThree,
	0xc: rl.KeyFour,
	0x4: rl.KeyQ,
	0x5: rl.KeyW,
	0x6: rl.KeyE,
	0xd: rl.KeyR,
	0x7: rl.KeyA,
	0x8: rl.KeyS,
	0x9: rl.KeyD,
	0xe: rl.KeyF,
	0xa: rl.KeyZ,
	0x0: rl.KeyX,
	0xb: rl.KeyC,
	0xf: rl.KeyV,
}

func (c *Chip8) SetRandomNumber(firstByte, secondByte byte) {
	randNumber := rand.Intn(256)
	c.Registers[firstByte&0xf] = byte(randNumber) & secondByte
}

// ClearScreen clears the screen by setting all pixels to 0.
func (c *Chip8) ClearScreen() {
	for i := range c.Screen {
		c.Screen[i] = c.SecondaryColor
	}
}

// StoreBCDRepresentationInMemory stores decimal number in Vx in Memory (Memory[i] = hundreds digit, Memory[i+1] = tens digit, Memory[i+2] = ones digit).
func (c *Chip8) StoreBCDRepresentationInMemory(firstByte, secondByte byte) {
	registerX := firstByte & 0xf
	valueRegisterX := c.Registers[registerX]

	var i int = 0
	for i = int(c.I + 2); i >= int(c.I); i-- {
		// get last digit from valueRegisterX
		c.Memory[i] = valueRegisterX % 10
		// divide to get the next digit from the right
		valueRegisterX /= 10
	}
}

// StoreValueOfVxPlusIInI adds values of index register and Vx and stores the result in index register I.
func (c *Chip8) StoreValueOfVxPlusIInI(firstByte, secondByte byte) {
	c.I = c.I + uint16(c.Registers[firstByte&0xf])
}

// LoadRegistersFromMemory loads x registers from memory starting at index register (I).
func (c *Chip8) LoadRegistersFromMemory(firstByte, secondByte byte) {
	numOfRegisters := (firstByte & 0xf) + 1

	for i := 0; i < int(numOfRegisters); i++ {
		c.Registers[i] = c.Memory[c.I]
		c.I += 1
	}
}

// LoadRegistersToMemory loads x registers to memory starting at index register I.
func (c *Chip8) LoadRegistersToMemory(firstByte, secondByte byte) {
	numOfRegisters := (firstByte & 0xf) + 1

	for i := 0; i < int(numOfRegisters); i++ {
		c.Memory[c.I] = c.Registers[i]
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
	// works for chip8 quirks
	register := c.Registers[0xf]
	c.Pc = get12BitValue(firstByte, secondByte) + uint16(register)
}

func (c *Chip8) Draw(firstByte, secondByte byte) {
	bytesToRead := secondByte & 0xf
	x := c.Registers[firstByte&0xf] % c.Width
	y := c.Registers[secondByte>>4] % c.Height
	c.Registers[0xf] = 0

	for i := c.I; i < (uint16(bytesToRead) + c.I); i++ {
		var currentByte byte = c.Memory[i]
		var color color.RGBA

		// check each bit in the current byte
		for j := 0; j < 8; j++ {
			pixel := currentByte >> 7
			// shift byte to access next bit from left
			currentByte = currentByte << 1

			// position in 1D array is based on x, y and width
			var position int = int(x) + (int(y) * int(c.Width))

			if (pixel == 1 && c.Screen[position] == c.SecondaryColor) || (pixel == 0 && c.Screen[position] == c.PrimaryColor) {
				color = c.PrimaryColor
			} else {
				color = c.SecondaryColor
			}
			// pixels are xored (^) onto the screen but xor is not defined for color.RGBA
			if pixel == 1 && c.Screen[position] == c.PrimaryColor {
				c.Registers[0xf] = 1
			}
			// c.Screen[position].R ^= color.R
			// c.Screen[position].G ^= color.G
			// c.Screen[position].B ^= color.B
			// c.Screen[position].A ^= color.A
			c.Screen[position] = color

			// increase x to draw in the next x coordinate
			x += 1
			if x > 63 {
				break
			}
		}
		// reset x
		x = c.Registers[firstByte&0xf] % c.Width
		// increase y to move down
		y += 1
		if y > 31 {
			break
		}

	}
}

func (c *Chip8) AddValueToRegister(firstByte, secondByte byte) {
	register := firstByte & 0xf
	value := secondByte

	c.Registers[register] += value
}

func (c *Chip8) VxOrVy(firstByte, secondByte byte) {
	registerX := firstByte & 0xf

	c.Registers[0xf] = 0
	value := c.Registers[registerX] | c.Registers[secondByte>>4]
	c.Registers[registerX] = value
}

// VxAndVy calculates result of Vx&Vy and stores the result in Vx.
func (c *Chip8) VxAndVy(firstByte, secondByte byte) {
	registerX := firstByte & 0xf

	c.Registers[0xf] = 0
	value := c.Registers[registerX] & c.Registers[secondByte>>4]
	c.Registers[registerX] = value
}

// VxXorVy calculates result of Vx^Vy and stores the result in Vx.
func (c *Chip8) VxXorVy(firstByte, secondByte byte) {
	registerX := firstByte & 0xf
	c.Registers[0xf] = 0

	value := c.Registers[registerX] ^ c.Registers[secondByte>>4]
	c.Registers[registerX] = value
}

// VxAddVy adds value of Vx and Vy, stores the result in Vx and sets Vf to 1 on overflow.
func (c *Chip8) VxAddVy(firstByte, secondByte byte) {
	registerX := firstByte & 0xf
	registerY := secondByte >> 4

	var xOverflow int = int(c.Registers[registerX])
	var yOverflow int = int(c.Registers[registerY])
	c.Registers[registerX] = c.Registers[registerX] + c.Registers[registerY]

	if xOverflow+yOverflow > 255 {
		c.Registers[0xf] = 1
	} else {
		c.Registers[0xf] = 0
	}

}

// VxSubVy sets Vf to 1 if Vx > Vy, stores result of Vx - Vy in Vx.
func (c *Chip8) VxSubVy(firstByte, secondByte byte) {
	registerX := firstByte & 0xf
	registerY := secondByte >> 4

	var xOverflow int = int(c.Registers[registerX])
	var yOverflow int = int(c.Registers[registerY])
	c.Registers[registerX] = c.Registers[registerX] - c.Registers[registerY]

	if xOverflow >= yOverflow {
		c.Registers[0xf] = 1
		// y > x sets Vf to 0 because of borrowing value
	} else if xOverflow < yOverflow {
		c.Registers[0xf] = 0
	}
}

// VySubVx sets Vf to 1 if Vy > Vx and sets Vx to Vy - Vx.
func (c *Chip8) VySubVx(firstByte, secondByte byte) {
	registerX := firstByte & 0xf
	registerY := secondByte >> 4

	c.Registers[registerX] = c.Registers[registerY] - c.Registers[registerX]

	if c.Registers[registerY] > c.Registers[registerX] {
		c.Registers[0xf] = 1
	} else {
		c.Registers[0xf] = 0
	}
}

// VxRightShift sets Vf to 1 if the least significant bit of Vx is 1 and divides Vx by 2.
func (c *Chip8) VxRightShift(firstByte, secondByte byte) {
	registerX := firstByte & 0xf

	c.Registers[registerX] = c.Registers[secondByte>>4]
	// find least significant bit and check if it's 1
	c.Registers[0xf] = c.Registers[registerX] & 0x1
	// right shift by 1 to divide by 2
	c.Registers[registerX] = c.Registers[registerX] >> 1
}

// VxLeftShift sets Vf to 1 if the most significant bit of Vx is 1 and multiplies Vx by 2.
func (c *Chip8) VxLeftShift(firstByte, secondByte byte) {
	registerX := firstByte & 0xf

	c.Registers[registerX] = c.Registers[secondByte>>4]
	// find most significant bit and check if it's 1
	c.Registers[0xf] = c.Registers[registerX] >> 7
	// left shift by 1 to multiply by 2
	c.Registers[registerX] = c.Registers[registerX] << 1

}

// SetDelayTimer sets delay timer to value in Vx.
func (c *Chip8) SetDelayTimer(firstByte byte) {
	c.Timers[0] = c.Registers[firstByte&0xf]
}

func (c *Chip8) PutTimerInRegister(firstByte byte) {
	c.Registers[firstByte&0xf] = c.Timers[0]
}

func (c *Chip8) SkipKeyNotPressed(firstByte byte) {
	targetKey := c.Registers[firstByte&0xf]
	if rl.IsKeyUp(keymap[targetKey]) {
		c.Pc += 2
	}
}

func (c *Chip8) SkipKeyPressed(firstByte byte) {
	targetKey := c.Registers[firstByte&0xf]
	if rl.IsKeyDown(keymap[targetKey]) {
		c.Pc += 2
	}
}

func (c *Chip8) WaitForKeyPress(firstByte byte) {
	if rl.GetKeyPressed() != 0 {
		value := rl.GetKeyPressed()
		c.Registers[firstByte&0xf] = getKeyFromKeymap(value)
	} else {
		c.Pc -= 2
	}
}

func getKeyFromKeymap(value int32) byte {
	for k, v := range keymap {
		if v == value {
			fmt.Fprintln(os.Stdout, []any{"v: %d", v}...)
			return k
		}
	}
	return 0
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
		case 0x4:
			e.VxAddVy(firstByte, secondByte)
		case 0x5:
			e.VxSubVy(firstByte, secondByte)
		case 0x6:
			e.VxRightShift(firstByte, secondByte)
		case 0x7:
			e.VySubVx(firstByte, secondByte)
		case 0xe:
			e.VxLeftShift(firstByte, secondByte)
		}
	case 0x9:
		e.SkipNotEqualRegisters(firstByte, secondByte)
	case 0xa:
		e.LoadIndexRegister(firstByte, secondByte)
	case 0xb:
		e.JumpPlusRegister(firstByte, secondByte)
	case 0xc:
		e.SetRandomNumber(firstByte, secondByte)
	case 0xd:
		e.Draw(firstByte, secondByte)
	case 0xe:
		switch secondByte {
		case 0x9e:
			e.SkipKeyPressed(firstByte)
		case 0xa1:
			e.SkipKeyNotPressed(firstByte)
		}
	case 0xf:
		switch secondByte {
		case 0x07:
			e.PutTimerInRegister(firstByte)
		case 0x0a:
			e.WaitForKeyPress(firstByte)
		case 0x15:
			e.SetDelayTimer(firstByte)
		case 0x18:
			panic(fmt.Sprintf("Instruction %x not implemented", uint16(firstByte)<<8|uint16(secondByte)))
		case 0x1e:
			e.StoreValueOfVxPlusIInI(firstByte, secondByte)
		case 0x33:
			e.StoreBCDRepresentationInMemory(firstByte, secondByte)
		case 0x55:
			e.LoadRegistersToMemory(firstByte, secondByte)
		case 0x65:
			e.LoadRegistersFromMemory(firstByte, secondByte)
		}
	}
}
